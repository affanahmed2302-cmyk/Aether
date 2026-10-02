package raft

import (
	"encoding/json"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/affanahmed2302-cmyk/Aether/internal/kv"
)

type State int

const (
	Follower State = iota
	Candidate
	Leader
)

func (s State) String() string {
	switch s {
	case Follower:
		return "Follower"
	case Candidate:
		return "Candidate"
	case Leader:
		return "Leader"
	default:
		return "Unknown"
	}
}

type LogEntry struct {
	Term    int64  `json:"term"`
	Index   int64  `json:"index"`
	Command []byte `json:"command"`
}

type Node struct {
	mu sync.Mutex

	ID    string
	Peers []string // addresses of other nodes

	// Persistent state
	CurrentTerm int64
	VotedFor    string
	Log         []LogEntry

	// Volatile
	CommitIndex int64
	LastApplied int64
	State       State

	// Leader state
	NextIndex  map[string]int64
	MatchIndex map[string]int64

	applyCh chan LogEntry
	stopCh  chan struct{}

	// Transport callbacks (set by the server)
	SendRequestVote   func(addr string, args RequestVoteArgs) (RequestVoteReply, error)
	SendAppendEntries func(addr string, args AppendEntriesArgs) (AppendEntriesReply, error)

	lastHeartbeat time.Time
	electionTimeout time.Duration
}

type RequestVoteArgs struct {
	Term         int64
	CandidateID  string
	LastLogIndex int64
	LastLogTerm  int64
}

type RequestVoteReply struct {
	Term        int64
	VoteGranted bool
}

type AppendEntriesArgs struct {
	Term         int64
	LeaderID     string
	PrevLogIndex int64
	PrevLogTerm  int64
	Entries      []LogEntry
	LeaderCommit int64
}

type AppendEntriesReply struct {
	Term    int64
	Success bool
}

func NewNode(id string, peers []string, applyCh chan LogEntry) *Node {
	n := &Node{
		ID:              id,
		Peers:           peers,
		State:           Follower,
		Log:             make([]LogEntry, 0),
		NextIndex:       make(map[string]int64),
		MatchIndex:      make(map[string]int64),
		applyCh:         applyCh,
		stopCh:          make(chan struct{}),
		lastHeartbeat:   time.Now(),
		electionTimeout: time.Duration(250+rand.Intn(200)) * time.Millisecond,
	}
	return n
}

func (n *Node) Start() {
	go n.loop()
}

func (n *Node) Stop() {
	close(n.stopCh)
}

func (n *Node) loop() {
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-n.stopCh:
			return
		case <-ticker.C:
			n.tick()
		}
	}
}

func (n *Node) tick() {
	n.mu.Lock()
	defer n.mu.Unlock()

	switch n.State {
	case Follower, Candidate:
		if time.Since(n.lastHeartbeat) > n.electionTimeout {
			n.startElection()
		}
	case Leader:
		n.broadcastAppendEntries()
	}
}

func (n *Node) startElection() {
	n.State = Candidate
	n.CurrentTerm++
	n.VotedFor = n.ID
	n.lastHeartbeat = time.Now()
	n.electionTimeout = time.Duration(250+rand.Intn(200)) * time.Millisecond
	term := n.CurrentTerm
	lastIdx, lastTerm := n.lastLogInfo()

	log.Printf("[%s] starting election for term %d", n.ID, term)

	votes := 1 // vote for self
	var mu sync.Mutex

	for _, peer := range n.Peers {
		go func(addr string) {
			if n.SendRequestVote == nil {
				return
			}
			args := RequestVoteArgs{
				Term:         term,
				CandidateID:  n.ID,
				LastLogIndex: lastIdx,
				LastLogTerm:  lastTerm,
			}
			reply, err := n.SendRequestVote(addr, args)
			if err != nil {
				return
			}
			n.mu.Lock()
			defer n.mu.Unlock()
			if reply.Term > n.CurrentTerm {
				n.becomeFollower(reply.Term)
				return
			}
			if n.State == Candidate && reply.VoteGranted && reply.Term == term {
				mu.Lock()
				votes++
				if votes > (len(n.Peers)+1)/2 {
					n.becomeLeader()
				}
				mu.Unlock()
			}
		}(peer)
	}
}

func (n *Node) becomeLeader() {
	n.State = Leader
	log.Printf("[%s] became LEADER for term %d", n.ID, n.CurrentTerm)
	for _, p := range n.Peers {
		n.NextIndex[p] = n.lastLogIndex() + 1
		n.MatchIndex[p] = 0
	}
}

func (n *Node) becomeFollower(term int64) {
	n.State = Follower
	n.CurrentTerm = term
	n.VotedFor = ""
	n.lastHeartbeat = time.Now()
}

func (n *Node) broadcastAppendEntries() {
	for _, peer := range n.Peers {
		go n.sendAppendEntriesTo(peer)
	}
}

func (n *Node) sendAppendEntriesTo(addr string) {
	n.mu.Lock()
	if n.State != Leader || n.SendAppendEntries == nil {
		n.mu.Unlock()
		return
	}
	nextIdx := n.NextIndex[addr]
	prevIdx := nextIdx - 1
	var prevTerm int64
	if prevIdx > 0 && int(prevIdx) <= len(n.Log) {
		prevTerm = n.Log[prevIdx-1].Term
	}
	entries := []LogEntry{}
	if nextIdx > 0 && int(nextIdx) <= len(n.Log) {
		entries = append(entries, n.Log[nextIdx-1:]...)
	}
	args := AppendEntriesArgs{
		Term:         n.CurrentTerm,
		LeaderID:     n.ID,
		PrevLogIndex: prevIdx,
		PrevLogTerm:  prevTerm,
		Entries:      entries,
		LeaderCommit: n.CommitIndex,
	}
	n.mu.Unlock()

	reply, err := n.SendAppendEntries(addr, args)
	if err != nil {
		return
	}

	n.mu.Lock()
	defer n.mu.Unlock()
	if reply.Term > n.CurrentTerm {
		n.becomeFollower(reply.Term)
		return
	}
	if reply.Success {
		if len(entries) > 0 {
			n.MatchIndex[addr] = entries[len(entries)-1].Index
			n.NextIndex[addr] = n.MatchIndex[addr] + 1
		}
		n.updateCommitIndex()
	} else {
		if n.NextIndex[addr] > 1 {
			n.NextIndex[addr]--
		}
	}
}

func (n *Node) updateCommitIndex() {
	// Find highest index replicated on majority
	for i := int64(len(n.Log)); i > n.CommitIndex; i-- {
		count := 1 // self
		for _, p := range n.Peers {
			if n.MatchIndex[p] >= i {
				count++
			}
		}
		if count > (len(n.Peers)+1)/2 && n.Log[i-1].Term == n.CurrentTerm {
			n.CommitIndex = i
			n.applyCommitted()
			break
		}
	}
}

func (n *Node) applyCommitted() {
	for n.LastApplied < n.CommitIndex {
		n.LastApplied++
		entry := n.Log[n.LastApplied-1]
		n.applyCh <- entry
	}
}

// Propose is called by the application when the node is leader.
func (n *Node) Propose(cmd []byte) (int64, error) {
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.State != Leader {
		return 0, ErrNotLeader
	}
	entry := LogEntry{
		Term:    n.CurrentTerm,
		Index:   n.lastLogIndex() + 1,
		Command: cmd,
	}
	n.Log = append(n.Log, entry)
	// Self match
	n.MatchIndex[n.ID] = entry.Index
	return entry.Index, nil
}

func (n *Node) HandleRequestVote(args RequestVoteArgs) RequestVoteReply {
	n.mu.Lock()
	defer n.mu.Unlock()
	reply := RequestVoteReply{Term: n.CurrentTerm}

	if args.Term < n.CurrentTerm {
		return reply
	}
	if args.Term > n.CurrentTerm {
		n.becomeFollower(args.Term)
	}
	lastIdx, lastTerm := n.lastLogInfo()
	upToDate := args.LastLogTerm > lastTerm || (args.LastLogTerm == lastTerm && args.LastLogIndex >= lastIdx)
	if (n.VotedFor == "" || n.VotedFor == args.CandidateID) && upToDate {
		n.VotedFor = args.CandidateID
		n.lastHeartbeat = time.Now()
		reply.VoteGranted = true
		reply.Term = n.CurrentTerm
		log.Printf("[%s] granted vote to %s for term %d", n.ID, args.CandidateID, args.Term)
	}
	return reply
}

func (n *Node) HandleAppendEntries(args AppendEntriesArgs) AppendEntriesReply {
	n.mu.Lock()
	defer n.mu.Unlock()
	reply := AppendEntriesReply{Term: n.CurrentTerm}

	if args.Term < n.CurrentTerm {
		return reply
	}
	if args.Term > n.CurrentTerm || n.State != Follower {
		n.becomeFollower(args.Term)
	}
	n.lastHeartbeat = time.Now()
	reply.Term = n.CurrentTerm

	// Simple append for educational version
	if len(args.Entries) > 0 {
		n.Log = append(n.Log, args.Entries...)
	}
	if args.LeaderCommit > n.CommitIndex {
		n.CommitIndex = min(args.LeaderCommit, n.lastLogIndex())
		n.applyCommitted()
	}
	reply.Success = true
	return reply
}

func (n *Node) lastLogIndex() int64 {
	if len(n.Log) == 0 {
		return 0
	}
	return n.Log[len(n.Log)-1].Index
}

func (n *Node) lastLogInfo() (int64, int64) {
	if len(n.Log) == 0 {
		return 0, 0
	}
	e := n.Log[len(n.Log)-1]
	return e.Index, e.Term
}

func (n *Node) IsLeader() bool {
	n.mu.Lock()
	defer n.mu.Unlock()
	return n.State == Leader
}

func (n *Node) GetState() (State, int64) {
	n.mu.Lock()
	defer n.mu.Unlock()
	return n.State, n.CurrentTerm
}

var ErrNotLeader = &Error{"not the leader"}

type Error struct{ s string }

func (e *Error) Error() string { return e.s }

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

// Helper for client commands
func EncodeCommand(op, key, value string) []byte {
	c := kv.Command{Op: op, Key: key, Value: value}
	b, _ := json.Marshal(c)
	return b
}

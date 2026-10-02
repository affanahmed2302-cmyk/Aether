package raft

import (
	"log"
	"sync"
	"time"
)

// NodeState represents the current role of a Raft node.
type NodeState int

const (
	Follower NodeState = iota
	Candidate
	Leader
)

func (s NodeState) String() string {
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

// LogEntry is a single entry in the replicated log.
type LogEntry struct {
	Term    int64
	Index   int64
	Command []byte
}

// Node is the core Raft state machine for one server.
type Node struct {
	mu sync.Mutex

	// Persistent state
	currentTerm int64
	votedFor    string
	log         []LogEntry

	// Volatile state
	commitIndex int64
	lastApplied int64

	// Leader state
	nextIndex  map[string]int64
	matchIndex map[string]int64

	id       string
	peers    []string
	state    NodeState
	electionTimeout time.Duration
	heartbeatInterval time.Duration

	applyCh chan LogEntry // applied entries go to the state machine
	stopCh  chan struct{}
}

// NewNode creates a new Raft node.
func NewNode(id string, peers []string, applyCh chan LogEntry) *Node {
	n := &Node{
		id:                id,
		peers:             peers,
		state:             Follower,
		log:               make([]LogEntry, 0),
		nextIndex:         make(map[string]int64),
		matchIndex:        make(map[string]int64),
		electionTimeout:   300 * time.Millisecond,
		heartbeatInterval: 100 * time.Millisecond,
		applyCh:           applyCh,
		stopCh:            make(chan struct{}),
	}
	return n
}

// Start begins the Raft main loop.
func (n *Node) Start() {
	go n.run()
}

func (n *Node) run() {
	for {
		select {
		case <-n.stopCh:
			return
		default:
			n.mu.Lock()
			state := n.state
			n.mu.Unlock()

			switch state {
			case Follower:
				n.runFollower()
			case Candidate:
				n.runCandidate()
			case Leader:
				n.runLeader()
			}
		}
	}
}

func (n *Node) runFollower() {
	timer := time.NewTimer(n.electionTimeout)
	defer timer.Stop()
	select {
	case <-timer.C:
		n.mu.Lock()
		n.state = Candidate
		log.Printf("[%s] election timeout → becoming Candidate (term %d)", n.id, n.currentTerm+1)
		n.mu.Unlock()
	case <-n.stopCh:
		return
	}
}

func (n *Node) runCandidate() {
	n.mu.Lock()
	n.currentTerm++
	n.votedFor = n.id
	term := n.currentTerm
	n.mu.Unlock()

	log.Printf("[%s] started election for term %d", n.id, term)

	// In a full implementation we would send RequestVote RPCs here.
	// For the educational skeleton we simulate a successful election after a short delay
	// when the node is the only one (or for demo purposes).
	time.Sleep(50 * time.Millisecond)

	n.mu.Lock()
	defer n.mu.Unlock()
	if n.state == Candidate && n.currentTerm == term {
		n.state = Leader
		log.Printf("[%s] became Leader for term %d", n.id, term)
		for _, p := range n.peers {
			n.nextIndex[p] = int64(len(n.log)) + 1
			n.matchIndex[p] = 0
		}
	}
}

func (n *Node) runLeader() {
	ticker := time.NewTicker(n.heartbeatInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			// Send heartbeats / AppendEntries (skeleton)
			log.Printf("[%s] heartbeat (Leader term %d)", n.id, n.currentTerm)
		case <-n.stopCh:
			return
		}
		n.mu.Lock()
		if n.state != Leader {
			n.mu.Unlock()
			return
		}
		n.mu.Unlock()
	}
}

// Propose is called by the leader to append a new command.
func (n *Node) Propose(cmd []byte) (int64, error) {
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.state != Leader {
		return 0, ErrNotLeader
	}
	entry := LogEntry{
		Term:    n.currentTerm,
		Index:   int64(len(n.log)) + 1,
		Command: cmd,
	}
	n.log = append(n.log, entry)
	// In full Raft we would replicate to majority before committing.
	// Here we immediately mark as committed for the single-node educational path.
	n.commitIndex = entry.Index
	n.applyCh <- entry
	return entry.Index, nil
}

var ErrNotLeader = &RaftError{"not the leader"}

type RaftError struct{ msg string }

func (e *RaftError) Error() string { return e.msg }

// Stop shuts down the node.
func (n *Node) Stop() {
	close(n.stopCh)
}

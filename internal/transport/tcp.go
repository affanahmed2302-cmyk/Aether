package transport

import (
	"encoding/gob"
	"fmt"
	"net"
	"time"

	"github.com/affanahmed2302-cmyk/Aether/internal/raft"
)

func init() {
	gob.Register(raft.RequestVoteArgs{})
	gob.Register(raft.RequestVoteReply{})
	gob.Register(raft.AppendEntriesArgs{})
	gob.Register(raft.AppendEntriesReply{})
	gob.Register(raft.LogEntry{})
}

type RPCType int

const (
	RPCRequestVote RPCType = iota
	RPCAppendEntries
	RPCClientPropose
	RPCClientGet
)

type Request struct {
	Type RPCType
	Args interface{}
}

type Response struct {
	Reply interface{}
	Err   string
}

// Server handles incoming TCP RPCs for a Raft node.
type Server struct {
	node     *raft.Node
	listener net.Listener
}

func NewServer(node *raft.Node) *Server {
	return &Server{node: node}
}

func (s *Server) Start(addr string) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	s.listener = ln
	go s.acceptLoop()
	return nil
}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()
	dec := gob.NewDecoder(conn)
	enc := gob.NewEncoder(conn)
	var req Request
	if err := dec.Decode(&req); err != nil {
		return
	}
	var resp Response
	switch req.Type {
	case RPCRequestVote:
		args := req.Args.(raft.RequestVoteArgs)
		reply := s.node.HandleRequestVote(args)
		resp.Reply = reply
	case RPCAppendEntries:
		args := req.Args.(raft.AppendEntriesArgs)
		reply := s.node.HandleAppendEntries(args)
		resp.Reply = reply
	case RPCClientPropose:
		cmd := req.Args.([]byte)
		idx, err := s.node.Propose(cmd)
		if err != nil {
			resp.Err = err.Error()
		} else {
			resp.Reply = idx
		}
	default:
		resp.Err = "unknown RPC"
	}
	_ = enc.Encode(resp)
}

func (s *Server) Close() {
	if s.listener != nil {
		s.listener.Close()
	}
}

// Client helpers
func call(addr string, typ RPCType, args interface{}) (interface{}, error) {
	conn, err := net.DialTimeout("tcp", addr, 800*time.Millisecond)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(1 * time.Second))
	enc := gob.NewEncoder(conn)
	dec := gob.NewDecoder(conn)
	if err := enc.Encode(Request{Type: typ, Args: args}); err != nil {
		return nil, err
	}
	var resp Response
	if err := dec.Decode(&resp); err != nil {
		return nil, err
	}
	if resp.Err != "" {
		return nil, fmt.Errorf(resp.Err)
	}
	return resp.Reply, nil
}

func SendRequestVote(addr string, args raft.RequestVoteArgs) (raft.RequestVoteReply, error) {
	reply, err := call(addr, RPCRequestVote, args)
	if err != nil {
		return raft.RequestVoteReply{}, err
	}
	return reply.(raft.RequestVoteReply), nil
}

func SendAppendEntries(addr string, args raft.AppendEntriesArgs) (raft.AppendEntriesReply, error) {
	reply, err := call(addr, RPCAppendEntries, args)
	if err != nil {
		return raft.AppendEntriesReply{}, err
	}
	return reply.(raft.AppendEntriesReply), nil
}

func Propose(addr string, cmd []byte) (int64, error) {
	reply, err := call(addr, RPCClientPropose, cmd)
	if err != nil {
		return 0, err
	}
	return reply.(int64), nil
}

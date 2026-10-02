package kv

import (
	"encoding/json"
	"sync"

	"github.com/affanahmed2302-cmyk/Aether/internal/raft"
)

// Command represents a single operation on the key-value store.
type Command struct {
	Op    string `json:"op"` // "put", "delete", "get"
	Key   string `json:"key"`
	Value string `json:"value,omitempty"`
}

// Store is the deterministic state machine that applies committed Raft log entries.
type Store struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewStore() *Store {
	return &Store{data: make(map[string]string)}
}

// Apply takes a committed Raft log entry and updates the state machine.
func (s *Store) Apply(entry raft.LogEntry) {
	var cmd Command
	if err := json.Unmarshal(entry.Command, &cmd); err != nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	switch cmd.Op {
	case "put":
		s.data[cmd.Key] = cmd.Value
	case "delete":
		delete(s.data, cmd.Key)
	}
}

func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.data[key]
	return v, ok
}

func (s *Store) Snapshot() map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	cp := make(map[string]string, len(s.data))
	for k, v := range s.data {
		cp[k] = v
	}
	return cp
}

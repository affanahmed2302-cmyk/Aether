package kv

import (
	"encoding/json"
	"sync"

	"github.com/affanahmed2302-cmyk/Aether/internal/raft"
)

type Command struct {
	Op    string `json:"op"`
	Key   string `json:"key"`
	Value string `json:"value,omitempty"`
}

type Store struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewStore() *Store {
	return &Store{data: make(map[string]string)}
}

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

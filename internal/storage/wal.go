package storage

import (
	"encoding/binary"
	"encoding/json"
	"os"
	"sync"

	"github.com/affanahmed2302-cmyk/Aether/internal/raft"
)

// WAL is a simple append-only write-ahead log for durability.
type WAL struct {
	mu   sync.Mutex
	file *os.File
}

func OpenWAL(path string) (*WAL, error) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	return &WAL{file: f}, nil
}

func (w *WAL) Append(entry raft.LogEntry) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	lenBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBuf, uint32(len(data)))
	if _, err := w.file.Write(lenBuf); err != nil {
		return err
	}
	if _, err := w.file.Write(data); err != nil {
		return err
	}
	return w.file.Sync() // force durability
}

func (w *WAL) Close() error {
	return w.file.Close()
}

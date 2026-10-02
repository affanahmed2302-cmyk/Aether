package metrics

import (
	"sync"
	"sync/atomic"
	"time"
)

// Metrics collects runtime signals that top-tier interns are expected to show.
type Metrics struct {
	LeaderElections   atomic.Int64
	CommandsProposed  atomic.Int64
	CommandsCommitted atomic.Int64
	ApplyCount        atomic.Int64

	mu          sync.Mutex
	latencies   []time.Duration
maxLatency  time.Duration
}

func New() *Metrics {
	return &Metrics{
		latencies: make([]time.Duration, 0, 1024),
	}
}

func (m *Metrics) IncElection() {
	m.LeaderElections.Add(1)
}

func (m *Metrics) IncPropose() {
	m.CommandsProposed.Add(1)
}

func (m *Metrics) IncCommit() {
	m.CommandsCommitted.Add(1)
}

func (m *Metrics) IncApply() {
	m.ApplyCount.Add(1)
}

func (m *Metrics) ObserveLatency(d time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.latencies = append(m.latencies, d)
	if d > m.maxLatency {
		m.maxLatency = d
	}
}

func (m *Metrics) Snapshot() map[string]interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()
	var total time.Duration
	for _, d := range m.latencies {
		total += d
	}
	avg := time.Duration(0)
	if len(m.latencies) > 0 {
		avg = total / time.Duration(len(m.latencies))
	}
	return map[string]interface{}{
		"leader_elections":    m.LeaderElections.Load(),
		"commands_proposed":   m.CommandsProposed.Load(),
		"commands_committed":  m.CommandsCommitted.Load(),
		"apply_count":         m.ApplyCount.Load(),
		"avg_latency_ms":      avg.Milliseconds(),
		"max_latency_ms":      m.maxLatency.Milliseconds(),
		"samples":             len(m.latencies),
	}
}

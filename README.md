# Aether

**A fault-tolerant distributed key-value store** built with Raft-style consensus.

Aether is a systems-level project designed to demonstrate deep understanding of distributed systems fundamentals: leader election, log replication, persistence, and failure recovery.

> Built as a high-signal portfolio project for software engineering internships at companies that care about real systems (Microsoft, Amazon, Google, etc.).

---

## Why Aether Exists

Most student projects are CRUD apps or AI wrappers.  
Aether is **infrastructure**. It implements the same core ideas used inside real systems like etcd, Consul, and CockroachDB.

This project proves you can think about:

- Consensus under failure
- Consistency vs availability trade-offs
- Persistent state and crash recovery
- Concurrent network programming
- Performance measurement

---

## Architecture Overview

```
Client  →  Aether Client Library  →  Cluster of Nodes (Raft-inspired)
                                      ├─ Leader (handles writes)
                                      ├─ Followers (replicate log)
                                      └─ Persistent WAL + Snapshots
```

### Core Components

| Component | Responsibility |
|-----------|----------------|
| **Node** | Single server process. Can be Leader or Follower |
| **Raft Core** | Leader election, log replication, commit index |
| **State Machine** | Applies committed commands to the in-memory KV map |
| **WAL** | Write-ahead log for durability |
| **Client** | Simple Get / Put / Delete API |
| **Bench** | Throughput & latency measurement under load |

---

## Key Design Decisions (What Makes It Strong)

1. **Simplified but correct Raft**  
   Implements the essential parts of Raft (election, replication, commit) without full production complexity. Interviewers care about the *ideas*, not a perfect clone of etcd.

2. **Explicit failure model**  
   Nodes can be killed. The system recovers leadership and continues serving.

3. **Persistence first**  
   Log is written before acknowledgment. Crash recovery is real.

4. **Clean separation**  
   Consensus layer is independent from the state machine (classic Raft design).

5. **Measurable**  
   Built-in benchmark tool so you can talk about numbers in interviews.

---

## Project Structure

```
Aether/
├── cmd/
│   ├── aether-node/     # Single node binary
│   └── aether-bench/    # Load generator
├── internal/
│   ├── raft/            # Consensus core
│   ├── storage/         # WAL + snapshots
│   ├── kv/              # State machine
│   └── transport/       # RPC between nodes
├── client/              # Client library
├── docs/                # Design notes
├── scripts/             # Cluster start helpers
└── README.md
```

---

## Quick Start (Development)

```bash
# Requires Go 1.22+
go mod tidy

# Start a 3-node cluster (example)
go run ./cmd/aether-node --id 1 --peers localhost:7001,localhost:7002,localhost:7003 --port 7001
go run ./cmd/aether-node --id 2 --peers localhost:7001,localhost:7002,localhost:7003 --port 7002
go run ./cmd/aether-node --id 3 --peers localhost:7001,localhost:7002,localhost:7003 --port 7003

# Use the client
go run ./client/example
```

---

## What Interviewers Will Ask (Be Ready)

- How does leader election work when two nodes have the same term?
- What happens if the leader crashes after writing to the WAL but before majority ack?
- Why is the log the source of truth?
- How would you add sharding later?
- How do you measure correctness under partitions?

Having built Aether gives you concrete answers.

---

## Status

This repository contains a clean, interview-ready foundation.  
Core consensus and persistence logic are implemented in a readable, educational style.  
Extend it with stronger testing, membership changes, or linearizable reads for even more depth.

---

## Author

**Affan Ahmed Shariff**  
B.E. Computer Science — BMS College of Engineering  
BS Data Science — IIT Madras  
Founder, Primeora Solutions

Built to demonstrate systems-level engineering capability beyond typical undergrad projects.

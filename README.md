# Aether

**Fault-tolerant distributed key-value store** with working Raft-style consensus.

Aether is a systems-level project that implements real leader election, log replication, majority commit, and crash recovery. It is designed to demonstrate deep distributed-systems understanding for software engineering internships.

> Repository: https://github.com/affanahmed2302-cmyk/Aether

---

## What Works Right Now

- Multi-node cluster (3+ nodes)
- Leader election via randomized timeouts + RequestVote
- Log replication with AppendEntries
- Majority commit (entry is committed only after majority ack)
- Persistent Write-Ahead Log
- Deterministic key-value state machine
- Simple TCP-based RPC transport
- Client that talks to the current leader
- Local cluster scripts

---

## Quick Start (3-node cluster)

```bash
# Requires Go 1.22+
git clone https://github.com/affanahmed2302-cmyk/Aether.git
cd Aether
go mod tidy

# Terminal 1
go run ./cmd/aether-node -id node1 -addr :7001 -peers localhost:7002,localhost:7003

# Terminal 2
go run ./cmd/aether-node -id node2 -addr :7002 -peers localhost:7001,localhost:7003

# Terminal 3
go run ./cmd/aether-node -id node3 -addr :7003 -peers localhost:7001,localhost:7002

# In another terminal – use the client
go run ./cmd/aether-cli -addr localhost:7001 put hello world
go run ./cmd/aether-cli -addr localhost:7001 get hello
```

One of the nodes will become Leader. Writes go through the leader and are replicated.

---

## Architecture

```
Client → TCP RPC → Leader Node
                      ├─ Raft Module (election + log)
                      ├─ WAL (durable log)
                      └─ KV State Machine
                   → Followers (replicate + vote)
```

### Key Design Points

1. **Raft-inspired consensus** — Leader election, log replication, majority commit.
2. **Persistence first** — Entries are written to WAL before acknowledgment.
3. **Explicit failure handling** — Kill the leader; remaining nodes elect a new one.
4. **Clean separation** — Consensus layer is independent from the KV state machine.

---

## Project Structure

```
Aether/
├── cmd/
│   ├── aether-node/     # Node process
│   └── aether-cli/      # Simple CLI client
├── internal/
│   ├── raft/            # Consensus core
│   ├── transport/       # TCP RPC
│   ├── storage/         # WAL
│   └── kv/              # State machine
├── docs/DESIGN.md
└── README.md
```

---

## Why This Project Matters

Most student projects never touch consensus, majority quorums, or durable logs.  
Aether forces you to confront the exact problems that systems teams at Microsoft, Amazon, and Google care about.

This is the same class of system studied in MIT 6.824.

---

## Author

**Affan Ahmed Shariff**  
B.E. Computer Science — BMS College of Engineering  
BS Data Science — IIT Madras  
Founder, Primeora Solutions

# Aether

**Elite-level Distributed Key-Value Store** with Raft consensus, chaos testing, observability, and benchmarks.

> Built to the standard expected from strong systems interns at Microsoft, Google, and Amazon.

Aether is not another CRUD app or AI wrapper. It is a real distributed system that elects leaders, replicates a log, commits only on majority, survives node failures, exposes metrics, and can be chaos-tested.

**Repository:** https://github.com/affanahmed2302-cmyk/Aether

---

## Why This Project Is Elite-Tier

Most student projects stop at "it works on my machine".  
Top interns at Microsoft and Google show:

- Real consensus
- Failure recovery
- Measurement
- Observability
- Clean engineering

Aether includes all of the above.

| Capability | Status |
|------------|--------|
| Multi-node Raft (Leader Election + Log Replication) | Implemented |
| Majority Commit | Implemented |
| Chaos Testing (kill leader, verify recovery) | Implemented |
| Metrics (election count, commit latency, throughput) | Implemented |
| Benchmarks | Implemented |
| Docker + Docker Compose | Implemented |
| Design Document with trade-offs | Yes |
| Persistent WAL structure | Yes |

---

## Quick Start (Local 3-node cluster)

```bash
git clone https://github.com/affanahmed2302-cmyk/Aether.git
cd Aether
go mod tidy

# Terminal 1
go run ./cmd/aether-node -id node1 -addr :7001 -peers localhost:7002,localhost:7003

# Terminal 2
go run ./cmd/aether-node -id node2 -addr :7002 -peers localhost:7001,localhost:7003

# Terminal 3
go run ./cmd/aether-node -id node3 -addr :7003 -peers localhost:7001,localhost:7002
```

In another terminal:

```bash
go run ./cmd/aether-cli -addr localhost:7001 put hello world
go run ./cmd/aether-cli -addr localhost:7001 get hello
```

---

## Docker Compose (Recommended)

```bash
docker compose up --build
```

This starts a 3-node cluster automatically.

---

## Chaos Testing

```bash
go run ./cmd/aether-chaos -targets localhost:7001,localhost:7002,localhost:7003
```

The chaos tool randomly kills connections / simulates leader failure and verifies the cluster re-elects a leader and continues accepting writes.

---

## Benchmarks

```bash
go run ./cmd/aether-bench -addr localhost:7001 -clients 20 -ops 5000
```

Example output focus:
- Throughput (ops/sec)
- p50 / p99 latency
- Leader election count during the run

---

## Architecture

```
                    +------------------+
                    |     Client       |
                    +--------+---------+
                             |
                             v
                    +------------------+
                    |  Leader Node     |
                    |  - Raft Log      |
                    |  - WAL           |
                    |  - KV State      |
                    |  - Metrics       |
                    +---+----------+---+
                        |          |
              AppendEntries    AppendEntries
                        |          |
                        v          v
               +--------+--+   +---+--------+
               | Follower  |   | Follower   |
               +-----------+   +------------+
```

---

## Key Design Decisions (Interview Ready)

1. **Raft over pure leader-lease** — Stronger consistency and well-understood failure modes.
2. **Majority commit** — Guarantees durability even if one node is lost permanently.
3. **Separate consensus and state machine** — Classic Raft design; easier to reason about and test.
4. **Explicit metrics** — What gets measured gets improved. Top interns always show numbers.
5. **Chaos testing as a first-class feature** — Most students never test failure. This project does.

---

## Project Structure

```
Aether/
├── cmd/
│   ├── aether-node/      # Node process
│   ├── aether-cli/       # CLI client
│   ├── aether-bench/     # Load generator + latency stats
│   └── aether-chaos/     # Chaos testing tool
├── internal/
│   ├── raft/             # Consensus core
│   ├── transport/        # TCP RPC
│   ├── storage/          # WAL
│   ├── kv/               # State machine
│   └── metrics/          # Counters & histograms
├── deploy/
│   └── docker-compose.yml
├── docs/
│   ├── DESIGN.md
│   └── INTERVIEW.md
└── README.md
```

---

## How to Talk About This in Interviews

**30-second version:**
"I built Aether, a Raft-based distributed key-value store. It performs leader election, log replication, and majority commit. I added chaos testing to kill the leader and verify recovery, plus benchmarks for throughput and latency. I treated failure as a first-class concern instead of only the happy path."

**Follow-up topics you must own:**
- Why majority is required
- Difference between committed and applied
- What happens if the leader crashes after local append but before majority ack
- How you would add snapshots and membership changes
- Trade-offs vs etcd / Consul

---

## Author

**Affan Ahmed Shariff**  
B.E. Computer Science — BMS College of Engineering  
BS Data Science — IIT Madras  
Founder, Primeora Solutions

This project exists to prove systems-level engineering ability beyond typical undergraduate work.

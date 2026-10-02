# Aether Design Document (Elite Version)

## Goal
Build a distributed key-value store that demonstrates the same concerns top systems interns are evaluated on:
- Consensus under partial failure
- Durability
- Observability
- Performance measurement
- Operational realism (chaos)

## Consensus
Raft-inspired:
- Leader election with randomized timeouts
- AppendEntries for replication and heartbeats
- Majority commit
- Deterministic state machine

## Observability
Metrics package tracks:
- Leader elections
- Proposed vs committed commands
- Apply count
- Latency samples

## Chaos
`aether-chaos` continuously writes while contacting random nodes. The goal is to show the system keeps making progress even when individual nodes are unreliable.

## Benchmarks
`aether-bench` reports throughput and latency under concurrent clients. Numbers matter in systems interviews.

## Non-Goals (Explicit)
- Full production Raft (joint consensus, PreVote, batched fsync, etc.)
- Multi-raft / sharding
- Disk-optimized storage engine

These are listed so interviewers see mature scoping.

## Future Work
1. Snapshots + log truncation
2. Membership changes
3. Linearizable read path
4. Persistent WAL replay on restart
5. Prometheus metrics endpoint

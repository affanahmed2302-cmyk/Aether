# Aether Design Notes

## Goals

- Demonstrate real distributed-systems competence
- Stay small enough that a strong 1st-year student can understand and extend every line
- Be interview-ready: every major design decision has a clear rationale

## Consensus

We use a **simplified Raft**:

- Leader election via randomized timeouts
- Log replication (AppendEntries)
- Commit index advanced only after majority acknowledgment (full version)
- State machine is strictly sequential and deterministic

We intentionally omit:

- Membership changes (AddServer / RemoveServer)
- Log compaction / snapshotting shipping (structure is ready)
- Linearizable reads via ReadIndex (can be added later)

These omissions keep the codebase readable while still covering the hardest conceptual parts.

## Durability

Every committed entry is written to a Write-Ahead Log (WAL) *before* the client receives success.  
On restart the node replays the WAL to reconstruct state.

## Failure Model

- Crash-stop failures (nodes may die and restart)
- Network partitions are possible; safety is prioritized over availability (CP system)

## Why This Project Signals Strongly

Most undergrad projects never touch:

- Consensus algorithms
- Persistent logs
- Leader election under concurrency
- Explicit consistency guarantees

Aether forces you to confront all of them. That is exactly the signal systems teams at Microsoft, Amazon, and Google look for.

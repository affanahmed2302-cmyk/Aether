# Aether Design

## Goals
- Real multi-node leader election and log replication
- Majority commit
- Readable educational codebase that still demonstrates production ideas
- Strong signal for systems interviews

## Consensus
Simplified Raft:
- Randomized election timeouts
- RequestVote RPC
- AppendEntries RPC (heartbeats + log entries)
- Commit index advanced only after majority acknowledgment
- State machine applied in log order

## Transport
TCP + gob encoding for RPCs. Simple and dependency-free.

## Persistence
WAL structure is present; full crash-recovery replay can be extended easily.

## Failure Handling
Kill the current leader process. Remaining nodes will time out and elect a new leader.

## What Interviewers Care About
- Why majority is required
- What happens if the leader crashes after local append but before majority
- Difference between committed and applied
- How you would add snapshots / membership changes later

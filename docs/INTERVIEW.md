# How to Present Aether in Microsoft / Google Interviews

## 30-second pitch
I built Aether, a Raft-based distributed key-value store. It performs leader election, log replication, and majority commit. I added chaos testing to verify the system recovers when the leader fails, and benchmarks that report throughput and latency. I focused on failure modes and measurement, not only the happy path.

## Strong follow-up answers

**Q: Why Raft?**  
A: Raft gives strong leadership and a clean decomposition (election, log replication, safety). It is easier to implement correctly than Multi-Paxos while teaching the same core ideas used in etcd and Consul.

**Q: What happens if the leader crashes after appending locally but before majority ack?**  
A: The entry is not committed. A new leader will be elected. Because the entry never reached majority, it is safe to discard or for the new leader to overwrite according to Raft log matching rules.

**Q: Difference between committed and applied?**  
A: Committed means the entry is durable on a majority and will never be lost. Applied means the state machine has executed it. Applied ≤ Committed always.

**Q: How would you improve it next?**  
A: Log compaction via snapshots, membership changes, linearizable reads (ReadIndex/LeaseRead), better quorum tracking, and disk-backed WAL with fsync batching.

## Signals this project gives
- You understand consensus
- You test failure
- You measure performance
- You can explain trade-offs
- You engineer beyond tutorials

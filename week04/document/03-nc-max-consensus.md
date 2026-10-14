# NC-Max Consensus

> **Curriculum Reference:** [truthixify/learn-ckb-in-45-minutes / 03-nc-max-consensus.md](https://github.com/truthixify/learn-ckb-in-45-minutes/blob/main/notes/03-nc-max-consensus.md)

## The Problem with Classic Nakamoto Consensus
- Two miners find blocks simultaneously → one becomes an **orphan** (wasted work)
- Faster block rate → more orphans → less security
- Bitcoin compensates by keeping blocks slow (~10 min)

## NC-Max: Two-Step Transaction Confirmation

NC-Max splits transaction inclusion into **propose** and **commit**:

### Step 1: Propose
- Transaction appears in a block's **proposal zone** as a short ID (just a few bytes)
- Signals to the network: "this transaction exists, start validating it"

### Step 2: Commit
- After a waiting window (2-4 blocks), the full transaction is **committed** in a later block
- By this time, all nodes have already validated and cached the transaction

### Why This Works
- Block validation bottleneck = verifying transactions inside it
- Proposal phase gives nodes time to pre-validate and cache
- When the commitment block arrives, validation is nearly instant
- Faster block propagation → fewer orphans → safely increase block rate

## Difficulty Adjustment
- NC-Max accounts for orphan blocks in difficulty calculation
- Targets a specific **orphan rate** rather than just block time
- Network self-tunes for maximum throughput without sacrificing security

## Why the Waiting Window is Mandatory
- If propose + commit happen in the same block, nodes receive transactions they've never seen
- They must validate from scratch before accepting → slow propagation → more orphans
- The window gives the network time to pre-fetch and pre-validate
- By commit time, nodes already have the tx cached → block propagates almost instantly

## Key Takeaways
- NC-Max solves the throughput vs security trade-off of classic Nakamoto Consensus
- Propose/commit separation lets nodes pre-validate, speeding up block propagation
- Orphan-aware difficulty adjustment allows the chain to safely push for higher throughput
- CKB achieves faster block times than Bitcoin while maintaining PoW security guarantees

---

## Questions

**Q1:** In the propose/commit model, why is there a mandatory waiting window between proposing and committing a transaction? What would go wrong if a transaction could be proposed and committed in the same block?

<details>
<summary>Hint</summary>
Think about what the receiving nodes need to do when they get a block with transactions they've never seen before.
</details>

<details>
<summary>Answer</summary>
If propose and commit happen in the same block, nodes receive transactions they've never seen before. They must download, validate, and execute those transactions from scratch before accepting the block. That validation takes time — during which other miners keep mining on the old tip, producing competing blocks (orphans). The waiting window gives nodes time to pre-fetch and pre-validate, so when the commitment block arrives, it propagates almost instantly.
</details>

---

**Q2:** Why does NC-Max target an orphan rate in its difficulty adjustment rather than just a block time?

<details>
<summary>Answer</summary>
Block time alone doesn't capture network health. If the network improves propagation speed (better connectivity, smaller blocks), it could safely handle faster blocks without more orphans. By targeting an orphan rate, the network self-tunes — when propagation is fast, difficulty adjusts to allow faster blocks; when orphans rise, it slows down. This maximizes throughput for the network's actual conditions.
</details>

---

**Q3:** Bitcoin keeps blocks at ~10 minutes to minimize orphans. Why can't Bitcoin simply adopt the propose/commit approach to speed up blocks?

<details>
<summary>Answer</summary>
It could in theory, but it would require a hard fork — a fundamental change to Bitcoin's consensus rules. Bitcoin's conservative governance makes protocol-level changes extremely difficult to push through. CKB had the advantage of designing NC-Max from the start, rather than retrofitting it onto an existing network.
</details>

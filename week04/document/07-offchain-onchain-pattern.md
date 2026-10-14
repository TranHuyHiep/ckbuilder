# Off-Chain Computation, On-Chain Verification

> **Curriculum Reference:** [truthixify/learn-ckb-in-45-minutes / 07-offchain-onchain-pattern.md](https://github.com/truthixify/learn-ckb-in-45-minutes/blob/main/notes/07-offchain-onchain-pattern.md)

## The Core Pattern

| | Ethereum (Account Model) | CKB (Cell Model) |
|---|---|---|
| **Who computes new state?** | The contract (on-chain) | The transaction builder (off-chain) |
| **What runs on-chain?** | Full computation | Pure validation only |
| **Script/contract role** | Compute + store + validate | Validate only (return 0 or non-zero) |

- Off-chain: wallet/dApp figures out the new cells, balances, and state transitions
- On-chain: scripts verify that the proposed state transition is valid
- Verification is almost always cheaper than computation → more efficient

## What You Can Build with RISC-V Scripts

Since CKB-VM runs real RISC-V, practically any program can be a validation script:
- Custom cryptographic schemes (BLS, Schnorr, RSA)
- Zero-knowledge proof verification
- Bitcoin/Ethereum signature verification (cross-chain interop)
- Spawning child scripts from parent scripts (CKB-VM v2)
- Only constraint: must finish within the cycle limit

## Security Advantage: No Re-Entrancy

- Ethereum: contract A calls B which calls back into A → re-entrancy attacks (e.g., the DAO hack)
- CKB: scripts can't call each other or modify state → entire vulnerability class eliminated

## Key Takeaways
- CKB separates state computation (off-chain) from state validation (on-chain)
- Transaction builders propose state, scripts approve or reject
- RISC-V means near-unlimited expressiveness for validation logic
- Pure validation model eliminates re-entrancy and other inter-contract attack vectors

---

## Questions

**Q1:** In the off-chain computation model, who is responsible for calculating the correct output cells in a token transfer? What happens if they calculate wrong?

<details>
<summary>Answer</summary>
The **transaction builder** (wallet or dApp) is responsible. They figure out the new cells, balances, and state transitions off-chain. If they calculate wrong (e.g., outputs don't sum correctly), the on-chain type script will detect the invalid state transition and reject the transaction. The script is the safety net — it doesn't need to know *how* to compute the right answer, just how to verify one.
</details>

---

**Q2:** Ethereum had the DAO hack, which exploited re-entrancy. Explain why this specific attack is impossible on CKB.

<details>
<summary>Answer</summary>
Re-entrancy happens when contract A calls contract B, and B calls back into A before A finishes executing — exploiting A's inconsistent intermediate state. On CKB, scripts cannot call other scripts, cannot modify state, and cannot trigger external execution. Each script runs independently, reads transaction data via syscalls, and returns a yes/no verdict. There's no execution flow between scripts to exploit.
</details>

---

**Q3:** A friend says "CKB scripts are less powerful than Ethereum contracts because they can't compute state." How would you respond?

<details>
<summary>Hint</summary>
Think about where the computation actually happens and whether the end result is any different.
</details>

<details>
<summary>Answer</summary>
The computation still happens — it's just done off-chain by the transaction builder instead of on-chain by the script. The end result is the same: valid state transitions. But CKB's approach is actually more efficient (verification < computation), more secure (no re-entrancy), and more flexible (scripts can run any RISC-V program, not just EVM bytecode). The "limitation" is really a design choice that shifts complexity to the right place.
</details>

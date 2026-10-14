# CKB-VM and Script Execution

> **Curriculum Reference:** [truthixify/learn-ckb-in-45-minutes / 06-ckb-vm-execution.md](https://github.com/truthixify/learn-ckb-in-45-minutes/blob/main/notes/06-ckb-vm-execution.md)

## CKB-VM vs EVM
- EVM: custom instruction set designed for Ethereum, requires Solidity/Vyper
- CKB-VM: full **RISC-V IMC** implementation — a real hardware CPU architecture
- Any language with a RISC-V backend works: C, Rust, Go, JavaScript (via QuickJS)
- Standard toolchains, debuggers, and testing tools apply

## Syscalls — How Scripts Read Transaction Data

Scripts are sandboxed — no network, filesystem, or outside access.
They read transaction data through syscalls:

| Syscall | What it does |
|---|---|
| `ckb_load_tx_hash` | Get current transaction hash |
| `ckb_load_script` | Get the executing script's own code_hash, hash_type, args |
| `ckb_load_cell` | Load a cell from inputs, outputs, or cell_deps |
| `ckb_load_cell_data` | Load the data field of a cell |
| `ckb_load_input` | Load an input (previous outpoint + since) |
| `ckb_load_witness` | Load a witness (signatures, proofs) |
| `ckb_load_header` | Load a block header from header_deps |

## The Execution Model

A script is a program that:
1. Starts in CKB-VM
2. Uses syscalls to read the transaction
3. Performs validation logic
4. **Returns 0** → validation passed (success)
5. **Returns non-zero** → transaction rejected (failure)

No storage writes. No state modifications. No calling other contracts.
Scripts simply say **yes or no** to the transaction.

## Cycles — CKB's "Gas"

- Each RISC-V instruction costs a number of **cycles**
- Transactions have a **cycle limit** (consensus rule)
- Prevents infinite loops, ensures termination
- Unlike Ethereum gas: you don't pay per cycle
- Cost = transaction fee (to miners), not per-cycle billing
- Transactions exceeding the cycle limit are simply invalid

## Key Takeaways
- CKB-VM is a real RISC-V CPU implementation, not a blockchain-specific VM
- Write scripts in any language that compiles to RISC-V
- Scripts are pure validation: read transaction via syscalls, return 0 or non-zero
- No side effects — scripts don't write state, they approve or reject
- Cycles limit computation, but aren't billed per-unit like ETH gas

---

## Questions

**Q1:** CKB scripts can't store data, modify state, or call other contracts — they just return 0 or non-zero. Does this feel like a limitation? Can you think of why this "pure validation" model might actually be an advantage?

<details>
<summary>Hint</summary>
Think about what happens on Ethereum when contract A calls contract B which calls back into contract A...
</details>

<details>
<summary>Answer</summary>
It's actually a superpower. Since the transaction builder proposes all state changes off-chain and scripts only validate, verification is cheaper than computation. And because scripts can't call each other or modify state, entire vulnerability classes are eliminated — most notably re-entrancy attacks (like the famous DAO hack on Ethereum). The pure validation model is simpler, more secure, and more efficient.
</details>

---

**Q2:** Why can CKB scripts be written in any programming language while Ethereum smart contracts are mostly limited to Solidity?

<details>
<summary>Answer</summary>
The EVM uses a custom instruction set designed specifically for Ethereum — you need a compiler that targets this custom architecture (hence Solidity/Vyper). CKB-VM implements RISC-V, a real hardware CPU architecture with mature compiler support. Any language with a RISC-V backend (C, Rust, Go, and even JavaScript via QuickJS compiled to RISC-V) can produce CKB scripts.
</details>

---

**Q3:** If cycles aren't billed per-unit like Ethereum gas, what prevents someone from spamming the network with computationally expensive transactions?

<details>
<summary>Answer</summary>
Two things: (1) there's a per-transaction cycle limit — transactions exceeding it are simply invalid and rejected. (2) Transaction fees (the capacity difference between inputs and outputs) still apply. Miners prioritize transactions with higher fees, so computationally expensive transactions need to pay proportionally to get included. The fee market still provides economic protection against spam.
</details>

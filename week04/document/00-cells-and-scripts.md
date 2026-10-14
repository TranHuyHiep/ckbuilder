# Cells and Scripts — The Foundation of CKB

> **Curriculum Reference:** [truthixify/learn-ckb-in-45-minutes / 00-cells-and-scripts.md](https://github.com/truthixify/learn-ckb-in-45-minutes/blob/main/notes/00-cells-and-scripts.md)

## What is CKB?

CKB (Common Knowledge Base) is the Layer 1 blockchain of the **Nervos Network**. It's a proof-of-work chain designed to be a secure, decentralized store of value and a foundation for Layer 2 scaling.

Unlike Ethereum's account model, CKB uses a **Cell Model** — a generalized version of Bitcoin's UTXO model.

## What is a Cell?

A cell is the basic unit of data on CKB. Think of it as a box that stores data on the blockchain. Every piece of state on CKB lives inside a cell.

A cell has four fields:

```rust
Cell {
    capacity:    u64,      # size limit in bytes (also the CKB token amount)
    lock:        Script,   # who can spend/consume this cell
    type:        Script?,  # optional — rules for how the cell can be used
    data:        Bytes     # arbitrary data stored in the cell
}
```

### capacity

This is both the **size limit** and the **CKB token amount**. 1 CKByte = 1 byte of storage. A cell's total size (all fields combined) cannot exceed its capacity value. So if you have a cell with `capacity: 100`, the cell's total serialized size must be ≤ 100 bytes.

### lock (Lock Script)

Every cell **must** have a lock script. It determines **who can spend this cell**. When a transaction tries to consume this cell as an input, the lock script runs and checks: "Is this authorized?"

The most common lock script checks for a valid digital signature from the cell owner.

### type (Type Script)

Optional. If present, it defines **rules for the cell's usage**. Type scripts validate state transitions — they can enforce things like "this cell represents a token, and you can't create tokens out of thin air."

### data

Arbitrary bytes. Can store anything: token balances, NFT metadata, script code for other cells to reference, or nothing at all.

## What is a Script?

A script is a program that runs in CKB-VM (a RISC-V virtual machine). It's defined by three fields:

```rust
Script {
    code_hash: H256,       # identifies which program to run
    hash_type: String,     # how to interpret code_hash
    args: Bytes            # arguments passed to the program
}
```

- **code_hash** + **hash_type** together locate the actual executable binary
- **args** customizes the script for a specific use (e.g., the owner's public key hash)

Scripts don't "live" inside the cell — they're just references. The actual executable code is deployed in a **separate cell** and looked up at runtime.

## Live Cells vs Dead Cells

- A **live cell** exists on-chain and hasn't been consumed. It's available as state.
- A **dead cell** has been consumed by a transaction. It no longer exists as usable state.

There's no "modify." To update a cell, you **consume** it (kill it) and **create** a new one with the updated data. This is why CKB uses the term "cell" — they're born (created as outputs) and die (consumed as inputs).

## How It All Connects

```text
┌─────────────────────────────────────────────────────┐
│                      Cell                           │
│                                                     │
│  capacity: 200 CKB                                  │
│                                                     │
│  lock: Script {                                     │
│      code_hash: 0x9bd7...  ← identifies the program │
│      hash_type: "type"                              │
│      args: 0xabcd...       ← owner's pubkey hash    │
│  }                                                  │
│                                                     │
│  type: null  (no type script for this cell)         │
│                                                     │
│  data: 0x   (empty)                                 │
└────────────────────────────��────────────────────────┘
```

This cell:
- Holds 200 CKB
- Can only be spent by whoever owns the private key matching `0xabcd...`
- Has no additional validation rules (no type script)
- Stores no extra data

## Key Takeaways
- Cells are the fundamental unit of state on CKB — everything lives in cells
- Each cell has capacity (size/value), lock (who can spend), optional type (validation rules), and data
- 1 CKByte = 1 byte of on-chain storage
- Scripts are references to programs, not the programs themselves
- Cells are immutable — to change state, destroy old cells and create new ones

---

## Questions

**Q1:** A cell has `capacity: 100`. Can you store 100 bytes of data in the `data` field?

<details>
<summary>Hint</summary>
Think about what else takes up space in the cell besides the data field.
</details>

<details>
<summary>Answer</summary>
No. The capacity covers the **total** serialized size of the cell, including the capacity field itself (8 bytes), the lock script fields, the optional type script, AND the data. So with a capacity of 100, you'd have significantly less than 100 bytes available for data after accounting for the other fields.
</details>

---

**Q2:** Why must every cell have a lock script but the type script is optional?

<details>
<summary>Answer</summary>
The lock script answers "who can spend this cell?" — without it, anyone could consume the cell and steal the CKB inside. It's the fundamental ownership mechanism. The type script answers "is this state transition valid?" — which is only needed when the cell represents something with rules (tokens, NFTs, etc.). A simple CKB transfer with no special logic doesn't need a type script.
</details>

---

**Q3:** If cells are immutable and can only be created or destroyed, how would you update a user's token balance from 500 to 300?

<details>
<summary>Hint</summary>
Think about what a transaction does with inputs and outputs.
</details>

<details>
<summary>Answer</summary>
You'd create a transaction that **consumes** the old cell (with balance 500) as an input and **creates** a new cell (with balance 300) as an output. The old cell dies, the new cell is born. There is no in-place modification — it's always destroy and recreate.
</details>

---

**Q4:** Script code doesn't live inside the cell's lock or type fields — it's stored in a separate cell and referenced. Why do you think CKB was designed this way instead of embedding code directly?

<details>
<summary>Answer</summary>
If every cell embedded its own script code, the same program would be duplicated millions of times across cells. By storing code in one cell and referencing it by hash, CKB achieves massive deduplication — one deployment serves all cells using that script. This saves on-chain storage (which directly costs CKB tokens) and makes the system more efficient.
</details>

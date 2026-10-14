# CKB Transaction Structure

> **Curriculum Reference:** [truthixify/learn-ckb-in-45-minutes / 01-transaction-structure.md](https://github.com/truthixify/learn-ckb-in-45-minutes/blob/main/notes/01-transaction-structure.md)

## Core Concept
A CKB transaction **consumes** existing cells (inputs) and **creates** new cells (outputs).
There is no "modify" — it's always destroy-and-recreate (like Bitcoin's UTXO model).

## Transaction Fields

| Field | Purpose |
|---|---|
| `cell_deps` | Cells referenced as read-only dependencies (e.g., deployed script code) |
| `header_deps` | Block headers the transaction depends on |
| `inputs` | Cells being consumed (destroyed) |
| `outputs` | New cells being created |
| `outputs_data` | Data field for each output cell |
| `witnesses` | Proofs/signatures (not stored on-chain in cells) |

## cell_deps — The Key Mechanism

- Scripts are referenced by **code hash**, but actual code lives in other cells
- `cell_deps` tells the transaction where to find the script code
- Script code is **deployed once** as data in a cell, then referenced by many transactions
- Upgradeable: deploy new code cell, point to it

## Validation Flow

1. CKB-VM loads the **lock script** of each input cell
2. Finds actual code via `cell_deps`
3. Executes the lock script — returns `0` means authorized to spend
4. **Type scripts** of both inputs AND outputs run (if present) to validate state transitions

## Lock Scripts vs Type Scripts

| | Lock Script | Type Script |
|---|---|---|
| **Runs on** | Inputs only | Both inputs AND outputs |
| **Question it answers** | "Is this tx authorized to consume this cell?" | "Is this state transition valid?" |
| **Example** | Signature verification | Token supply conservation |

### Why type scripts see both sides
Because they can **compare** inputs vs outputs, they enforce:
- **Conservation rules** — e.g., total UDT in inputs = total UDT in outputs
- **Uniqueness** — e.g., only one NFT cell with a given ID in outputs
- **State machines** — e.g., a cell's status can only go from "pending" → "approved"

## Key Takeaways
- Transactions are stateless transformations: old cells in, new cells out
- Code is separated from references to code (cell_deps pattern)
- Lock scripts = authorization ("can you spend this?")
- Type scripts = validation ("is this state transition valid?")

---

## Questions

**Q1:** When a type script runs during transaction validation, it has access to both input and output cells. Why is this important? What can a type script enforce that a lock script cannot?

<details>
<summary>Hint</summary>
Think about what "seeing both sides" enables. What kinds of rules require comparing old state to new state?
</details>

<details>
<summary>Answer</summary>
Because type scripts see both inputs and outputs, they can **compare** old state vs new state. This lets them enforce conservation rules (e.g., total tokens in = total tokens out), uniqueness constraints (only one NFT with a given ID), and state machine transitions (status can only go from "pending" to "approved"). Lock scripts only see inputs, so they can only answer "is this spend authorized?" — they can't validate what the new state looks like.
</details>

---

**Q2:** If a transaction has 3 inputs totaling 300 CKB and 2 outputs totaling 280 CKB, where do the remaining 20 CKB go?

<details>
<summary>Answer</summary>
The difference (inputs - outputs) becomes the **transaction fee** paid to the miner who includes the transaction in a block. This is similar to Bitcoin's fee model — there's no explicit "gas price" field; the fee is implicit in the difference between input and output capacity.
</details>

---

**Q3:** Why does CKB separate script code from the cells that reference it (via `cell_deps`) instead of embedding code directly in each cell?

<details>
<summary>Answer</summary>
If every cell embedded its own script code, the same lock script binary would be duplicated across millions of cells. With `cell_deps`, the code is deployed **once** in a single cell, and every transaction that needs it just references that cell. This saves massive amounts of on-chain storage — and since 1 CKByte = 1 byte, that directly saves real cost.
</details>

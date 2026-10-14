# Witnesses — Proof Data in Transactions

> **Curriculum Reference:** [truthixify/learn-ckb-in-45-minutes / 10-witness-layout.md](https://github.com/truthixify/learn-ckb-in-45-minutes/blob/main/notes/10-witness-layout.md)

## What Are Witnesses?

Witnesses are the **proof** section of a transaction. They carry evidence that scripts need to verify authorization or validity.

Most commonly: a **digital signature** proving the spender owns the cell.

## Why Separate From Inputs/Outputs?

- Inputs = which cells to consume
- Outputs = what cells to create
- Witnesses = proof that this is all authorized
- Signatures aren't "state" — they don't need to persist on-chain after validation
- Keeping proofs separate keeps the data model clean

## Not Always a Signature

Since CKB-VM can run any program, witnesses can carry anything:
- Digital signatures (most common)
- Merkle proofs
- Zero-knowledge proofs
- Any arbitrary data the script expects

The script decides what to expect — witnesses are just the delivery mechanism.

## The Common Flow (Signature Example)

1. Alice builds a transaction spending her cell
2. She signs the transaction hash with her private key
3. Puts the signature in the `witnesses` field
4. Lock script runs → reads witness → extracts signature
5. Checks signature against Alice's public key (from lock script `args`)
6. Valid → return 0 → authorized

## Witness-to-Input Matching

Matched by **position** (index):

```text
inputs:    [ input_0,   input_1,   input_2 ]
witnesses: [ witness_0, witness_1, witness_2 ]
               ↕            ↕           ↕
           proof for    proof for   proof for
           input_0      input_1     input_2
```

- There can be **more witnesses than inputs** — extras carry additional data
- Scripts can read any witness, not just the one at their index

## WitnessArgs — Internal Structure

Each witness is typically structured as:

```json
WitnessArgs {
    lock:        Bytes,   ← proof for the lock script (e.g., signature)
    input_type:  Bytes,   ← proof for the input's type script
    output_type: Bytes,   ← proof for the output's type script
}
```

- A cell can have both a lock and type script, each needing different proofs
- WitnessArgs gives each script its own slot instead of sharing one blob

## Witnesses + Script Groups

The script reads the witness at the **first input index** of its group.

```text
inputs:    [ Alice_cell_0, Alice_cell_1, Alice_cell_2, Bob_cell_3 ]
witnesses: [ Alice's sig,   empty,        empty,        Bob's sig  ]
               ↑                                           ↑
         Alice's group                                Bob's group
         first index = 0                              first index = 3
```

- Alice's lock reads `witness_0` → one signature covers inputs 0, 1, 2
- Bob's lock reads `witness_3` → one signature covers input 3
- `witness_1` and `witness_2` can be empty or carry type script proof data
- One signature per signer, no matter how many cells they spend

## Missing/Empty Witness
- If a witness is missing or empty, the lock script has no proof to verify
- Lock script returns non-zero → transaction is **rejected** (never hits the chain)
- CKB doesn't "revert" — invalid transactions are simply thrown out before inclusion

## Key Takeaways
- Witnesses = proof data (signatures, proofs, etc.) attached to a transaction
- They exist only in the transaction, not stored in cells
- Scripts read witnesses via `ckb_load_witness` syscall
- What a witness contains depends entirely on what the script expects

---

## Questions

**Q1:** What would happen if someone submitted a transaction with an empty or missing witness for an input?

<details>
<summary>Answer</summary>
The lock script would try to read the witness, find nothing (or empty data), and have no signature to verify. It returns a non-zero exit code, and the transaction is **rejected** — it never hits the chain. CKB doesn't "revert" like Ethereum; invalid transactions are simply thrown out before inclusion.
</details>

---

**Q2:** The WitnessArgs has three separate fields: `lock`, `input_type`, and `output_type`. Why is the type script's proof split into two fields rather than just one "type" field?

<details>
<summary>Hint</summary>
Can the input's type script and the output's type script be different scripts entirely?
</details>

<details>
<summary>Answer</summary>
Yes — the input cell's type script and the output cell's type script at the same index can be **completely different scripts**. Each needs its own proof data. `input_type` carries proof for the input's type script, `output_type` carries proof for the output's type script. Three scripts, three slots, no collisions.
</details>

---

**Q3:** Alice's one signature in `witness_0` covers all three of her inputs (in a script group). What exactly is being signed — and why is that sufficient?

<details>
<summary>Answer</summary>
Just the **transaction hash**. This is sufficient because the transaction hash already commits to all inputs, outputs, cell_deps, and witnesses (with the signature zeroed out). By signing the hash, Alice is implicitly approving the entire transaction. If anyone changed even one input or output, the hash would change and her signature would be invalid.
</details>

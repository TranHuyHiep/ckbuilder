# Dep Groups — Bundling Cell Dependencies

> **Curriculum Reference:** [truthixify/learn-ckb-in-45-minutes / 09-dep-groups.md](https://github.com/truthixify/learn-ckb-in-45-minutes/blob/main/notes/09-dep-groups.md)

## CellDep Structure

Each entry in `cell_deps` looks like:

```rust
CellDep {
    out_point: OutPoint,    # which cell to reference
    dep_type: "code" | "dep_group"
}
```

## dep_type: "code"

The simple case — the referenced cell's data is used directly as script code or data.

```json
cell_deps: [
    { out_point: "cell with secp256k1 binary", dep_type: "code" }
]
```

## dep_type: "dep_group"

A dep group is a **cell whose data contains a list of OutPoints**.
CKB reads that list and loads all referenced cells, as if you'd listed them individually.

```rust
Dep Group Cell:
    data: [ OutPoint_A, OutPoint_B, OutPoint_C ]

Transaction uses:
    cell_deps: [{ out_point: "dep group cell", dep_type: "dep_group" }]

CKB expands to effectively:
    cell_deps: [ Cell_A, Cell_B, Cell_C ]
```

## Why Dep Groups Exist

- Scripts sometimes need multiple cells (binary + helper data + shared libraries)
- Listing them individually in every transaction is verbose and repetitive
- Dep groups let you bundle them — one reference instead of many
- Especially useful for system scripts used by nearly every transaction

## Real-World Example

The default lock script (secp256k1-blake2b-sighash) needs two cells:
1. The script binary
2. A precomputed secp256k1 data table (for faster signature math)

A dep group cell bundles both. Every transaction references just the one dep group.

## Who Pays for Dep Group Cells?

- **System scripts** (default lock, etc.): deployed by the **Nervos Foundation** from genesis allocation — public infrastructure
- **Custom scripts**: the **script developer** deploys and pays capacity for their own dep group cells
- Since CKB is locked, not spent — deployers can destroy the cell and reclaim CKB anytime (but breaks references)

## Key Takeaways
- `dep_type: "code"` = use this cell's data directly
- `dep_type: "dep_group"` = this cell contains a list of OutPoints, expand them all
- Dep groups are a convenience mechanism — they don't add new capability, just reduce verbosity
- Commonly used for system scripts with multiple components

---

## Questions

**Q1:** A dep group cell is just a regular cell with data. Given CKB's economic model (1 CKByte = 1 byte), who pays for and maintains these dep group cells — and why would they bother?

<details>
<summary>Answer</summary>
For system scripts (like the default lock), the **Nervos Foundation** deploys and maintains them from the genesis block allocation — it's public infrastructure the whole network needs. For custom scripts, the **script developer** pays. They bother because without the dep group, every user of their script would need to list multiple cell_deps manually, making the developer experience worse.
</details>

---

**Q2:** What happens if someone destroys a dep group cell that many transactions depend on?

<details>
<summary>Answer</summary>
Any **new** transactions trying to reference that dep group would fail — CKB wouldn't be able to resolve the cell_dep. Existing confirmed transactions are unaffected (they're already on-chain). For system scripts, the Foundation's dep group cells are critical infrastructure — destroying them would break the most common lock script. In practice, the lock script on the dep group cell itself prevents unauthorized destruction.
</details>

---

**Q3:** Could you nest dep groups — a dep group that points to another dep group?

<details>
<summary>Hint</summary>
Think about what CKB does when it encounters `dep_type: "dep_group"` — does it recursively expand?
</details>

<details>
<summary>Answer</summary>
No. CKB only expands one level. When it encounters a dep group, it reads the list of OutPoints and loads those cells as `dep_type: "code"`. It does not recursively check if any of those cells are themselves dep groups. This keeps the resolution simple and predictable.
</details>

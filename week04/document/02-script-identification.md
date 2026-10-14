# Script Identification — How CKB Finds and Runs Code

> **Curriculum Reference:** [truthixify/learn-ckb-in-45-minutes / 02-script-identification.md](https://github.com/truthixify/learn-ckb-in-45-minutes/blob/main/notes/02-script-identification.md)

## The Script Structure

```rust
Script {
    code_hash: H256,
    hash_type: "data" | "data1" | "data2" | "type",
    args: Bytes
}
```

## code_hash + hash_type

These two fields work together to locate the executable binary in `cell_deps`.

### hash_type: "data" / "data1" / "data2"
- `code_hash` = hash of the **code binary itself** (cell data)
- CKB searches `cell_deps` for a cell whose data hashes to that value
- Content-addressed lookup: "find the cell containing code with this exact hash"
- `data` / `data1` / `data2` select CKB-VM version v0, v1, v2 respectively

### hash_type: "type"
- `code_hash` = hash of the **type script** of the cell containing the code
- Lookup: "find the cell whose type script hashes to this, then load its data as code"
- Enables **upgradeable scripts** — deploy new code in a new cell with the same type script, and the code_hash still resolves

## args — Runtime Parameters

`args` is passed to the script at execution time. It customizes a generic script:

- **Default lock (secp256k1-blake2b):** args = owner's public key hash
- **UDT type script:** args = token issuer's lock hash

Same binary can secure millions of cells, each with different `args`.

## Full Resolution Example

```rust
Cell {
    capacity: 200 CKB,
    lock: {
        code_hash: 0x9bd7...,
        hash_type: "type",
        args: 0xabcd1234...
    },
    type: null,
    data: 0x
}
```

1. `hash_type: "type"` → search `cell_deps` for cell whose **type script** hashes to `0x9bd7...`
2. Load that cell's **data** as executable code
3. Run in CKB-VM, passing `args` (`0xabcd1234...`)
4. Script checks witness for valid signature matching that pubkey hash

## Upgradeability: "data" vs "type"

### Why "data" is immutable
- `code_hash` = hash of the binary bytes
- New code → new hash → existing cells still point to old hash
- You'd have to consume and recreate every cell to update them

### Why "type" is upgradeable
- `code_hash` = hash of the **type script** on the code cell
- Deploy a new cell with the **same type script** but **new code** in its data
- Existing cells resolve to the new code automatically

### Analogy
- **"data"** = "give me book with ISBN X" — new edition gets a new ISBN, old request still fetches old book
- **"type"** = "give me whatever book is on shelf A, row 3" — librarian swaps in new edition, you get the update

### Trade-off

| | `data` | `type` |
|---|---|---|
| Security | Stronger — code can never change | Requires trusting whoever controls the type script |
| Flexibility | None — locked forever | Can fix bugs, add features |
| Use case | Core/high-security scripts | Application scripts that may evolve |

## How the Type Script Hash is Generated

The hash is computed over the **type script fields** of the code cell — NOT the cell's data.

```rust
Code Cell {
    capacity: 500 CKB,
    lock: { ... },
    type: {                          ← THESE fields get hashed
        code_hash: 0xabc...,
        hash_type: "data1",
        args: 0x1234...
    },
    data: <the actual RISC-V binary>  ← NOT part of the hash
}
```

hash = blake2b(type.code_hash + type.hash_type + type.args)

### Two layers
```text
Layer 1 (identity, stable):    Type Script { code_hash, hash_type, args }  → hashed
Layer 2 (content, swappable):  Cell Data = <RISC-V binary>                 → NOT hashed
```

- The type script acts as a **permanent address** for the code cell
- On upgrade: consume old code cell, create new one with same type script but new data
- The type script typically enforces that only the deployer can update the cell
- Since type script fields don't change, the hash stays the same, and all references still resolve

## Key Takeaways
- `code_hash` + `hash_type` = how CKB locates script code in cell_deps
- "data" variants = content-addressed (immutable scripts)
- "type" = indirection through type script hash (upgradeable scripts)
- `args` = per-cell configuration passed to the script at runtime
- Immutable scripts are more secure but can't be patched; upgradeable scripts are flexible but require trust

---

## Questions

**Q1:** Why would using `hash_type: "type"` make a script "upgradeable," while `hash_type: "data"` makes it immutable? Think about what happens if you want to fix a bug in your deployed script.

<details>
<summary>Hint</summary>
With "data", the code_hash is the hash of the binary bytes themselves. With "type", it's the hash of something else entirely. What happens to each hash when the binary changes?
</details>

<details>
<summary>Answer</summary>
With `hash_type: "data"`, the `code_hash` is the hash of the actual binary. If you deploy fixed code, the bytes change, so the hash changes. All existing cells still point to the old hash — they'll keep using the old buggy code forever. With `hash_type: "type"`, the `code_hash` is the hash of the **type script** on the code cell. You can deploy a new cell with the same type script but different code in its data field. The hash doesn't change, so all existing cells automatically resolve to the new code.
</details>

---

**Q2:** If you're building a high-stakes vault contract where users need to be 100% certain the rules can never change, would you use `hash_type: "data"` or `"type"`? Why?

<details>
<summary>Answer</summary>
`hash_type: "data"`. Since the code_hash is tied to the exact binary, no one — not even the deployer — can change the script code after deployment. Users can verify the code hash matches the audited source code and know with certainty it will never change. With "type", you'd need to trust the deployer not to swap in malicious code.
</details>

---

**Q3:** If someone uses `hash_type: "type"` for their lock script, who do they need to trust? What could go wrong if that trust is misplaced?

<details>
<summary>Answer</summary>
They need to trust whoever controls the **type script on the code cell** — typically the script deployer. If that trust is misplaced, the deployer could swap in malicious code that bypasses signature checks. Every cell using that lock script via `hash_type: "type"` would then be vulnerable — an attacker could spend anyone's cells without their private key.
</details>

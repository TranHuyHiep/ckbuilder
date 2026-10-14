# UDT — User Defined Tokens

> **Curriculum Reference:** [truthixify/learn-ckb-in-45-minutes / 05-user-defined-tokens.md](https://github.com/truthixify/learn-ckb-in-45-minutes/blob/main/notes/05-user-defined-tokens.md)

## Ethereum vs CKB Token Model

### Ethereum (ERC-20)
- Single contract holds a `mapping(address => balance)` for everyone
- All state in one place, token creator pays for deployment
- Receiving tokens is "free" — just a number update in the mapping

### CKB (UDT)
- **Each user holds their own token cell**
- Balance is in the cell's `data` field (uint128, 16 bytes)
- The `type script` identifies which token it is
- The `lock script` determines who can spend it

## UDT Cell Structure

```rust
Cell {
    capacity: 142 CKB,
    lock: { <owner's lock script> },
    type: {
        code_hash: 0xUDT_SCRIPT...,
        hash_type: "type",
        args: 0xISSUER_LOCK_HASH     # identifies which token
    },
    data: <balance as uint128 little-endian>
}
```

## Transfer Example: Alice sends 200 of 500 tokens to Bob

```text
Inputs:                          Outputs:
┌─────────────────────┐          ┌─────────────────────┐
│ Alice's UDT cell    │          │ Alice's UDT cell    │
│ lock: Alice         │          │ lock: Alice         │
│ type: UDT           │    →     │ type: UDT           │
│ data: 500           │          │ data: 300           │
└─────────────────────┘          ├─────────────────────┤
                                 │ Bob's UDT cell      │
                                 │ lock: Bob           │
                                 │ type: UDT           │
                                 │ data: 200           │
                                 └─────────────────────┘
```

### What validates the transfer
- **Type script** (runs on inputs + outputs): checks total in = total out (500 = 300 + 200)
- **Lock script** (runs on Alice's input): verifies Alice authorized the spend

## Storage Cost for Token Holders

Minimum UDT cell size ≈ **142 bytes → 142 CKB locked**

Breakdown: 8 (capacity) + 32+1+20 (lock) + 32+1+32 (type) + 16 (data)

### Implications
- Receiving tokens is NOT free — receiver needs a cell with locked CKB
- Either:
  - Receiver already has an empty UDT cell ready
  - Sender includes extra CKB to create the receiver's cell
  - dApp handles cell creation for the user

## Minting — How New Tokens Are Created

The type script's `args` = **issuer's lock hash**. This controls who can mint.

### Type script logic:
1. UDT cells in inputs exist AND total in = total out → **normal transfer** ✓
2. Output total > input total (or no UDT inputs) → **minting** → check if issuer authorized it

### How issuer authorization works:
- Type script checks if any input cell's **lock script hash matches `args`**
- If yes → issuer's lock script ran and passed → issuer is present → mint allowed

```text
Minting 1000 tokens:

Inputs:                              Outputs:
┌──────────────────────┐             ┌────────────────────────┐
│ Issuer's regular cell│             │ Issuer's regular cell  │
│ lock: Issuer         │  ← proves   │ lock: Issuer           │
│ type: null           │  issuer     │ type: null             │
│ data: 0x             │  approved   │ data: 0x               │
└──────────────────────┘    this     ├────────────────────────┤
                            tx       │ New UDT cell           │
                                     │ lock: Issuer           │
                                     │ type: UDT(args=Issuer) │
                                     │ data: 1000             │
                                     └────────────────────────┘
```

- No special mint function needed
- Issuer just includes one of their cells in inputs
- Lock script validates their identity, type script sees "issuer is present"
- **Ownership proof = having a cell in inputs whose lock script passes**

## Key Takeaways
- UDT state is distributed — each user holds their own balance in their own cell
- Type script enforces token rules (conservation, issuance)
- Lock script enforces ownership
- Holding any token requires locking ~142 CKB for storage
- Fundamentally different UX from Ethereum — "someone has to pay for the cell"

---

## Questions

**Q1:** In the transfer example, the type script checks that input tokens equal output tokens. But what about when tokens are first created (minted)? There are no input UDT cells, only outputs. How does the type script decide who is allowed to mint?

<details>
<summary>Hint</summary>
Look at what's stored in the type script's `args` field and think about what the type script could check in the transaction's inputs.
</details>

<details>
<summary>Answer</summary>
The type script's `args` contains the issuer's lock hash. When it detects a mint (output total > input total), it checks if any input cell's lock script hash matches the `args`. If yes, the issuer's lock script ran and passed — meaning the issuer authorized this transaction. The issuer doesn't need a special mint function; they just include one of their own cells in the inputs.
</details>

---

**Q2:** Could anyone other than the issuer accidentally trigger a mint?

<details>
<summary>Answer</summary>
No. The only way to mint is to have an input cell whose lock script hash matches the type script's `args` — and that `args` is set to the issuer's lock hash. Unless someone breaks the issuer's cryptography (steals their private key), they can't produce a valid input cell that would pass the check.
</details>

---

**Q3:** Alice wants to send tokens to Bob, but Bob has never held this token before and doesn't have a UDT cell. What needs to happen?

<details>
<summary>Hint</summary>
Remember: every cell needs CKB capacity. A UDT cell requires ~142 CKB minimum.
</details>

<details>
<summary>Answer</summary>
Someone needs to provide the CKB to create Bob's UDT cell. Either: (1) Alice includes extra CKB in the transaction to fund Bob's new cell, (2) Bob has already set up an empty UDT cell ahead of time, or (3) the dApp handles cell creation on Bob's behalf. This is a key UX difference from Ethereum where receiving tokens is "free."
</details>

---

**Q4:** On Ethereum, all ERC-20 balances live in one contract. On CKB, each user has their own cell. What are the implications for querying "how many holders does this token have?"

<details>
<summary>Answer</summary>
It's much harder on CKB. On Ethereum, you can read the contract's storage. On CKB, you'd need to scan the entire set of live cells to find all cells with the matching UDT type script. This requires an indexer — a service that scans and indexes cells by type script. This is a trade-off of the distributed state model.
</details>

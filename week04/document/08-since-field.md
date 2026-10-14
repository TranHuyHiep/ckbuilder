# The `since` Field — Time Locks

> **Curriculum Reference:** [truthixify/learn-ckb-in-45-minutes / 08-since-field.md](https://github.com/truthixify/learn-ckb-in-45-minutes/blob/main/notes/08-since-field.md)

## What It Does
Every transaction input has an optional `since` field that prevents the input from being consumed until a condition is met.

Similar to Bitcoin's `nLockTime`/`nSequence` but more flexible.

## Lock Metrics

| Lock by | Meaning |
|---|---|
| **Block number** | "Can't spend until block N" |
| **Epoch** | "Can't spend until epoch N" |
| **Timestamp** | "Can't spend until Unix time N" |

## Relative vs Absolute

- **Absolute** — "can't spend until block 500,000 (globally)"
- **Relative** — "can't spend until 100 blocks **after this cell was created**"
- Relative locks start counting from the cell's confirmation time

## Bit Layout

```text
since (64 bits):
┌──────────┬──────────┬─────────────────────────────┐
│ bit 63   │ bits 62  │ bits 0-61                    │
│ relative │ metric   │ value                        │
│ flag     │ type     │                              │
├──────────┼──────────┼─────────────────────────────┤
│ 0=abs    │ 00=block │ the actual number            │
│ 1=rel    │ 01=epoch │                              │
│          │ 10=time  │                              │
└──────────┴──────────┴─────────────────────────────┘
```

## Use Cases
- **Vesting contracts** — tokens unlock N epochs after deposit
- **Payment channels** — either party can close after a timeout
- **DAO withdrawals** — Nervos DAO enforces minimum lock period via `since`
- **Dispute windows** — challenge period must pass before finalization

## Connection to `header_deps`
- Timestamp/epoch-based locks need current chain time
- Scripts read this from block headers via `header_deps` and `ckb_load_header` syscall

## Why the DAO Uses Epoch-Based Locks

An epoch in CKB ≈ **4 hours of real time**, regardless of hash rate changes:
- Hash rate doubles → faster blocks, but more blocks per epoch → still ~4 hours
- Hash rate drops → slower blocks, fewer per epoch → still ~4 hours

| Metric | Problem for DAO |
|---|---|
| Block number | Hash rate spikes → 1000 blocks pass in minutes, lock period unpredictable |
| Timestamp | Miners can slightly manipulate reported timestamps |
| Epoch | Stable real-world duration, not manipulable, predictable |

DAO lock ≈ 180 epochs ≈ 30 days, consistent regardless of network conditions.

## Key Takeaways
- `since` adds conditional time locks to cell spending
- Supports block number, epoch, or timestamp — both absolute and relative
- Enables time-dependent logic without custom script code
- Relative locks are especially powerful for contracts with countdown mechanics

---

## Questions

**Q1:** The Nervos DAO uses epoch-based locks rather than block-number or timestamp-based. Why?

<details>
<summary>Hint</summary>
What happens to block numbers if hash rate suddenly doubles? What can miners do with timestamps?
</details>

<details>
<summary>Answer</summary>
An epoch in CKB ≈ 4 hours of real time, regardless of hash rate changes. Block numbers are unpredictable in real time — a hash rate spike means 1000 blocks could pass in minutes. Timestamps can be slightly manipulated by miners. Epochs provide a stable, manipulation-resistant, real-world time measurement. The DAO lock of ~180 epochs ≈ 30 days stays consistent regardless of network conditions.
</details>

---

**Q2:** Design a simple vesting contract: an employer locks tokens for an employee, and the tokens unlock 1 year after deposit. Which `since` configuration would you use?

<details>
<summary>Answer</summary>
Use a **relative epoch-based** since. Set bit 63 = 1 (relative), metric type = epoch, and value ≈ 2190 epochs (365 days × 6 epochs/day). This means the cell can't be spent until ~2190 epochs after it was confirmed on-chain. Relative is key here — it starts counting from when the cell was created, not from a fixed point in time.
</details>

---

**Q3:** Can the `since` field alone implement complex time-lock logic, or does it need cooperation from scripts?

<details>
<summary>Answer</summary>
The `since` field is enforced by the **CKB consensus layer** itself — no script code needed for the basic time lock. However, for more complex logic (e.g., "unlock after 30 days OR immediately if both parties agree"), you'd need the lock script to implement the alternative conditions. The `since` field handles the time constraint, and the script handles the rest.
</details>

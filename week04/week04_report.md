# CKB Weekly Report - Week 4

**Reporting period:** 08 October - 14 October 2026  
**Publication date:** 14 October 2026  
**Participant:** Tran Huy Hiep  

---

## 1. Week 4 Overview

The goal of Week 4 was to study the **core architecture and system mechanics** of Nervos CKB by following the [learn-ckb-in-45-minutes](https://github.com/truthixify/learn-ckb-in-45-minutes) curriculum.

Key topics covered:
- **Cell Model:** Cell anatomy (capacity, lock, type, data), state immutability, Live vs Dead cells.
- **Transaction Structure:** `cell_deps`, input/output lifecycle, Lock Script vs Type Script roles.
- **Script Identification:** `code_hash` + `hash_type` resolution; `data` (immutable) vs `type` (upgradeable).
- **NC-Max Consensus:** Two-step propose/commit confirmation; orphan-rate targeted difficulty.
- **Economic Model:** 1 CKByte = 1 Byte state rent; Base & Secondary Issuance; Nervos DAO.
- **User Defined Tokens (UDT):** Cell-owned token model vs ERC-20; supply conservation invariant.
- **CKB-VM:** RISC-V IMC sandbox; syscall-based transaction access; cycle limits.
- **Off-Chain / On-Chain Pattern:** Off-chain computation + on-chain verification; no re-entrancy.
- **`since` Field:** 64-bit time-lock encoding (block/epoch/timestamp, relative/absolute).
- **Dep Groups:** Bundling multiple `cell_deps` into a single OutPoint reference.
- **Witness Layout:** Separating cryptographic proofs from persistent on-chain state.

All notes have been formatted and stored in [`document/`](./document/).

---

## 2. Development Environment

- **OS:** macOS (Apple Silicon arm64)
- **CKB Tools:** CKB-CLI v1.12.0, CKB Debugger v1.1.1, OffCKB v0.4.13
- **Networks:** Nervos Meepo Testnet, Local Devnet (`http://127.0.0.1:8114`)
- **Study Source:** [`truthixify/learn-ckb-in-45-minutes`](https://github.com/truthixify/learn-ckb-in-45-minutes)

---

## 3. Evidence

| Module | Topic | Reference |
|---|---|---|
| 00 | Cells & Scripts | [`00-cells-and-scripts.md`](./document/00-cells-and-scripts.md) |
| 01 | Transaction Structure | [`01-transaction-structure.md`](./document/01-transaction-structure.md) |
| 02 | Script Identification | [`02-script-identification.md`](./document/02-script-identification.md) |
| 03 | NC-Max Consensus | [`03-nc-max-consensus.md`](./document/03-nc-max-consensus.md) |
| 04 | Economic Model | [`04-economic-model.md`](./document/04-economic-model.md) |
| 05 | User Defined Tokens | [`05-user-defined-tokens.md`](./document/05-user-defined-tokens.md) |
| 06 | CKB-VM Execution | [`06-ckb-vm-execution.md`](./document/06-ckb-vm-execution.md) |
| 07 | Off-Chain / On-Chain | [`07-offchain-onchain-pattern.md`](./document/07-offchain-onchain-pattern.md) |
| 08 | `since` Time Locks | [`08-since-field.md`](./document/08-since-field.md) |
| 09 | Dep Groups | [`09-dep-groups.md`](./document/09-dep-groups.md) |
| 10 | Witness Layout | [`10-witness-layout.md`](./document/10-witness-layout.md) |

---

## 4. Key Concepts Learned

- **Cell Model vs Account Model:** CKB state is explicit, user-owned, and immutable. Transactions destroy old cells and create new ones — no in-place mutation. Enables native parallelism.
- **Lock vs Type Script:** Lock = authorization (runs on inputs only); Type = state validation (runs on both inputs and outputs). This separation cleanly divides *who can spend* from *is the transition valid*.
- **State Rent:** Locking 1 CKByte stores 1 byte on-chain. Secondary issuance dilutes locked CKB, effectively charging rent to miners. Nervos DAO deposits neutralize dilution for non-storing holders.
- **NC-Max:** Propose/commit separation lets nodes pre-validate transactions, eliminating block propagation bottlenecks and allowing safe faster block times.
- **CKB-VM:** Pure RISC-V validation — no side effects, no inter-contract calls. Re-entrancy attacks are architecturally impossible.

---

## 5. Week 4 Development Log

| Date | Activity | Result |
|---|---|---|
| 08 Oct 2026 | Cell Model & Scripts study | Completed |
| 09 Oct 2026 | Transaction structure & script identification | Completed |
| 10 Oct 2026 | NC-Max consensus & economic model | Completed |
| 11 Oct 2026 | UDT token model | Completed |
| 12 Oct 2026 | CKB-VM execution & syscalls | Completed |
| 13 Oct 2026 | Off-chain pattern, `since` field, dep groups | Completed |
| 14 Oct 2026 | Witness layout, document formatting & report | Completed |

---

## 6. Challenges & Solutions

- **Cell Model mental shift:** Coming from Ethereum's account model, thinking in terms of cell creation/destruction rather than variable mutation required significant reframing.
- **`since` field bitwise layout:** Parsing the 64-bit encoding (relative/absolute flag, metric type, value) required careful bit-masking analysis.
- **State Rent economics:** Understanding that secondary issuance acts as implicit storage rent (not explicit gas) clarified why Nervos DAO is essential for long-term CKB holders.

---

## 7. Final Reflection

Week 4 provided the theoretical foundation connecting all prior hands-on work. Understanding CKB-VM, NC-Max, state rent, and the Cell Model enables confident, idiomatic protocol-level development going forward.

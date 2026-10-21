# CKB Weekly Report - Week 5

**Reporting period:** 15 October - 21 October 2026  
**Publication date:** 21 October 2026  
**Participant:** Tran Huy Hiep  

---

## 1. Week 5 Overview

The goal of Week 5 was to study the **Perun state channel framework** and run it hands-on, first on a local Ethereum chain and then on a **CKB devnet**, connecting it to the Fiber payment channels of Week 3 and the Cell Model of Week 4.

Key accomplishments:
- Studied the Perun protocol and framework ([Perun Introduction](https://hyperledger-labs.github.io/perun-doc/introduction.html)), with notes in English and Vietnamese, updated against the current source code.
- Wrote **5 labs** with ready-to-run code (`week05/code/run.sh`), each with real output, exercises and answers.
- Ran a full payment channel lifecycle (open → off-chain updates → settle) on Hardhat, and counted which steps touch the chain.
- Attacked the channel and reproduced a **real bug in the official example** (missing `return` after `Reject`), which crashes the program; verified the fix.
- Compared the cooperative settle path (~0.4 s) with the **on-chain dispute path** (15–40 s, one extra transaction).
- Built the Perun CKB contracts (RISC-V), started a CKB devnet and ran channels with **CKByte and SUDT**, including inspecting the live channel cell and funds cells and a full dispute → ForceClose.
- Found and fixed **5 problems** that blocked the CKB example on this server (toolchain, headers, devnet script, a race condition in `perun-ckb-backend`).

---

## 2. Development Environment

- **OS:** Ubuntu 26.04.1 LTS aarch64 (VMware Fusion VM, disk expanded from 20 GB to 49 GB) / macOS host
- **Languages & tools:** Go 1.26.8, Node.js 22.22.1, Hardhat 2.29.1, Rust 1.85.0 (pinned for the contracts), clang-18, `riscv64-unknown-elf-gcc`
- **CKB tools:** `ckb` v0.210.0, `ckb-cli` v1.16.0 (upgraded from 1.12.0), tmux/tmuxp
- **Perun libraries:** go-perun v0.15.0, perun-eth-backend v0.6.0 (Ethereum); perun-ckb-backend `v0.2.1-0.20250603…` with go-perun `v0.13.1-…` (CKB)
- **Source:** [`perun-network/perun-examples`](https://github.com/perun-network/perun-examples) at commit `6e975e7`
- **Networks:** Local Hardhat chain (`ws://127.0.0.1:8545`, chain ID 1337) and local CKB devnet (`http://127.0.0.1:8114`)

---

## 3. Deliverables

| Item | Contents | Reference |
|---|---|---|
| Reading notes | Perun protocol, framework components, status vs. current code | [`perun_introduction_en.md`](./document/perun_introduction_en.md), [`perun_introduction_vi.md`](./document/perun_introduction_vi.md) |
| Lab guide | Roadmap, server state, how to run | [`labs/README.md`](./labs/README.md) |
| Lab 01 | Environment check and setup | [`lab01_setup.md`](./labs/lab01_setup.md) |
| Lab 02 | First payment channel & client architecture | [`lab02_payment_channel.md`](./labs/lab02_payment_channel.md) |
| Lab 03 | Validation, bug hunting & two-sided funding (10 scenarios) | [`lab03_validation_and_funding.md`](./labs/lab03_validation_and_funding.md) |
| Lab 04 | On-chain dispute | [`lab04_dispute.md`](./labs/lab04_dispute.md) |
| Lab 05 | Perun on CKB: contracts and devnet | [`lab05_ckb.md`](./labs/lab05_ckb.md) |
| Lab code | `run.sh`, Go modules for Ethereum (`eth/`) and CKB (`ckb/`), helper scripts | [`code/`](./code/) |

---

## 4. Evidence

| ID | Task / Operation | Key Result | Evidence |
|---|---|---|---|
| 01 | Start local Ethereum chain | Hardhat node on `127.0.0.1:8545`, 3 accounts (Deployer, Alice, Bob) with 10 ETH each | [`1_list_account.png`](./evidence/1_list_account.png) |
| 02 | Run the original payment channel demo | Deploy → open → pay → settle; two `ConcludedEvent`; wallets 10/10 → **6.9997 / 12.9999 ETH** | [`2_payment_channel.png`](./evidence/2_payment_channel.png) |
| 03 | Lab 03 attack `mint` | `MintForMe(100)` rejected by go-perun: `unequal allocation: invalid state transition`; balances unchanged | [`3_mint.png`](./evidence/3_mint.png) |
| 04 | Lab 04 fast settle path | Finalize + settle in **479 ms**; 6 on-chain transactions (1 deposit, 1 conclude, 2 withdraw) | [`4_fast.png`](./evidence/4_fast.png) |
| 05 | Lab 05 on CKB devnet | `dispute`: Dispute on-chain → ForceClose in **34.0 s**; `default`: channel 100/100 → 90/110 → 95/105 CKB, settle in **4.4 s** | [`5_exchange.png`](./evidence/5_exchange.png) |

---

## 5. Key Results

### Ethereum (Hardhat)

| Experiment | Result |
|---|---|
| 3 payments in a channel (Lab 02) | 4 channel transactions on-chain (1 deposit, 1 conclude, 2 withdraw), the 3 payments themselves are off-chain only |
| 100 micro-payments of 0.01 ETH (Lab 03) | ~42.7 s off-chain (~430 ms each on this VM), still only **6** on-chain transactions in total; as on-chain transactions with 5 s blocks it would take ≥ 8 min 20 s |
| `steal-buggy` / `fund52-buggy` (Lab 03) | Original example code panics: `multiple calls on channel update responder` and `nil pointer dereference` |
| `steal` / `fund52-default` (Lab 03, fixed code) | Bob rejects cleanly, Alice receives the error, program continues |
| Fast settle vs dispute (Lab 04) | Fast: ~0.4 s, 6 transactions. Dispute: 15–20 s (challenge 10 s) and 34–40 s (challenge 30 s), 7 transactions (extra `register`) |

### CKB devnet

| Scenario | Result |
|---|---|
| `default` | Each side deposits 100 CKB; balances 100/100 → 90/110 → 95/105; open ~11 s, settle ~4.4 s |
| `inspect` | While open: 1 channel cell (804 CKB, lock PCLS, type PCTS, 256 bytes data) + 2 funds cells (100 CKB each, lock PFLS); all consumed after settle |
| `sudt` | SUDT in channel 50/50 → 40/60; wallets: Alice −10 SUDT, Bob +10 SUDT |
| `dispute` | Dispute tx on-chain, both watchers receive `Registered` (version 2), ForceClose after the challenge; ~34 s total |

---

## 6. Key Concepts Learned

- **State channels:** Funds are locked on-chain once; any number of updates are exchanged off-chain as states signed by both parties; only the final state is settled on-chain. The highest-version state signed by both wins.
- **Two layers of validation:** go-perun enforces protocol invariants (conserved total, increasing version, no updates after final, valid signatures). Business rules (who pays whom, deposit limits, minimum payment) must be written by the application in `HandleProposal` / `HandleUpdate`. A signed bad state cannot be undone.
- **Dispute path:** Without cooperation, a party registers the latest state on-chain, waits out the challenge duration, then concludes. The counterparty's watcher refutes outdated states during that window, so the challenge duration must be long enough for the victim to come back online.
- **Blockchain-agnostic design:** go-perun's core knows no chain; each chain provides `Funder`, `Adjudicator` and `Wallet` (perun-eth-backend, perun-ckb-backend, …).
- **Perun on the Cell Model:** A channel is a group of cells. The **channel cell** (lock PCLS, type PCTS) stores `ChannelStatus` in its data and `ChannelConstants` in the type script args; **funds cells** (lock PFLS) hold each party's deposit and can only be spent together with the channel cell. Uniqueness comes from a thread token (an OutPoint, same idea as Type ID); ForceClose is time-locked with the `since` field (Week 04).
- **Compared with Fiber (Week 3):** both are CKB payment channel systems; Perun is a general framework (state channels, virtual channels, multiple chains), while Fiber is a Lightning-style network with HTLC/TLC routing.

---

## 7. Week 5 Development Log

| Date | Activity | Result |
|---|---|---|
| 15–21 Oct 2026 | Read Perun introduction, protocol and framework docs; wrote notes (EN/VI) | Completed |
| 21 Oct 2026 | Updated notes and labs to the server's real state; reduced labs from 10 to 5 | Completed |
| 21 Oct 2026 | Expanded the VM disk (20 → 49 GB) in VMware Fusion | Completed |
| 21 Oct 2026 | Coded and ran Labs 01–04 on Hardhat (`run.sh`, 10 Lab 03 scenarios, dispute timing) | Completed |
| 21 Oct 2026 | Installed CKB tooling, built contracts, ran Lab 05 on CKB devnet (4 scenarios) | Completed |
| 21 Oct 2026 | Translated labs and code to English; weekly report | Completed |

---

## 8. Challenges & Solutions

| Challenge | Cause | Solution |
|---|---|---|
| Disk only had ~430 MB free | 19 GB VM disk almost full | Expanded the disk in VMware Fusion, then `growpart` + `resize2fs` online |
| Outdated documentation | perun-doc still uses `ganache-cli`, `UpdateBy`, says virtual channels are unimplemented | Labs follow the current code (Hardhat, `ch.Update`, go-perun v0.15.0 with virtual channels) |
| Bug in the official example | `HandleProposal` / `HandleUpdate` miss `return` after `Reject` → panic | Added `return`; kept a `Buggy` flag so the crash can be reproduced in Lab 03 |
| `Cannot select: AtomicLoadAdd` building contracts | Contracts build with `-C target-feature=-a`; Rust 1.99 can no longer compile the atomics in `bytes` under it | Pinned Rust **1.85.0** (README-recommended) for the contracts directory only, via `rustup override` |
| `gnu/stubs-lp64.h` not found | clang-18 targets soft-float riscv64 (`lp64`), Ubuntu only ships `lp64d` headers | Added `lp64` copies of the two headers |
| `Method not found` during contract deployment | `setup-devnet.sh` failed to enable the Indexer RPC because ckb 0.210 adds `"Terminal"` to the module list | Fixed the `sed` line in `setup-devnet.sh` |
| SUDT balance 0 after devnet start | SUDT issued 1 s after deployment, before deploy transactions were committed | `devnet.sh` re-issues SUDT until the balance is non-zero |
| `dispute` scenario hangs forever | Race in `perun-ckb-backend`: its subscription's first poll comes after the Dispute tx is committed, so `Registered` is never emitted to go-perun | `ForceSettle` bounds the first `Settle` with a timeout, then calls `Settle` again (phase is already `Registered`) to run ForceClose |
| libp2p relay `failed to dial` (occasional) | go-perun's libp2p wire uses a public relay server | Re-run; documented in the lab troubleshooting table |

---

## 9. Final Reflection

Week 5 tied the earlier weeks together: the payment channel idea from Fiber (Week 3) became concrete at the code level with Perun, and the Cell Model, `since`, Type ID and witnesses from Week 4 appeared directly in the Perun CKB contracts. Running the same channel on Ethereum and on CKB made the difference between the account model (one Adjudicator contract with storage) and the Cell Model (a channel cell plus funds cells) very clear. The most valuable lessons came from things that broke: a missing `return` in the official example that crashes both clients, and a timing race in the CKB backend that only appears on a fast devnet. Both show that in channel systems, correctness depends as much on the off-chain client code as on the on-chain contracts.

# Perun Labs — Hands-on Perun, step by step

These labs follow the reading [Perun Introduction](../document/perun_introduction_en.md).
There are **5 labs**, and each one builds on the previous. The first four run on a local Ethereum chain (Hardhat), the last one on a CKB devnet.

## Roadmap

| # | Lab | Main topics | Chain | Time | Difficulty |
|---|---|---|---|---|---|
| 01 | [Set up the server environment](./lab01_setup.md) | Check the machine, install Hardhat, test build (`./run.sh check`) | — | 15 min | ⭐ |
| 02 | [First payment channel & client architecture](./lab02_payment_channel.md) | Channel lifecycle open → update → settle, Funder/Adjudicator/Watcher/Bus | Hardhat | 30–45 min | ⭐⭐ |
| 03 | [Validation, bug hunting & two-sided funding](./lab03_validation_and_funding.md) | `HandleProposal`/`HandleUpdate`, missing-`return` bug, attacking the channel | Hardhat | 45–60 min | ⭐⭐⭐ |
| 04 | [On-chain dispute](./lab04_dispute.md) | Uncooperative settle, `RegisteredEvent`, challenge duration | Hardhat | 30 min | ⭐⭐⭐ |
| 05 | [Perun on CKB](./lab05_ckb.md) | Read the Rust contracts (cell model) and run a channel on a CKB devnet | CKB devnet | 1.5–2 h | ⭐⭐⭐⭐ |

App channels (Tic-Tac-Toe) and virtual channels are no longer separate labs. Lab 04 has a further-reading section on both.

## Server state (checked 2026-10-21)

| Item | Value |
|---|---|
| OS | Ubuntu 26.04.1 LTS, **aarch64 (ARM64)**, VMware Fusion VM |
| User | `root` |
| CPU / RAM | 2 cores / 3.3 GB |
| Disk `/` | 49 GB (expanded in VMware Fusion), ~30 GB free |
| Example code | Already cloned at **`/root/perun-examples`** (commit `6e975e7`) |

| Tool | Version on the server | Needed for | Notes |
|---|---|---|---|
| `go` | 1.26.8 (`/usr/bin/go`, apt package) | Labs 01–05 | ✅ enough (needs ≥ 1.23) |
| `node` / `npm` | 22.22.1 / 9.2.0 | Labs 01–04 | ✅ enough |
| `git`, `jq`, `make`, `curl`, `wget` | present | all | ✅ |
| `cargo` / `rustc` | 1.99.0 (rustup `stable`) | Lab 05 | too new for the contracts; `lab05-tools` pins **1.85.0** for the contracts directory only |
| `clang` | 21.1.8 | Lab 05 | the build script needs exactly `clang-18` → install it |
| `ckb-cli` | 1.12.0 (`/usr/local/bin`) | Lab 05 | ❌ needs ≥ 1.13.0 |
| `ckb` | not installed | Lab 05 | ❌ needs ≥ 0.201.0 |
| `tmux`, `tmuxp`, `expect` | not installed | Lab 05 | ❌ install with apt |
| `riscv64-unknown-elf-gcc` | not installed | Lab 05 | ❌ install with apt |
| `python` (the command) | missing, only `python3` | Lab 05 | ❌ `devnet/fund_accounts.sh` calls `python`; install `python-is-python3` |

`./run.sh lab05-tools` installs everything marked ❌.

## Reference source versions

| Component | Version |
|---|---|
| `perun-examples/payment-channel` | go-perun v0.15.0, perun-eth-backend v0.6.0, go-ethereum v1.15.11 |
| Hardhat (installed from `package.json`, `^2.0.0`) | 2.29.1 |
| `perun-examples/payment-channel-ckb` | go-perun `v0.13.1-0.20250528…`, perun-ckb-backend `v0.2.1-0.20250603…`, ckb-sdk-go v2.2.0 (perun-network fork) |
| CKB contracts | Already present in `payment-channel-ckb/devnet/contracts` (not a submodule) |

> **Note:** perun-doc still describes `ganache-cli` and old APIs such as `UpdateBy`. The current code uses **Hardhat**, `ch.Update(...)` and addresses of the form `map[wallet.BackendID]...`. These labs follow the current code.

## How to do the labs
Every exercise is **already coded** in [`week05/code/`](../code/). There is no need to edit `/root/perun-examples`. Each step is one command:

```sh
cd /root/ckbuilder/week05/code
./run.sh help            # list every command
./run.sh check           # Lab 01
./run.sh lab02           # Lab 02
./run.sh lab03 all       # Lab 03, all 10 scenarios
./run.sh lab04 all       # Lab 04, fast path and dispute path
./run.sh lab05-tools     # Lab 05, install the CKB tooling (asks first)
./run.sh lab05-devnet start && ./run.sh lab05 inspect
```
For Labs 02–04, `run.sh` starts a fresh Hardhat, runs the scenario, prints a table of on-chain transactions and stops Hardhat, so you do not need two terminals. Each run's log is saved in `week05/code/logs/`. Your job is to read the output, read the code it points to, and answer the questions in each lab.

## Verified on the server (2026-10-21)
- ✅ `./run.sh check`: Go 1.26.8, Node 22, Hardhat 2.29.1, builds `payment-channel` and `code/eth`.
- ✅ `./run.sh lab02`, `./run.sh lab03 all` (10 scenarios), `./run.sh lab04 all`: all finish; the output quoted in each lab is real output.
- ✅ `code/ckb` builds. The aarch64 downloads for `ckb` v0.210.0 and `ckb-cli` v1.16.0 exist. All apt packages are available on Ubuntu 26.04.
- ✅ Lab 05: `lab05-tools`, `lab05-build` (Rust 1.85.0, about 25 s after the first build), `lab05-devnet start` on a fresh devnet, and the scenarios `default`, `inspect`, `sudt`, `dispute`. Five problems found while running it are fixed; see the table at the top of [Lab 05](./lab05_ckb.md).

## Troubleshooting
| Symptom | Cause | Fix |
|---|---|---|
| `panic: connecting to the relay server: failed to dial ...` | The go-perun libp2p wire connects Alice and Bob through a public relay server; the connection failed briefly | Run the same command again |
| `no space left on device` | Disk full during the Go build | Free space, then `go clean -cache` |
| `Port 8545 is in use ...` | Another Hardhat is running | Nothing to do, `run.sh` stops it and starts a fresh one |

## Shared variables
The labs use two variables. Add them to `~/.bashrc` once, in Lab 01:
```sh
export PERUN=/root/perun-examples
export LABS=/root/ckbuilder/week05
```

## Structure of each lab
- **Goals**
- **Prepared code:** which file holds what, to read.
- **Run:** the `./run.sh ...` command and its real output.
- **Exercises:** questions based on the output and the code.
- **Check questions:** hints are in a `<details>` block at the end. Answer first, then open it.
- **Checklist.**

Keep your notes and answers in `$LABS/notes/` (`mkdir -p $LABS/notes`).

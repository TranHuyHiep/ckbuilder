# Lab 05 — Perun on CKB: read the contracts and run a channel on a devnet

## Goals
- Understand how Perun is implemented on CKB's **Cell (UTXO) model**, compared with Ethereum's account model.
- Connect Week 04 (lock script, type script, witness, `since`) to a real system.
- Run an Alice ↔ Bob channel with **CKByte** and **SUDT** on a CKB devnet, and inspect the real channel cell and funds cells.

The lab has two parts. **Part A** is code reading only and can be done right away. **Part B** needs extra tooling (the install script is prepared).

> **Verification status (2026-10-21):** run end to end on the server: `lab05-tools`, `lab05-build`, `lab05-devnet start` on a fresh devnet, and all four scenarios `default`, `inspect`, `sudt`, `dispute`. The output quoted below is real output.

### Problems found and fixed while running it
| Problem | Cause | Fix (already applied) |
|---|---|---|
| `rustc-LLVM ERROR: Cannot select: AtomicLoadAdd` building `bytes` | The contracts build with `-C target-feature=-a` (CKB-VM has no atomics). Rust 1.99 (`stable` on the server) can no longer compile the atomics in `bytes` under that flag | `tools.sh` installs Rust **1.85.0** (the version the contract README recommends) and pins it with `rustup override` **only** for the contracts directory. `stable` and `1.95.0` (used by `/root/fiber`) are untouched |
| `fatal error: 'gnu/stubs-lp64.h' file not found` building `blake2b-rs` | clang-18 compiles for soft-float riscv64 (ABI `lp64`), but Ubuntu's `libc6-dev-riscv64-cross` only ships `lp64d` headers | `tools.sh` copies `stubs-lp64d.h` → `stubs-lp64.h` and `lib-names-lp64d.h` → `lib-names-lp64.h` in `/usr/riscv64-linux-gnu/include/gnu/` (adds 2 files, changes nothing existing) |
| `jsonrpc error: Method not found` during `deploy gen-txs` | `setup-devnet.sh` enables the Indexer RPC with `sed` on `"Debug"]`, but ckb 0.210 lists `"Terminal"` after `"Debug"`, so nothing matched | Fixed line 94 of `$PERUN/payment-channel-ckb/devnet/setup-devnet.sh` to append `"Indexer"` to the end of the `modules` list (uncommitted change in `/root/perun-examples`) |
| SUDT balance stays 0 after start | The tmux session runs `sudt_helper.sh fund` only 1 s after deploying, before the deploy transactions are committed | `devnet.sh start` checks Alice's SUDT balance and re-issues SUDT until it is non-zero |
| `dispute` hangs forever | Race in perun-ckb-backend: its polling subscription first polls 4 s after creation; on a fast devnet the Dispute tx is already committed, so it never emits `Registered` and go-perun waits forever | `ForceSettle` in `code/ckb/client/channel.go` bounds the first `Settle` with a timeout, then calls `Settle` again: the phase is already `Registered`, so it goes straight to ForceClose |

## Prepared code
| File | Contents |
|---|---|
| `code/ckb/tools.sh` | Installs everything: apt (tmux, tmuxp, expect, python-is-python3, RISC-V gcc, clang-18), Rust target, `ckb` v0.210.0, `ckb-cli` v1.16.0 (backs up the old one) |
| `code/ckb/build_contracts.sh` | Builds the RISC-V contracts |
| `code/ckb/devnet.sh` | `start` (like `make dev` but in the background, waiting until deployment is done), `stop`, `status` |
| `code/ckb/cells.sh` | Lists the live channel cell and funds cells via the `get_cells` RPC |
| `code/ckb/main.go` | Scenarios `default`, `sudt`, `inspect`, `dispute` |
| `code/ckb/client/` | Copy of the CKB example client, with `return` added after `Reject`, detailed event logs, and `ForceSettle` |
| `code/ckb/devnet` | Symlink to `$PERUN/payment-channel-ckb/devnet` |

---

## Part A — Read the Perun contracts on CKB

### Step A1. Print the code to read
```sh
cd $LABS/code
./run.sh lab05-code
```
The script prints the list of contracts, the Molecule types in `types.mol`, the `ChannelWitness` union, the PCTS function list, the whole PCTS `main()`, the `PChannelAction` enum, the call to `verify_time_lock_expired`, and the length of the two lock scripts.

The source is at `$PERUN/payment-channel-ckb/devnet/contracts` (already present, no clone needed).

| Script | Attached to | Role |
|---|---|---|
| `perun-channel-lockscript` (PCLS) | Lock of the **channel cell** | Only allows transactions with an input owned by one of the two parties |
| `perun-channel-typescript` (PCTS) | Type of the **channel cell** | State machine: checks Start, Progress and Close |
| `perun-funds-lockscript` (PFLS) | Lock of the **funds cells** | Funds can only be unlocked when the channel cell (PCTS) is in the same transaction |

`perun-vchannel-*` are for virtual channels; `sample-udt` is the SUDT used in tests.

### Part A exercises
Answer from the `lab05-code` output:

**A.1.** Draw the branch diagram of PCTS `main()`. Check it against the `PChannelAction` enum.

**A.2. Ethereum ↔ CKB mapping**

| Concept | Ethereum (Labs 02–04) | CKB |
|---|---|---|
| Where channel funds are held | `AssetHolder` contract | ? |
| Where the registered state is stored | `Adjudicator` storage | ? |
| Register and dispute | Call `register` | ? |
| Conclude | Call `conclude` | ? |
| Preventing duplicate channels | Channel ID is a hash of the params | ? |
| Measuring the challenge duration | `block.timestamp` | ? |

**A.3. Thread token.** `ChannelToken` holds an `OutPoint`. Why does that guarantee each channel is unique? Relate it to Type ID.

**A.4. ForceClose.** How does `check_normal_force_close` check time? Recall the `since` field from Week 04.

---

## Part B — Payment channel on a CKB devnet

### Step B1. Install the tooling (once, about 5–10 minutes)
```sh
cd $LABS/code
./run.sh lab05-tools
```
The script prints what it will do and asks `Continue? [y/N]`. Type `y`. The script:
- backs up the old `ckb-cli` 1.12.0 as `/usr/local/bin/ckb-cli-1.12.0.bak` and the keystore `~/.ckb-cli` as `~/.ckb-cli.bak`,
- does not touch the Rust `1.95.0` toolchain that `/root/fiber` uses,
- installs `python-is-python3`, because `devnet/fund_accounts.sh` calls `python` and the server only has `python3`.

It also applies the Rust 1.85.0 pin and the `lp64` header fix described above.

Run `./run.sh check` again; everything under Lab 05 should be ✔.

### Step B2. Build the contracts (5–15 minutes the first time)
```sh
./run.sh lab05-build
```
It ends with the list of binaries in `build/release`: `perun-channel-lockscript`, `perun-channel-typescript`, `perun-funds-lockscript`, `perun-vchannel-*`, `sample-udt`.

Check that the contracts directory uses Rust 1.85.0: `cd $PERUN/payment-channel-ckb/devnet/contracts && rustc --version` should print `rustc 1.85.0`.

### Step B3. Start the devnet (about 1–2 minutes)
```sh
./run.sh lab05-devnet start
```
The script runs `setup-devnet.sh` (rebuilds the contracts, creates the miner/alice/bob accounts, `ckb init --chain dev`), opens the `devnet` tmux session in the background, then waits until the contracts are deployed and SUDT has been issued. Every `start` **deletes the old devnet** and creates a fresh one.

Output on the server (fresh devnet):
```
>>> Alice has no SUDT yet, issuing SUDT (attempt 1)
>>> Alice has no SUDT yet, issuing SUDT (attempt 2)
>>> SUDT balance of Alice: 100000000
tmux session devnet: running
RPC http://127.0.0.1:8114: block 45
Contracts: deployed
RPC http://127.0.0.1:8114: block ...
Accounts: alice.pk alice.txt bob.pk bob.txt genesis-1.pk ... sudt-owner-lock-hash.txt
```
To watch the panes live: `tmux attach -t devnet` (switch panes with `Ctrl+b` then an arrow key, switch windows with `Ctrl+b n`, detach with `Ctrl+b d`).

### Step B4. Run the scenarios
```sh
./run.sh lab05 default    # the original scenario
./run.sh lab05 inspect    # pauses 30 s while the channel is open, inspects cells before and after settling
./run.sh lab05 sudt       # a channel holding both CKB and SUDT
./run.sh lab05 dispute    # Alice settles without finalizing
```
The first run downloads the perun-ckb-backend Go modules, which takes a few minutes. Logs are saved as `code/logs/lab05-*.log`.

The `default` scenario (same as the original `payment-channel-ckb/main.go`):
1. Alice opens a channel, **each side deposits 100 CKB** (the CKB example's `HandleProposal` requires equal funding).
2. Alice pays Bob 10 CKB.
3. Bob pays Alice 5 CKB.
4. Alice settles.

The backend prints many debug lines (`signature without padding`, `getExactChannelLiveCell`, …). To see only the important ones:
```sh
./run.sh lab05 default 2>&1 | grep -E "allocation|==>|EVENT|Wallet"
```

Output on the server (2026-10-21), shortened:

**`default`**
```
==> Opening the channel (Start + Fund on-chain) took 11.4s
Alice's allocation: 100.00 CKBytes   Bob's allocation: 100.00 CKBytes
Alice's allocation:  90.00 CKBytes   Bob's allocation: 110.00 CKBytes
Alice's allocation:  95.00 CKBytes   Bob's allocation: 105.00 CKBytes
Alice Settle (finalize + Close)
==> Settle took 4.431s
EVENT Concluded  version=0 client=...
```

**`inspect`** (cells while the channel is open; after settling both lists are empty)
```
== Channel cell (type script = PCTS)
  out_point 0x9e31ec643fbb…:0  capacity 804.0 CKB  lock=PCLS  type=PCTS  data=256 byte
== Funds cell (lock script = PFLS)
  out_point 0x93c7990362f1…:0  capacity 100.0 CKB  lock=PFLS  type=-  data=0 byte
  out_point 0x9e31ec643fbb…:1  capacity 100.0 CKB  lock=PFLS  type=-  data=0 byte
```

**`sudt`**
```
Alice's allocation: 50.00 SUDT   Bob's allocation: 50.00 SUDT
Alice's allocation: 40.00 SUDT   Bob's allocation: 60.00 SUDT
Wallet after  | Alice: 50121.00 CKByte  99999990 SUDT
Wallet after  | Bob:   50154.00 CKByte  100000010 SUDT
```

**`dispute`**
```
Alice ForceSettle, NO finalize (simulating an uncooperative Bob)
Dispute called
EVENT Registered version=2 timeout=<Timeout: ...> client=028ef3...   (Bob's watcher)
EVENT Registered version=2 timeout=<Timeout: ...> client=028510...   (Alice's watcher)
Dispute is on-chain (phase Registered) but no Registered event arrived within 30s (perun-ckb-backend subscription race); calling Settle again to ForceClose
ForceClose called
==> Force settle took 34.042s
EVENT Concluded  version=0 client=...
```

### Step B5. Stop the devnet when done
```sh
./run.sh lab05-devnet stop
```

### Part B exercises
**B.1. Check the balances (`default`).** Work out the in-channel balances after each step and compare with `=== Allocation Balances ===` in the log. Compare the wallet balances before and after: why does Alice not gain exactly 5 CKB?

**B.2. Inspect the cells (`inspect`).** From the `cells.sh` output:
- How many channel cells and funds cells exist while the channel is open? What lock/type does each have?
- The channel cell holds 804 CKB but only 256 bytes of data. Where does that capacity come from, and who gets it back? Why does each funds cell hold exactly 100 CKB?
- One funds cell shares its transaction hash with the channel cell (`…:1` vs `…:0`). Which transaction created it, Start or Fund? Whose deposit is the other one?
- Do these cells still exist after settling?

**B.3. SUDT (`sudt`).** Each side deposits 50 SUDT, Alice pays 10 SUDT. Check the SUDT balances in the channel and in the wallets. Which type script do the in-channel cells holding SUDT have?

**B.4. Dispute (`dispute`).** Both watchers receive `Registered` with version 2. Why version 2, when there were two payments after opening? `Force settle took` reports about 34 s with `challengeDuration = 10` (`code/ckb/client/client.go` line 190): how much of that is the challenge duration and how much is the workaround's timeout? Relate the ForceClose to A.4.

**B.5. (Code reading) The race.** Read `run()` and `emitEventIfNecessary()` in `$(go env GOMODCACHE)/perun.network/perun-ckb-backend@v0.2.1-0.20250603085027-d25a5eb09110/channel/adjudicator/adjudicator_sub.go`, and `ensureRegistered` in go-perun's `client/adjudicate.go`. Why do the watchers (subscribed when the channel opened) see `Registered`, while the subscription that `Settle` creates does not? Read `ForceSettle` in `code/ckb/client/channel.go` and explain why calling `Settle` a second time works.

---

## Check questions
1. Why does CKB need **two different lock scripts**, one for the channel cell and one for the funds cells?
2. Where does the channel's on-chain state live: cell data, args, or witness?
3. Why does a channel on CKB need **several cells** instead of one contract as on Ethereum?
4. Who pays the capacity for the channel cell, and who gets it back when the channel closes?
5. To use a Perun channel in a real CKB dApp, what would you replace `devnet` with?

<details>
<summary>Hints</summary>

- **A.1:** `get_channel_action()` → `Start` → `check_valid_start`; `Progress` → load witness → `check_valid_progress` (Fund | Dispute); `Close` → load witness → `check_valid_close` (Abort | Close | ForceClose).
- **A.2:** funds sit in **funds cells** locked by PFLS. The registered state sits in the channel cell's **cell data** (`ChannelStatus`). Register/dispute is a transaction that consumes the old channel cell and creates a new one with a `Dispute` witness. Conclude is a transaction with a `Close` or `ForceClose` witness. Uniqueness comes from the thread token (OutPoint). Time is based on `since` and block headers.
- **A.3:** an OutPoint can only be spent **once**, so a duplicate channel cannot be created. Type ID uses the same idea.
- **A.4:** `verify_time_lock_expired(challenge_duration)` checks the input's `since` field: a ForceClose transaction is only valid once the challenge duration has passed since the channel cell was created.
- **B.1:** in the channel: 100/100 → 90/110 → 95/105. Alice's wallet: −100 deposited, +95 returned, minus transaction fees and the capacity temporarily locked in cells.
- **B.2:** 1 channel cell (lock PCLS, type PCTS) and 2 funds cells (lock PFLS), one per party. The channel cell's capacity is paid by Alice (the opener) to store the channel state on-chain and is returned to her on close. Each funds cell holds exactly the 100 CKB deposit. The funds cell with the same tx hash as the channel cell is Bob's: the Fund transaction consumes the channel cell, recreates it with `funded = true` and adds Bob's funds cell in the same transaction. The other one is Alice's, created by the Start transaction. After settling all three are consumed.
- **B.4:** version 0 is the opening state, then two payments give version 2. About 10 s is the challenge duration; the rest is the workaround's 30 s timeout (challenge + 20 s) plus the ForceClose transaction. ForceClose passes because by then `since` shows the challenge duration has elapsed.
- **B.5:** a subscription only emits `Registered` when it has already seen the funded status and then sees a disputed one. The watchers' subscriptions have been polling since the channel opened, so they see that transition. The subscription created inside `Settle` first polls 4 s later, when the Dispute is already committed, so it records the disputed status as its starting point and emits nothing. `registerDispute` has already set the phase to `Registered`, so the second `Settle` skips waiting (`ensureRegistered` returns immediately) and calls `Withdraw` → ForceClose.

Check questions:
1. The channel cell must be updatable by either party. The funds cells must move together with the channel cell, so funds cannot be withdrawn on their own.
2. `ChannelStatus` is in the channel cell's **data**. `ChannelConstants` is in the type script **args**. Signatures and the action type are in the **witness**.
3. CKB has no shared contract storage. State must live in cells. Each party needs its own funds cell to deposit independently, and the channel cell holds the shared state.
4. The party opening the channel creates the channel cell and pays its capacity, and gets it back on close (see `payment_min_capacity` in `Participant` and `check_valid_close`).
5. Testnet or mainnet, using the contracts deployed there. The deployment config must point to the correct `code_hash` and `cell_dep`.
</details>

## Checklist
- [ ] Part A: PCTS branch diagram, ETH ↔ CKB mapping table
- [ ] `./run.sh check` shows ✔ for everything under Lab 05
- [ ] `./run.sh lab05-devnet start` reports "Contracts: deployed"
- [ ] Ran `default`, `inspect`, `sudt`
- [ ] (Advanced) Ran `dispute` and saw `Registered`

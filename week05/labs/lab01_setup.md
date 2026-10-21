# Lab 01 — Set up the server environment

## Goals
- Know what the server has and what it lacks.
- Install Hardhat, test-build the example and the lab code.
- Understand the layout of `perun-examples` and of the lab code in `week05/code`.

This lab only prepares Labs 02–04 (Ethereum). The CKB tooling is installed in Lab 05.

## Steps

### Step 1. Set environment variables (once)
Open a terminal in VS Code (`` Ctrl+` ``):
```sh
cat >> ~/.bashrc <<'EOF'
export PERUN=/root/perun-examples
export LABS=/root/ckbuilder/week05
EOF
source ~/.bashrc
```

### Step 2. Install Hardhat (if missing)
```sh
cd $PERUN/payment-channel && npm install
```
Takes about 30 seconds and creates about 290 MB of `node_modules`. The warning `npm WARN deprecated uuid@8.3.2` is normal.

### Step 3. Check the whole environment
```sh
cd $LABS/code
./run.sh check
```
The script prints the state of each tool and the free disk space, test-builds `payment-channel` and `code/eth`, then prints library versions and the list of examples (the answers to exercises 1–3). The first build takes 1–3 minutes.

Output on the server (2026-10-21), Labs 01–04 part:
```
== System
  ✔ OS                           Ubuntu 26.04.1 LTS aarch64
  ✔ CPU / RAM                    2 cores / 3.3Gi
  ✔ Disk / free                  29 GB
== Labs 01–04 (Ethereum)
  ✔ go                           go version go1.26.8 linux/arm64
  ✔ node                         v22.22.1
  ✔ npm                          9.2.0
  ✔ perun-examples               /root/perun-examples (6e975e7 Merge pull request #51 ...)
  ✔ hardhat                      2.29.1
== Test build
  ✔ go build                     /root/perun-examples/payment-channel
  ✔ go build                     /root/ckbuilder/week05/code/eth
```
The `✘` lines under Lab 05 (ckb, tmux, riscv…) are expected. Those tools are installed in Lab 05.

### Step 4. See the available commands
```sh
./run.sh help
```

## Where the lab code lives
All exercises are already coded in `week05/code/`. **You do not need to edit `/root/perun-examples`.**

```
week05/code/
├── run.sh              # the only command you need: ./run.sh <lab> [scenario]
├── eth/                # Labs 02–04, Go module based on payment-channel
│   ├── client/         #   extended client: Balances, Version, MintForMe, Steal,
│   │                   #   ForceSettle, Policy, Buggy flag to restore the original bug
│   ├── demo/           #   chain setup, contract deployment, Alice and Bob
│   └── cmd/lab02 lab03 lab04
├── ckb/                # Lab 05, based on payment-channel-ckb + install and devnet scripts
├── tools/              # check.sh, read_code.sh, txsummary.py (counts on-chain transactions)
└── logs/               # one log per run (created automatically)
```

For every Ethereum scenario, `run.sh` does four things:
1. Stops any old Hardhat holding port 8545.
2. Starts a fresh Hardhat, so balances always start at 10/10.
3. Runs the Go scenario and saves the log to `code/logs/`.
4. Reads the Hardhat log, prints a table of on-chain transactions, then stops Hardhat.

So **you do not need two terminals**.

## Exercises
All answers come from the output of `./run.sh check`:
1. Which versions of `go-perun`, `perun-eth-backend` and `go-ethereum` does `payment-channel` use?
2. Which `go-perun` version does `payment-channel-ckb` use? What is the `replace` line for?
3. Which blockchain does each `payment-channel-*` directory target?
4. Open `$PERUN/payment-channel/hardhat.config.js`. Which account belongs to the Deployer, Alice and Bob? (Compare with the constants in `code/eth/demo/demo.go`.)

## Check questions
1. Why is go-perun a separate module while the Ethereum part lives in `perun-eth-backend`?
2. `hardhat.config.js` sets `mining: { auto: true, interval: 5000 }`. What do the two parameters mean?

<details>
<summary>Hints</summary>

- **Exercise 1:** go-perun v0.15.0, perun-eth-backend v0.6.0, go-ethereum v1.15.11.
- **Exercise 2:** pseudo-version `v0.13.1-0.20250528124331-…`, older than v0.15.0. `replace` points `ckb-sdk-go` to perun-network's fork, because the backend needs changes not yet in upstream.
- **Exercise 3:** `-ckb` CKB, `-ckb-eth` CKB ↔ Ethereum, `-xlm` Stellar, `-dot` Polkadot, `-icp` Internet Computer, `-cc` cross-contract, no suffix is Ethereum.
- **Exercise 4:** Account #0 `0xe84d…` = Deployer, #1 `0x56FD…` = Alice, #2 `0x6536…` = Bob, 10 ETH each.
1. go-perun is **blockchain-agnostic**. The protocol core knows nothing about Ethereum. Each chain has a backend implementing the `Funder`, `Adjudicator` and `Wallet` interfaces.
2. `auto: true` mines immediately whenever a transaction arrives. `interval: 5000` also mines an empty block every 5 seconds, so chain time keeps moving (needed for the challenge duration in Lab 04).
</details>

## Checklist
- [ ] `./run.sh check` shows ✔ for everything under Labs 01–04
- [ ] Answered the 4 exercises

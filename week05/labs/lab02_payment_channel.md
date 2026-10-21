# Lab 02 — First payment channel & client architecture

## Goals
- Run the Alice ↔ Bob demo on a local Ethereum chain (Hardhat).
- See the channel lifecycle: **open → fund → off-chain update → finalize → settle**.
- Know which actions go on-chain and which do not.
- Understand the client components: Funder, Adjudicator, Watcher, Wire Bus, Wallet.

## Background
- Only **opening** and **settling** touch the blockchain. The payments in between are just signatures exchanged by the two parties.
- Every state has an increasing `version`. The highest-version state signed by both parties is the valid one.

## Part A — Run the demo

### Step 1. Run
```sh
cd $LABS/code
./run.sh lab02
```
The scenario (`code/eth/cmd/lab02/main.go`) is identical to the original `payment-channel/main.go`: Alice opens a 5 ETH channel, Alice pays 3, Bob pays 1, Alice pays 1, then they settle. The difference is that after every step it prints the in-channel balances, the version and the elapsed time.

### Step 2. Read the result
Output on the server (2026-10-21):
```
Wallets before  Alice=10.000000 ETH  Bob=10.000000 ETH

| Step | Action | Alice (channel) | Bob (channel) | Version | IsFinal | Time |
|---|---|---|---|---|---|---|
| 0 | OpenChannel(bob, 5)          |  5.00 |  0.00 | 0 | false |   2.201s |
| 1 | chAlice.SendPayment(3)       |  2.00 |  3.00 | 1 | false |    376ms |
| 2 | chBob.SendPayment(1)         |  3.00 |  2.00 | 2 | false |    846ms |
| 3 | chAlice.SendPayment(1)       |  2.00 |  3.00 | 3 | false |    495ms |
| 4 | Finalize (IsFinal = true)    |  2.00 |  3.00 | 4 | true  |    522ms |

[Bob] EVENT Concluded   version=4
[Alice] EVENT Concluded   version=4
Settle (conclude + 2 withdraw) took 85ms
Wallets after   Alice=6.999594 ETH  Bob=12.999939 ETH

On-chain transactions (from the Hardhat log):
| # | Block | Sender | Target | Value | Gas | Kind |
|---|---|---|---|---|---|---|
| 1 | #1 | Deployer | Adjudicator | 0 ETH | 2977006 | deploy Adjudicator |
| 2 | #2 | Deployer | AssetHolder | 0 ETH | 941410 | deploy AssetHolder |
| 3 | #4 | Alice | AssetHolder | 5 ETH | 45801 | deposit |
| 4 | #5 | Alice | Adjudicator | 0 ETH | 170916 | register / conclude |
| 5 | #6 | Alice | AssetHolder | 0 ETH | 42816 | withdraw |
| 6 | #7 | Bob | AssetHolder | 0 ETH | 42836 | withdraw |
Total: 6 transactions.
```
Your timings will differ a little. The full logs are in `code/logs/lab02-run-*.log` and `code/logs/hardhat-lab02-run-*.log`.

## Part B — Read the client architecture

### Step 3. Print the code to read
```sh
./run.sh lab02-code
```
The script prints:
- the lines that create the Funder, Depositor, Adjudicator, Watcher and `client.New` in `payment-channel/client/client.go`,
- the three interfaces `ProposalHandler`, `UpdateHandler`, `AdjudicatorEventHandler` in go-perun,
- the steps of `Update` (sign, send, collect the signature),
- the list of channel phases in `channel/machine.go`.

Each snippet comes with its file path. Ctrl+click in the VS Code terminal to open the file.

## Exercises
Answer from the output of Steps 2 and 3, in `$LABS/notes/lab02.md`:
1. Explain each row of the lifecycle table: why does the version go up by 1 after every step, even the Finalize step that moves no money?
2. There were 3 payments but none of them shows up in the on-chain table. If Alice and Bob paid each other 3 times directly on-chain, how many transactions would that take? How many does the channel need, excluding the 2 deployments?
3. Why does Bob have no deposit transaction? Why do both Alice and Bob have a withdraw transaction?
4. Alice ends with 6.999594 ETH, not exactly 7. Where did the difference go? How much less does Bob have, and why is it less than Alice's difference?
5. Fill in the component table (using the `lab02-code` output):

| Component | Talks to | Role |
|---|---|---|
| Contract backend | ? | ? |
| Funder + Depositor | ? | ? |
| Adjudicator | ? | ? |
| Watcher | ? | ? |
| Wire bus | ? | ? |
| Wallet | ? | ? |

6. Put the steps of one `SendPayment` in order (using `updateGeneric` in the `lab02-code` output): *Bob checks and signs*, *Alice signs the new state*, *Alice sends the update over the bus*, *Alice applies the new state to the state machine*, *Alice receives Bob's signature and enables the new state*.
7. Which phases in the list need an on-chain transaction?

## Check questions
1. Why does Bob only call `AcceptedChannel()` and never `OpenChannel`?
2. Where does the watcher run? If Alice shuts down her machine, does her watcher still protect her?
3. What kinds of messages travel over the wire bus? Does any of them reach the blockchain?
4. `challengeDuration` is 10 seconds. Why does production need a much larger value?

<details>
<summary>Hints</summary>

1. Every off-chain update, even one that only sets `IsFinal`, is a new state both parties must sign, so the version always increases.
2. Directly on-chain: 3 transactions, each paying fees and waiting for a block. Channel: 4 transactions (1 deposit, 1 conclude, 2 withdraw) **whether there are 3 or 3000 payments**.
3. Bob deposits 0 ETH, so there is nothing to deposit. Each party withdraws its own share from the AssetHolder.
4. Gas for deposit, conclude and withdraw (about 0.0004 ETH). Bob only pays gas for one withdraw (about 0.00006 ETH).
5. Contract backend: the Ethereum node at `ws://127.0.0.1:8545`; signs and sends transactions. Funder + Depositor: the AssetHolder; deposit ETH when opening. Adjudicator: the Adjudicator contract; register/conclude/withdraw. Watcher: listens for on-chain events and refutes old states. Wire bus: the peer's client over libp2p; carries proposals/updates/signatures. Wallet: signs states and transactions.
6. Apply the new state to the machine (`machine.Update`) → Alice signs (`machine.Sig`) → send over the bus (`conn.Send`) → Bob checks (`HandleUpdate`) and signs → Alice receives the signature (`AddSig`) and enables the state (`EnableUpdate`/`EnableFinal`).
7. `Funding` (deposit), and the closing or dispute phases: `Registering`/`Registered`, `Progressing`/`Progressed` (app channels), `Withdrawing`/`Withdrawn`. `Acting`, `Signing` and `Final` are purely off-chain.

Check questions:
1. Alice is the proposer. Bob receives the proposal in `HandleProposal`, accepts it and pushes the channel into the Go channel `c.channels`. `AcceptedChannel()` just takes it out.
2. `local.NewWatcher` runs inside the client's own process. If Alice shuts down, so does her watcher. That is why the challenge duration must be long enough, or a remote watcher is needed.
3. Proposals, accept/reject, updates, signatures. None of them reaches the chain.
4. It is the window for the other party to object after someone registers a state on-chain. In practice it should be hours to days, so a cheated party has time to come back online.
</details>

## Checklist
- [ ] `./run.sh lab02` finishes with two `EVENT Concluded`
- [ ] Answered exercises 1–4 from the lifecycle and transaction tables
- [ ] Filled in the component table and the order of `SendPayment` steps

# Lab 04 — On-chain dispute

## Goals
- See the **uncooperative** settle path: register the state on-chain, wait out the challenge duration, then conclude.
- Tell the adjudicator events apart: `Registered`, `Progressed`, `Concluded`.
- Measure the cost (time, number of transactions) of the dispute path compared with the fast path.

## Background
In go-perun, `Channel.Settle(ctx, secondary)` works like this:
- If the state is **already final**: conclude immediately. This is the fast path, used in Lab 02.
- If the state is **not final**: register the latest state on-chain, **wait out the challenge duration**, then conclude.

While waiting, the other party's watcher can submit a higher-version state if the registered one is outdated. Hardhat mines a block every 5 seconds, so chain time moves in 5-second steps.

## Prepared code
| File | Contents |
|---|---|
| `code/eth/client/channel.go` | `Settle` (line 70): finalize, then settle. `ForceSettle` (line 82): settle **without** finalizing. `WithdrawOnly` (line 92): `Settle(ctx, true)`, withdraw only |
| `code/eth/client/handle.go` | `HandleAdjudicatorEvent` (line 111) logs each event type in detail |
| `code/eth/cmd/lab04/main.go` | `-mode fast` or `-mode dispute`, `-challenge <seconds>` |

Read the go-perun code:
```sh
GP=$(go env GOMODCACHE)/perun.network/go-perun@v0.15.0
sed -n '/^func (c \*Channel) Settle/,/^}/p' $GP/client/adjudicate.go
```

## Run
Run all three (fast path, dispute with 10 seconds, dispute with 30 seconds) and compare, about 1.5 minutes:
```sh
cd $LABS/code
./run.sh lab04 all
```
Or one at a time:
```sh
./run.sh lab04 fast
./run.sh lab04 dispute        # challenge 10 seconds
./run.sh lab04 dispute 30     # challenge 30 seconds
./run.sh lab04 dispute 60     # try another value
```

## Result (run on the server on 2026-10-21)
The "Comparing the 3 runs" section at the end of the output:
```
--- lab04-fast-...log
[Bob] EVENT Concluded   version=4
[Alice] EVENT Concluded   version=4
==> Alice settle took 386ms
Total: 6 transactions. ... deposit: 1, register / conclude: 1, withdraw: 2

--- lab04-dispute10-...log
[Bob] EVENT Registered  version=3 timeout=<Block timeout: 1792577659>
[Alice] EVENT Registered  version=3 timeout=<Block timeout: 1792577659>
[Bob] EVENT Concluded   version=3
[Alice] EVENT Concluded   version=3
==> Alice force settle took 19.678s (challenge duration 10s, block time 5s)
Total: 7 transactions. ... deposit: 1, register / conclude: 2, withdraw: 2

--- lab04-dispute30-...log
[Bob] EVENT Registered  version=3 ...
[Bob] EVENT Concluded   version=3
==> Alice force settle took 34.421s (challenge duration 30s, block time 5s)
Total: 7 transactions. ... register / conclude: 2, withdraw: 2
```
Timings vary between runs: across several runs on this server, force settle took 15–20 s with a 10 s challenge and 34–40 s with a 30 s challenge. In the transaction table, the register and conclude transactions are several blocks apart, with only empty blocks in between.

## Exercises
1. Write down the event order for each client in the three runs. Which event does the dispute path have that the fast path does not?
2. The version in the fast-path events is 4, in the dispute-path events it is 3. Why do they differ by 1?
3. With a 10 s challenge, force settle took 15–20 s. With 30 s, it took 34–40 s. Explain the extra time beyond the challenge duration, relating it to the 5-second block time.
4. The dispute path has 7 transactions, the fast path 6. What is the extra transaction? In the transaction table, what are the two `register / conclude` rows of the dispute path?
5. Run `./run.sh lab04 dispute 60`, predicting the force settle time first.
6. **(Code reading)** If Alice registers an **old** state, who reacts? Read `$GP/watcher/local/watcher.go` and find where the watcher compares versions and registers again. Describe it in 3–4 sentences.

## Check questions
1. Why is the dispute path still **safe** for Bob even though he does nothing?
2. If Alice registers an old state, e.g. version 1 when she still had more money, who reacts and how?
3. Why do real applications always try to take the fast path?

<details>
<summary>Hints</summary>

1. The dispute path has `Registered` before `Concluded`. The fast path only has `Concluded`.
2. The fast path has one extra off-chain update to set `IsFinal` (version 3 → 4). The dispute path registers version 3 directly, not final.
3. Waiting time ≈ challenge duration + waiting for the next block whose timestamp passes the timeout (up to 5 s) + a few transactions.
4. The extra transaction is `register`. The first row is `register` of state v3, the second is `conclude` after the challenge ends. On the fast path a single `concludeFinal` transaction does both.
5. Roughly 60–70 seconds.
6. The watcher listens for `Registered` events. If the on-chain version is lower than the newest version the client holds, it automatically sends a transaction registering the newer state within the challenge window.

Check questions:
1. A registered state must carry both parties' signatures. Bob holds the signatures of the latest state, and Bob's watcher monitors the chain.
2. Bob's watcher sees a `Registered` event with a version older than the state Bob holds. It automatically registers the newer state within the challenge window. The highest version wins.
3. The dispute path is slow (it waits out the challenge duration) and costs extra transactions and gas.
</details>

## Further reading (optional)

### App channel: Tic-Tac-Toe
A general state channel, where the state carries application data, not just balances. Source in `$PERUN/app-channel`:

| File | What to understand |
|---|---|
| `app/app.go` | `ValidInit`, `ValidTransition`: game rules checked off-chain |
| `app/util.go` | `CheckFinal`. The "Check all set" loop uses `v != notSet`: does that really mean "every cell is filled"? |
| `client/channel.go` | `Set` uses `ch.Update`, `ForceSet` uses `ch.ForceUpdate`: advances the state **on-chain** when the opponent stops cooperating (emits `Progressed`) |

### Virtual channels
go-perun v0.15.0 **does include** virtual channels (`ls $GP/client/virtual_channel*.go`). Alice and Bob have no direct channel; each has a channel with Ingrid. The Alice–Bob virtual channel is funded off-chain from the two parent channels. Ingrid is only involved when opening and closing.

## Checklist
- [ ] `./run.sh lab04 all` finishes, showing `Registered` then `Concluded`
- [ ] Can explain the force settle time for both challenge durations
- [ ] Can explain the extra transaction on the dispute path

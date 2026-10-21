# Lab 03 — Validation, bug hunting & two-sided funding

## Goals
- Understand **who checks what** on an update: the go-perun state machine or the application logic.
- Attack your own channel and watch the protection layers work.
- See a real bug in the example code: the missing `return` after `Reject`.
- Let both parties deposit when opening, under a policy Bob sets himself.

## Background
go-perun enforces the general rules in `channel/machine.go`: the total is conserved, the version increases, a final state cannot change.
Application logic, such as "only accept updates that increase my balance", lives in `HandleUpdate` and `HandleProposal`, written by the application.

## Prepared code
| File | Contents |
|---|---|
| `code/eth/client/channel.go` | `MintForMe` (line 100), `Steal` (line 110): the two attack functions |
| `code/eth/client/handle.go` | `HandleProposal`, `HandleUpdate` with the `return` after `Reject` added. Setting `Buggy = true` restores the original behavior |
| `code/eth/client/client.go` | `Policy`: `MaxPeerFunding` (maximum deposit), `MinPayment` (anti-spam). `OpenChannel` takes `peerAmount` |
| `code/eth/cmd/lab03/main.go` | 10 scenarios |

Compare with the original code:
```sh
diff $PERUN/payment-channel/client/handle.go $LABS/code/eth/client/handle.go
```

## Run
Run one scenario at a time (recommended, so you can read each carefully):
```sh
cd $LABS/code
./run.sh lab03 mint
```
or all 10 in a row (about 2 minutes):
```sh
./run.sh lab03 all
```
List the scenarios: `./run.sh lab03`.

## Results and questions per scenario
The output below was produced on the server on 2026-10-21 (only the important lines are kept).

### 1. `mint`: money out of thin air
```
In channel after open        Alice=5.0000 Bob=0.0000 version=0
==> MintForMe(100) returned: updating machine: unequal allocation: invalid state transition
In channel after MintForMe   Alice=5.0000 Bob=0.0000 version=0
```
**Question:** does the error come from Alice's side or Bob's? Did Bob ever receive this update? (Hint: there is no `[Bob] REJECT` line in the log.)

### 2. `steal-buggy`: taking Bob's money, Bob runs the original code
```
In channel after SendPayment(3) Alice=2.0000 Bob=3.0000 version=1
[Bob] REJECT update v2: invalid balance: 1000000000000000000
panic: multiple calls on channel update responder
>>> Program exited with code 1 (see the log above)
```
**Question:** the total is conserved, so who blocked the update? Bob rejected correctly, so why did the program still crash? Read `HandleUpdate` in `code/eth/client/handle.go`: with `Buggy = true`, which line runs after `Reject`?

### 3. `steal`: same, with the fixed code
```
[Bob] REJECT update v2: invalid balance: 1000000000000000000
==> Steal(2) returned: channel update rejected by peer: invalid balance: 1000000000000000000
In channel after Steal   Alice=2.0000 Bob=3.0000 version=1
==> Program ran to the end, no panic.
```
**Question:** what is the version after `Steal`? Why is it not 2?

### 4. `fund52-buggy`: Alice asks Bob to deposit 2 ETH, original policy, original code
```
[Bob] REJECT proposal: peer funding too high: 2 ETH > 0 ETH
panic: runtime error: invalid memory address or nil pointer dereference
```
**Question:** what is the value of `lcp` when the proposal is rejected? Which line causes the nil pointer?

### 5. `fund52-default`: same, with the fixed code
```
[Bob] REJECT proposal: peer funding too high: 2 ETH > 0 ETH
==> Alice got an error opening the channel: channel proposal rejected by peer: peer funding too high: 2 ETH > 0 ETH
```
**Question:** what is Bob's original policy? Find the matching condition in `$PERUN/payment-channel/client/handle.go`.

### 6. `fund52`: deposit 5/2, new policy (Bob deposits at most 3 ETH)
```
In channel after open        Alice=5.0000 Bob=2.0000 version=0
In channel after 3 payments  Alice=2.0000 Bob=5.0000 version=3
Wallets after   Alice=6.999630 ETH  Bob=12.999866 ETH
Total: 7 transactions. ... deposit: 2, register / conclude: 1, withdraw: 2
```
**Question:** work out the in-channel and wallet balances yourself, then compare with the output. Why does Bob's wallet still end near 13 ETH, as in Lab 02, even though he deposited 2 more? Which transaction is new compared with Lab 02?

### 7. `fund54`: deposit 5/4, above the policy
```
[Bob] REJECT proposal: peer funding too high: 4 ETH > 3 ETH
==> Alice got an error opening the channel: ...
Total: 2 transactions. deploy Adjudicator: 1, deploy AssetHolder: 1
```
**Question:** why is there no deposit transaction at all?

### 8. `spam`: Bob rejects payments below 0.001 ETH
```
[Bob] REJECT update v1: payment too small: 0.0001 ETH < 0.001 ETH
==> SendPayment(0.0001) returned: channel update rejected by peer: payment too small ...
==> SendPayment(0.01)   returned: <nil>
```
**Question:** the Finalize update also goes through Bob's `HandleUpdate`. Why does the anti-spam rule not block it? (Read how `received` is computed in `HandleUpdate`.)

### 9. `afterfinal`: a payment after the state is final
```
In channel after Finalize   Alice=2.0000 Bob=3.0000 version=2 final=true
==> SendPayment(1) after final returned: updating machine: not in correct phase: invalid channel phase transition (... current: Final ...)
```
**Question:** which layer blocks it, go-perun or the application?

### 10. `micro`: 100 payments of 0.01 ETH
```
In channel after 100 payments  Alice=4.0000 Bob=1.0000 version=100
==> 100 off-chain payments took 42.718s (427.18 ms each)
==> If each were one 5-second on-chain block: at least 8m20s
Total: 6 transactions.
```
**Question:** 100 payments but only 6 on-chain transactions, exactly as with 3 payments in Lab 02. Explain. (The ~400 ms per payment on this server is mostly libp2p latency between the two clients plus persistence writes. Other machines will show different numbers.)

## Check questions
1. Distinguish the two layers of checks: which does go-perun guarantee, which does the application write itself? Sort the 10 scenarios into the two layers.
2. If Bob has a bug and accidentally signs a state that hurts him, can it be undone? Why?
3. How does `FundingAgreement` differ from `InitBals`?
4. Why must Bob check the proposal himself instead of trusting Alice?

<details>
<summary>Hints</summary>

1. **mint:** the error happens on Alice's side. The state machine checks the total (`unequal allocation`) before sending, so Bob never receives it.
2. **steal-buggy:** Bob's `HandleUpdate` blocks it because Bob's balance would drop. But the original code has no `return`, so after `Reject` it still calls `r.Accept`. The responder is used twice, `Accept` returns an error and `panic(err)` crashes the whole program, Alice included, since both run in the same process.
3. **steal:** the version stays 1. The update was rejected, so no new state was signed.
4. **fund52-buggy:** `lcp` is `nil` (the check function returns `nil, err`). The line `lcp.Accept(...)` causes the nil pointer.
5. **fund52-default:** Bob only accepts if his own deposit is 0: reject when `FundingAgreement[0][1].Cmp(big.NewInt(0)) != 0`.
6. **fund52:** in the channel, 5−3+1−1 = 2 for Alice and 2+3−1+1 = 5 for Bob. Alice's wallet 10−5+2 = 7, Bob's 10−2+5 = 13, minus gas. Bob deposits 2 and gets it back on settle; the money is only locked temporarily. The new transaction is Bob's deposit.
7. **fund54:** the proposal is rejected off-chain, before funding, so nobody deposits.
8. **spam:** the Finalize update moves no money, so `received = 0`. The rule only applies when `received > 0`.
9. **afterfinal:** go-perun blocks it. The channel is in the `Final` phase and accepts no more updates.
10. **micro:** on-chain transactions only happen on open (deposit) and close (conclude, withdraw). The number of payments in between does not matter.

Check questions:
1. go-perun: total conserved, version increases, no change to a final state, valid signatures (mint, afterfinal). Application: who may pay whom, deposit limits, minimum payment (steal, fund*, spam).
2. No. Bob's signature on that state is valid, and the highest version wins a dispute. That is why `HandleUpdate` is Bob's most important line of defense.
3. `InitBals` is the starting balance in the channel. `FundingAgreement` is how much each party commits to deposit on-chain. In this example they are equal.
4. A proposal is an off-chain message created by Alice. If Bob does not check it, Alice could make Bob lock a large amount or use an unknown asset.
</details>

## Checklist
- [ ] Ran all 10 scenarios (`./run.sh lab03 all`)
- [ ] Can explain both panics in `steal-buggy` and `fund52-buggy`
- [ ] Can point to the added `return` lines in `code/eth/client/handle.go`
- [ ] Answered the question for each scenario

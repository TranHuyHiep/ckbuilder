# CKB Weekly Report - Week 10

**Reporting period:** 19 - 25 November 2026  
**Publication date:** 24/11/2026
**Participant:** Tran Huy Hiep  

---

## 1. Summary

Week 9 put the bot on mainnet but left two gaps: I could not see how it was used, and some questions a signer really asks were not answered. This week I added five small features, each checked against real mainnet data. Every number in the replies still comes from code, not from an LLM.

**What was added**
- **A. Usage metrics and `/stats`:** one JSON line per message (command, language, ok or error, time). No message text, no addresses, and the Telegram id is replaced by a salted hash. Files are kept 14 days. `/stats` works only for the ids in `ADMIN_IDS`.
- **B. Cache:** repeated lookups are served from memory. A transaction is cached for 1 hour only once it has at least 24 confirmations; an address or DAO lookup for 30 seconds; a token for 60 seconds. Errors are never cached and parallel identical requests share one chain call.
- **D. `/tx <hash> <address>`:** adds a "Your view" block: what that transaction did to that address, in the same words `/sign` uses. An address that is not in the transaction is told so, instead of showing 0.
- **C. `/sign` reads `since` and checks new recipients:** an input locked by a block, epoch or time condition is reported as met, not yet met (yellow) or not checkable. A payment of 100 CKB or more to a lock with no transaction history gets a yellow warning to check the address.
- **E. `/dao` timing:** for each deposit, the epoch at which the current 180-epoch cycle ends, and what that means: request the withdrawal by then to claim at that epoch, later and the claim moves a full cycle (180 epochs) further. No future compensation is predicted.
- **Tests:** 435 → 488, all offline. The service on the server was restarted with the new code and runs on mainnet. I checked the new commands in Telegram on the live bot (section 2, cases 9 to 13, 7 screenshots).

---

## 2. Evidence

| # | Case | Result | Evidence |
|---|---|---|---|
| 1 | All tests and type check | 488 passed, typecheck clean | [`evidence/0_tests_488.txt`](./evidence/0_tests_488.txt) |
| 2 | A: handler run against mainnet with 3 users, then `/stats` | Admin gets 7 messages from 3 users, commands, languages, 1 usage error, p50 727 ms / p90 1861 ms, cache 1 hit / 3 misses. A non-admin gets the help text. The raw event file holds only time, hashed user, kind, command, language, ok, error code and milliseconds | [`evidence/A_stats_local_run.txt`](./evidence/A_stats_local_run.txt) |
| 3 | B: the same tools with the cache off and on, real mainnet | Repeat of `/tx`: 930 ms → 0 ms; `/addr`: 336 ms → 0 ms; `/dao`: 339 ms → 0 ms. Outputs identical as JSON | [`evidence/B_cache_latency.txt`](./evidence/B_cache_latency.txt) |
| 4 | D: `/tx` of a 2021 mainnet transaction as sender, receiver and an unrelated address | Sender: `You pay 800.00002228 CKB in total, of which 0.00002228 CKB is the network fee`. Receiver: `You receive 800 CKB net` plus two tokens. Unrelated: `This address is not part of this transaction`. Each number equals that party's line in the base summary | [`evidence/D_perspective_mainnet.txt`](./evidence/D_perspective_mainnet.txt) |
| 5 | C: `npm run check:sign` after the changes | **20/20 match** between the unsigned and the committed path | [`evidence/C_check_sign.txt`](./evidence/C_check_sign.txt) |
| 6 | C: the 20 fixture transactions through the new checks | 3 of 20 have an input with `since` (all Nervos DAO claims, already met). **No new finding on any of the 20** | [`evidence/C_since_recipient_mainnet.txt`](./evidence/C_since_recipient_mainnet.txt), part 1 |
| 7 | C: real live mainnet cell spent in an unsigned transaction (nothing signed or sent) | 500 CKB to a lock with no history: yellow `no transaction history`. 500 CKB to an address with history: no finding. Absolute block `since` in the future: yellow `time-locked … from block …`. In the past: info `already met`. Relative 1,000,000 blocks: yellow | same file, part 2; same transactions through the bot in en and zh: [`evidence/C_sign_examples_via_bot.txt`](./evidence/C_sign_examples_via_bot.txt) |
| 8 | E: `/dao` of an address with a 130000 CKB deposit, in en, zh and vi | Cycle ends at epoch 15485.37 (about 531 h), next claim epoch 15665.37. Equal to `ccc.calcDaoClaimEpoch` computed from the real headers | [`evidence/E_dao_timing_mainnet.txt`](./evidence/E_dao_timing_mainnet.txt) |
| 9 | Live bot in Telegram: `/tx <hash> <address>` as sender, then with an address that is not in the transaction (D) | Sender: `Your view (ckb1qzda…udafpe)`, `You pay 800.00002228 CKB in total, of which 0.00002228 CKB is the network fee`. Other address: `This address is not part of this transaction.` | [`telegram/01_tx_view_sender.png`](./telegram/01_tx_view_sender.png), [`02_tx_view_not_involved.png`](./telegram/02_tx_view_not_involved.png) |
| 10 | Live `/dao` in English and Vietnamese (E) | Cycle ends at epoch 15485.37, next claim epoch 15665.37, the same as case 8. The phase label reads `deposited` / `đã gửi` | [`telegram/03_dao_timing_en.png`](./telegram/03_dao_timing_en.png), [`05_dao_timing_vi.png`](./telegram/05_dao_timing_vi.png) |
| 11 | Live `/stats` as admin (A, B) | 19 messages from 1 user: tx 5, dao 4, sign 3; en 17, vi 2; 0 errors; p50 323 ms, p90 1155 ms; cache 5 hits, 4 misses (56%). Cache counters above 0 show the cache working in the live bot | [`telegram/04_stats.png`](./telegram/04_stats.png) |
| 12 | Live `/sign` with a 500 CKB payment to a lock with no history (C) | 🟡 `500 CKB goes to ckb1qzda0c…0szp0w, which has no transaction history on chain. Check the address…`, plus the expected cell dep warning (the hand-built transaction has no cell deps) | [`telegram/06_sign_new_recipient.png`](./telegram/06_sign_new_recipient.png) |
| 13 | Live `/sign` with an input locked until block 99999999 (C) | 🟡 `Input 0xf221d2b1…858e:1 is time-locked: it can only be spent from block 99999999. The network rejects this transaction until then.` | [`telegram/07_sign_since_locked.png`](./telegram/07_sign_since_locked.png) |


**What the evidence does not show**
- **Telegram coverage is partial.** Cases 9 to 13 went through the live bot. Not captured in Telegram: `/tx` with the receiver's address, `/sign` in Chinese or Vietnamese for the new findings, and the cache's speed (case 11 shows hits, not time). Those were checked only by calling the handler and tools directly (cases 3, 4 and 7).
- **A:** the numbers in case 2 come from a local run with 6 messages. Case 11 is the live bot, but with only one user (me) and 19 messages, so it shows that `/stats` works, not how people use the bot.
- **B:** the times are direct calls to the tools, not what a Telegram user sees. First-call times move a lot from run to run on the public RPC, so the useful number is the repeat. The 30 second TTL means a balance can be up to 30 seconds old.
- **C:** I found no real mainnet transaction with a non-DAO `since`. The future, past and relative cases in case 7 are built by me from a real live cell. Epoch and time rules are my approximation (see section 4); they were checked by unit tests I wrote, not against a node that actually verifies a transaction.
- **E:** only the epoch boundary is checked, against CCC. It is **not compared with Neuron** (not available). Compensation after the boundary is not predicted.
- **Thresholds are my own defaults, not calibrated:** 100 CKB and at most 5 recipients for the new-recipient check; a 10 minute margin for time-based `since`; the cache TTLs; 14 days retention.
- **LLM:** still not configured, so the 50 question eval has still produced no accuracy figure.

---

## 3. Key Learnings

- **A passing run can hide a feature that is off.** `ADMIN_IDS` was appended to a `.env` file that had no final newline, so it became part of the `CKB_NETWORK=mainnet` line. Everything still started, and `/stats` would have stayed disabled. I only noticed because I counted the matching lines and got 0. I fixed the line and I now check configuration by reading it back.
- **A real run found an old display bug.** In `/dao` the phase label was looked up with the key `deposit`, but the data says `deposited`. So the English-looking word `deposited` appeared in the Chinese and Vietnamese replies since week 9. The tests had never checked that label. It is fixed and covered in all three languages.
- **Reuse the same logic for the same question.** The "your side" view in `/tx` and in `/sign` now comes from one function. A test checks that both give the same number for the same address, so they cannot drift apart.
- **Say "not checked" instead of guessing.** For a relative `since` without the input's block header, a time too close to the tip, or a failed recipient lookup, the finding says it could not check. A false red or yellow would cost more trust than a missing one.
- **A cache needs a rule for what must not be cached.** A transaction with fewer than 24 confirmations, a failed call and a zero TTL are all tested to be re-read next time.
- **Keep the cached object clean.** The address view is added to a copy, so one user's address never leaks into another user's cached transaction.

---

## 4. Challenges

| Challenge | Solution |
|---|---|
| A `since` can count in blocks, epochs or time, absolute or relative, and the node's own rule for time uses a median of recent blocks | Blocks and epochs are compared exactly (epochs as fractions). For time I only decide when the gap to the tip is over 10 minutes and report "not checked" inside it. This is my approximation, not the node's code |
| Nervos DAO claims carry a `since` on purpose; warning on them would be a false alarm | A DAO input whose `since` is met adds no finding; one that is too early is still warned (both tested) |
| A "new recipient" lookup on an already committed transaction always finds that transaction | Real checks use unsigned transactions; the 20 committed fixtures are therefore not a test of this check, only of "no unexpected warning" |
| A shell command with `pkill -f` matched my own command line and killed it | Do not use `pkill -f` with a name that appears in the command itself |
| A script that uses the RPC client kept the process alive after printing its results | End scripts with an explicit exit |
| `systemctl restart` logs the old process as `status=143` (SIGTERM passed on by the Node wrapper) | Left as is: the bot stops after logging "SIGTERM received", and the new process starts. Cosmetic |
| Telegram token printed in a terminal during week 9 | Still pending on my side: rotate it in BotFather |

---

## 5. Limits and Next Steps

- Send the remaining rows of `docs/telegram-test.md` section G through the live bot (`/tx` with the receiver's address, `/sign` in zh and vi) and add screenshots.
- Let the bot collect a few days of real usage, then read `/stats` to decide whether the cache TTLs and the quota fit what people use.
- Look for a real mainnet transaction with a non-DAO `since` (timelock or multisig) to test case 7 on something I did not build myself.
- Chinese and Vietnamese texts still need a native speaker's review.
- Not done and out of scope this week: a fallback RPC, an alert when the bot is down, backups, an MCP server, address labels, and the LLM eval (needs a key).

---

## 6. Reflection

This week every feature was small, and the useful work was in the edges: what to say when the check cannot be done, what must not be cached, and what an address that is not in the transaction should see. The two real mistakes (the `.env` line and the `deposited` label) were again not found by tests but by looking at the actual result. I will keep reading back the real output, and I will count a feature as done only when I can say which part I have not verified.

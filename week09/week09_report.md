# CKB Weekly Report - Week 9

**Reporting period:** 12 - 18 November 2026  
**Publication date:** 16/11/2026
**Participant:** Tran Huy Hiep  

---

## Summary

Weeks 6 to 8 built the checks. This week I made the bot something that can run unattended and be used in more than one language, and I prepared a way to measure the free-text `/ask` path. Every number in the replies still comes from code, not from the LLM.

**What was added**
- **Token decimals in `/tx` and `/addr`:** token amounts now show with decimals and symbol (from the token's info cell), as `/sign` already did. When no info cell exists, the amount stays in raw units and says so.
- **Persistent quota:** the daily limits and each user's language are saved to `data/*.json`, written atomically. A restart no longer resets everyone's allowance. A corrupt file makes the bot start clean instead of crashing.
- **English, Chinese and Vietnamese:** all templates exist in en, zh and vi. `/lang en|zh|vi` sets the language per user; until then it follows the Telegram language of the user. `/sign` findings are translated from their code and values, not by parsing English text. English output is unchanged (the old snapshots still pass).
- **Running as a service:** a systemd unit, `/health` (uptime, network, LLM on or off), a clean shutdown on SIGTERM, and an RPC check at start-up. The bot is now installed on the server and runs on **mainnet** (it ran on testnet before).
- **Eval set for `/ask`:** 50 questions and an automatic scorer (`npm run eval`). See section 2, case 7 for what has and has not been measured.
- **Tests:** 383 → 435, all offline.
- **Checked in Telegram:** I sent the commands to the live mainnet bot and kept 13 screenshots (section 2, cases 9 to 13).

---

## Evidence

| # | Case | Result | Evidence |
|---|---|---|---|
| 1 | Token decimals in `/tx` and `/addr` on mainnet | `xudt 0x178f…a0a7 -725000000000000` became `-7250000 Seal` (8 decimals). `/addr` of a JoyId address lists `RGB++ Protocol … 1000 RGB++`; a token-rich address is in case 11 | [`evidence/2_token_decimals.txt`](./evidence/2_token_decimals.txt) |
| 2 | `/tx` on the 20 fixture transactions in en, zh and vi (60 runs) against mainnet | 60/60: explorer link present, correct fee, under Telegram's 4096-character limit | run by hand, summarized here |
| 3 | `/sign` on real mainnet transactions (witnesses removed) in zh, vi and en | Warnings translated; the script type is translated too after a fix (section 3) | [`evidence/4_sign_zh_vi_en.txt`](./evidence/4_sign_zh_vi_en.txt) |
| 4 | `npm run check:sign` after the changes | **20/20 match** between the unsigned and committed paths | [`evidence/3_check_sign_20of20.txt`](./evidence/3_check_sign_20of20.txt) |
| 5 | Quota and language across two separate processes (limit set to 2) | Process 2 still knew the 2 commands were used and the language was zh | [`evidence/5_persistence_two_processes.txt`](./evidence/5_persistence_two_processes.txt) |
| 6 | systemd: install, `restart`, `kill -9` | Restart logged "SIGTERM received, stopping" and came back; after `kill -9` systemd restarted it in 5 seconds | [`evidence/6_service_journal.txt`](./evidence/6_service_journal.txt) |
| 7 | Eval set and scorer | 50 questions (13 tx, 8 address, 5 DAO, 6 token, 6 investment advice, 4 prompt injection, 3 bad input, 5 out of scope; 31 en, 10 vi, 9 zh). The 32 questions that need a tool were run directly against mainnet: all worked and each expected fee appears in the tool output | [`evidence/7_eval_set_breakdown.txt`](./evidence/7_eval_set_breakdown.txt), [`eval/questions.json`](../../ckb-analyst/eval/questions.json) |
| 8 | All tests | 435 passed (after the two display fixes in section 3) | [`evidence/1_tests_435.txt`](./evidence/1_tests_435.txt) |
| 9 | Live bot in Telegram: `/help`, `/lang` without a value, `/health` | Help lists `/lang`; `/lang` alone replies `Usage: /lang en \| zh \| vi`; `/health` replies `OK · uptime 1m · network mainnet · LLM off` | [`telegram/01_help.png`](./telegram/01_help.png), [`02_lang.png`](./telegram/02_lang.png), [`03_health.png`](./telegram/03_health.png) |
| 10 | Live `/tx` of the RGB++ token transfer in English | Fee 0.0007293 CKB, 253 parties; token amounts read `-7250000 Seal` and `+1000 Seal` instead of raw units; Bitcoin outputs linked | [`telegram/04_en_tx.png`](./telegram/04_en_tx.png) |
| 11 | Live `/addr` and `/token` in English | `/addr`: 24151.21245522 CKB, 200 CKB in Nervos DAO, tokens `Seal`, `CANDY`, `LP-UTXOSwap V1` with decimals. `/token`: `LP-UTXOSwap V1`, Decimals 8, 3 holders, top holder 99.99% | [`telegram/05_en_addr.png`](./telegram/05_en_addr.png), [`06_en_token.png`](./telegram/06_en_token.png) |
| 12 | Same three commands in Chinese and in Vietnamese | Same numbers as English, labels translated, status shown as `已确认` / `đã xác nhận`. In the token reply the Facts lines stay English (known limit) | zh: [`07`](./telegram/07_zh_tx.png), [`08`](./telegram/08_zh_addr.png), [`09`](./telegram/09_zh_token.png); vi: [`10`](./telegram/10_vi_tx.png), [`11`](./telegram/11_vi_addr.png), [`12`](./telegram/12_vi_token.png) |
| 13 | Live `/sign` of a real mainnet DAO claim, witnesses removed, with the signer address | 🔴, type Nervos DAO claim (phase 2), `You receive 2447.18365693 CKB net`, NervosDao (official, code is fixed), the finding "Input … is already spent" and "including 2447.18366629 CKB compensation" | [`telegram/13_sign.png`](./telegram/13_sign.png) |

**What case 7 does not show.** I have no LLM configured, so **no accuracy figure exists yet**. I did run the whole eval pipeline once against a small fake server that pretends to be an LLM. That only proved the runner, the scorer and the report file work over HTTP. Its pass rate means nothing and I deleted the result.

**What cases 9 to 13 show and do not show.** They show the real bot answering through Telegram on mainnet, with the numbers I had already checked from code (for example the fee and the 2447.18366629 CKB compensation match the earlier runs). They do not cover `/dao`, `/sign` in Chinese or Vietnamese, the daily limit, or a restart while a user keeps a chosen language; those were checked only by calling the same handler or process directly (cases 3, 5 and 6). The `/sign` screenshot is red because the transaction is an old one with its signatures removed, so its input is already spent. That red result is the correct one.

---

## Key Learnings

- **My first check was empty.** I tested "token decimals" on a transaction whose net token change was zero, so no token was shown and the check passed for the wrong reason. I only saw it when I printed the output. I then searched mainnet for real tokens and re-ran. A passing check needs a case that could fail.
- **A real run found a bug that tests did not.** In the Chinese `/sign` output, the script type still read "lock" in English. Unit tests covered the template, but not this parameter. I fixed it and added a test.
- **One of my own eval questions was wrong.** For "is this token a good investment?" I had required that no tool be called. Calling the token tool for facts is correct. A bad expectation would have counted a good answer as a failure.
- **Look at the real reply, not only the tests.** Reading actual output found two display problems that all tests allowed: a token with an empty name showed a blank before its symbol, and the transaction status stayed `committed` in Chinese and Vietnamese. Both are fixed and covered by tests.
- **Keep the translation next to the data.** Findings keep their English text, and carry a code plus values. That let me translate them without parsing sentences, and English output stayed exactly as before.
- **Persist small state simply.** One bot process and a few hundred bytes do not need a database. An atomic file write (write temp file, rename) is enough, and a corrupt file must never stop the bot.


## Reflection

Most of this week's work was easy to write and easy to get wrong without noticing. The two mistakes I caught (an empty token check and a wrong eval expectation) were both mine, and both looked like success at first. I will keep one rule: before saying something works, show the output that would have looked different if it did not.

# CKB Weekly Report - Week 6

**Reporting period:** 22 October - 28 October 2026  
**Publication date:** 28 October 2026  
**Participant:** Tran Huy Hiep  

---

## 1. Overview

This week I moved from learning to building. I shipped **CKB Chain Analyst**, a Telegram bot that explains CKB mainnet data in plain language. You send it a tx hash, an address or a question.

- **Problem:** Raw CKB transactions are hard to read. There are many lock types (JoyID, Omnilock, RGB++), many token standards (sUDT, xUDT, Spore) and the Nervos DAO. General-purpose LLMs often explain CKB in EVM terms.
- **Approach:** Deterministic code computes every number. The LLM only picks tools and explains their results, and it never gives investment advice.

**Accomplishments:**
- A script registry with 29 entries, built from CCC, the Spore SDK and rgbpp-sdk.
- Analytics for decoding transactions, profiling addresses, reporting DAO status (compensation and unlock time) and showing sUDT/xUDT token facts.
- An LLM agent with 4 tools and a grounding check that flags any number not present in the tool output.
- A Telegram bot with `/tx`, `/addr`, `/dao` and `/token`, plus a daily free quota.
- 356 passing tests that use 20 real mainnet transactions as fixtures.

**Stack:** TypeScript, Node.js 22, `@ckb-ccc/core`, `@ckb-ccc/spore`, `grammy`, and the public RPC at `mainnet.ckb.dev`. The design is RPC-first, so it needs no node and no database.

---

## 2. Architecture

```
CKB RPC/indexer → chain/ → registry/ → analytics/ ─┬→ format/ → Telegram commands
                                                   └→ agent/ (LLM + grounding) → free-text questions
```

---

## 3. Evidence

| Command | Cases tested | Evidence |
|---|---|---|
| `/tx` | CKB transfer, JoyID transfer, xUDT transfer, sUDT mint, DAO withdraw and claim, Spore mint, Spore cluster, RGB++ airdrop (253 parties), RGB++ Spore, unregistered scripts | [`tx/`](./tx/) |
| `/addr` | JoyID wallet: 868.99894316 CKB and 1 xUDT | [`addr/`](./addr/) |
| `/dao` | Large depositor: 3,093,035.22 CKB in 29 deposits. Also an Omnilock wallet | [`dao/`](./dao/) |
| `/token` | Seal (3,722 holders) and RGB++ (1,272 holders) | [`token/`](./token/) |

**Verification:**
- Fees match the explorer, including DAO claims.
- Balances are conserved in all 20 fixture transactions.
- The DAO math matches CCC in 200 of 200 cases.

---

## 4. Key Learnings

- **Cell Model:** You have to rebuild "who sent what" by matching inputs with outputs and grouping them by lock and by type.
- **DAO fee:** The formula `inputs − outputs` gives a negative result for a claim, so the compensation has to be added. The compensation comes from the header_deps.
- **xUDT owner mode (RFC 0052):** The type args encode who can mint the token.
- **RGB++:** The RgbppLock args link a CKB cell to a Bitcoin UTXO.
- **Safety:** On-chain text is untrusted data that can carry prompt injection. Numbers must come from code, never from the LLM.

---

## 5. Challenges

| Challenge | Solution |
|---|---|
| The VM lacks the disk space for a full node and database | Used an RPC-first design on the public indexer |
| DAO claim transactions showed a negative fee | Added the DAO compensation to the fee formula |
| Block scanning timed out | Fetched blocks in parallel batches of 25 |
| Supply scans are slow for large tokens | Capped the scan at 5,000 cells and labeled the figures as partial |

---

## 6. Next Steps

1. Display token amounts with their decimals and symbols.
2. Run the bot as a service (pm2/systemd).
3. Connect a real LLM and evaluate it on a set of about 50 questions.
4. Invite 5 to 10 community users to try the bot.
5. In phase 2, charge per use in CKB, add streaming payments over **Fiber**, and offer an API/MCP for AI agents.

---

## 7. Reflection

This was the first week of building a real product, and it used every earlier week:
- the Cell Model and UDT drive the decoder;
- the DAO economics drive the compensation calculator;
- Fiber and Perun define where the product goes next.

The main lesson is to verify against the chain itself and keep every number in deterministic code.

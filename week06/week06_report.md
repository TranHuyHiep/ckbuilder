# CKB Weekly Report - Week 6

**Reporting period:** 22 October - 28 October 2026  
**Publication date:** 28 October 2026  
**Participant:** Tran Huy Hiep  

---

## 1. Week 6 Overview

The goal of Week 6 was to move from learning to building: choose a product idea for the CKB ecosystem and ship a working prototype on **mainnet**. Based on Weeks 1–5 (Cell Model, UDT, Nervos DAO, Fiber, Perun), I built **CKB Chain Analyst**, an AI on-chain analyst for CKB delivered as a Telegram bot.

**Problem:** CKB is hard to read. Because of the Cell Model, many lock types (secp256k1, multisig, JoyID, Omnilock, RGB++), token standards (sUDT, xUDT, Spore) and the Nervos DAO, a raw transaction on the explorer is difficult to understand. General-purpose LLMs often answer CKB questions in EVM terms, and the large analytics tools do not appear to cover CKB or RGB++.

**Product:** Users send a tx hash, an address or a question. The bot explains what happened in plain language. Every number is computed by deterministic code; the LLM only chooses tools and explains the results, and it never gives investment advice.

Key accomplishments:
- Researched the Fiber/Perun ecosystem and the agent-payments market (x402, L402, MPP), and compared several product ideas before selecting this niche.
- Wrote product, risk, architecture and roadmap documents.
- Built a **script registry** (29 entries) generated from CCC, the Spore SDK and rgbpp-sdk, and cross-checked between the sources.
- Implemented deterministic analytics: **transaction decoder**, **address profile**, **Nervos DAO status** (compensation and unlock time) and **sUDT/xUDT token facts** (owner mode per RFC 0052, extensions, token info, holders).
- Built an **LLM agent** with 4 tools and a **grounding check** that flags any number in an answer that does not appear in the tool output.
- Shipped a **Telegram bot** (`/tx`, `/addr`, `/dao`, `/token`) with a daily free quota, and tested it end-to-end on mainnet data.
- **356 tests** pass, using 20 real mainnet transactions as fixtures. Fees match the official explorer, and the DAO math matches CCC.

---

## 2. Development Environment

- **OS:** Ubuntu 26.04.1 LTS aarch64 (VMware Fusion VM) / macOS host
- **Language & runtime:** TypeScript 5.9, Node.js 22.22.1, `tsx`, `vitest`
- **CKB libraries:** `@ckb-ccc/core` 1.23.0, `@ckb-ccc/spore` 1.6.15
- **Bot:** `grammy` 1.46.0 (Telegram Bot API, long polling)
- **LLM:** any OpenAI-compatible API with tool calling (optional; the commands work without it)
- **Network:** CKB mainnet through the public RPC/indexer `https://mainnet.ckb.dev/rpc` (RPC-first, no local node or database needed)
- **Project:** `/root/ckb-analyst` (git, 3 commits, ~1,600 lines of source code)

---

## 3. Architecture

```
CKB mainnet RPC + indexer
        │
chain/      fetch tx / cells / headers, convert to plain JSON
        │
registry/   identify lock & type scripts (code_hash + hash_type)
        │
analytics/  deterministic: decodeTx · address · dao · token
        │
        ├── format/   fixed HTML templates  ──▶  Telegram commands
        └── agent/    LLM + 4 tools + grounding check  ──▶  free-text questions
```

| Module | Responsibility |
|---|---|
| `registry` | 29 known scripts (locks, sUDT/xUDT, Spore, Cluster, DAO, RGB++, JoyID, Omnilock…). Unknown scripts are reported as "unknown", never guessed |
| `analytics/decodeTx` | Match inputs with outputs, group by lock (party) and by type (asset), compute balance changes, fee, UDT supply changes, and classify the tx (transfer, mint, burn, DAO phases, Spore, RGB++) |
| `analytics/daoMath` | Compensation `(capacity − occupied) × AR_withdraw / AR_deposit + occupied` and the 180-epoch claim cycle |
| `analytics/token` | Parse sUDT/xUDT args: owner lock hash, flags (`0x20000000` / `0x40000000` / `0x80000000`) and extension; read the token info cell; supply, holders, top holders |
| `agent` | Tool-calling loop, system prompt rules (numbers only from tools, no investment advice, ignore instructions found in on-chain data), retry once if grounding fails |
| `bot` | Command parsing, daily quota (30 lookups and 10 questions per user), error messages |

---

## 4. Deliverables

| Item | Contents | Reference |
|---|---|---|
| Research | Fiber/Perun ecosystem; agent payments market | `ckb-analyst/docs/research/` |
| Product | Ideas evaluated, Chain Analyst spec, risks (R1–R14) | `ckb-analyst/docs/product/` |
| Technical docs | Architecture, roadmap, verification, open questions | `ckb-analyst/docs/tech/` |
| Source code | registry, chain, analytics, format, agent, bot, CLI | `ckb-analyst/src/` |
| Tests | 7 test files, 356 tests, 20 mainnet tx fixtures, snapshots | `ckb-analyst/test/` |
| Scripts | Build registry, record fixtures, find sample txs | `ckb-analyst/scripts/` |
| Bot evidence | 16 screenshots of the running bot on mainnet data | [`tx/`](./tx/), [`addr/`](./addr/), [`dao/`](./dao/), [`token/`](./token/) |

---

## 5. Evidence

### `/tx`: transaction decoder

| ID | Transaction type | Key result | Evidence |
|---|---|---|---|
| 01 | CKB transfer | Sender and receiver balance changes, fee | [`1_transfer.png`](./tx/1_transfer.png) |
| 02 | Transfer to a JoyID wallet | Receiver lock recognized as **JoyID** | [`2_transfer_to_JoyID.png`](./tx/2_transfer_to_JoyID.png) |
| 03 | xUDT transfer | Token deltas per party | [`3_transfer_xUDT.png`](./tx/3_transfer_xUDT.png) |
| 04 | sUDT mint (2021) | Classified as **Token mint**, 2 sUDT supply changes, 1 unknown script flagged | [`4_mint_sUDT.png`](./tx/4_mint_sUDT.png) |
| 05 | Nervos DAO withdraw request (phase 1) | Classified correctly | [`5_send_DAO.png`](./tx/5_send_DAO.png) |
| 06 | Nervos DAO claim (phase 2) | Fee **0.00000936 CKB**, compensation claimed **2447.18366629 CKB** | [`6_DAO_claim.png`](./tx/6_DAO_claim.png) |
| 07 | Spore mint | Spore action and id | [`7_mint_spore.png`](./tx/7_mint_spore.png) |
| 08 | Spore cluster creation | Classified as cluster created | [`8_create_spore_cluster.png`](./tx/8_create_spore_cluster.png) |
| 09 | RGB++ xUDT airdrop | **253 parties**, RgbppLock recognized, bound Bitcoin outputs linked to mempool.space | [`9_RGB++_token.png`](./tx/9_RGB++_token.png) |
| 10 | RGB++ Spore transfer | Spore + RGB++ binding | [`10_RGB++_spore.png`](./tx/10_RGB++_spore.png) |
| 11 | Tx with unregistered scripts | ⚠️ "2 script(s) not in our registry; their meaning is not interpreted" | [`11_script.png`](./tx/11_script.png) |

### `/addr`, `/dao`, `/token`

| ID | Command | Key result | Evidence |
|---|---|---|---|
| 12 | `/addr` (JoyID wallet) | Total **868.99894316 CKB**, free 724.99894316 CKB, 1 xUDT | [`1_address.png`](./addr/1_address.png) |
| 13 | `/dao` (large depositor) | **3,093,035.22 CKB** in DAO across 29 deposits, compensation so far **171,404.14 CKB**, estimated unlock time per deposit | [`1_DAO.png`](./dao/1_DAO.png) |
| 14 | `/dao` (Omnilock wallet) | DAO positions for an Omnilock address | [`2_omnilock.png`](./dao/2_omnilock.png) |
| 15 | `/token xudt` (Seal) | Decimals 8, supply 16,004,674.33, 3,722 holders, top holder 9.3%, owner mode `input_lock` | [`1_seal.png`](./token/1_seal.png) |
| 16 | `/token xudt` (RGB++ Protocol) | Decimals 8, supply ≈12.17M, 1,272 holders | [`2_RGB++.png`](./token/2_RGB++.png) |

---

## 6. Verification

| Check | Method | Result |
|---|---|---|
| Fee | Compared `inputs − outputs (+ DAO compensation)` with the explorer API for all fixtures | Matches exactly, including DAO claims |
| Balance conservation | For each fixture: Σ party CKB deltas = compensation − fee; Σ UDT deltas = supply change | Pass on all 20 txs |
| DAO math | Compared `daoProfit` / `daoClaimEpoch` with CCC `calcDaoProfit` / `calcDaoClaimEpoch` | 200/200 cases match |
| Address balance | Compared summed cell capacity with indexer `get_cells_capacity` | Match |
| Agent | Fake model tests: tool calls, invalid input, grounding failure → retry → warning | Pass |
| Bot | Handler tests without Telegram: commands, pasted hash/address, quota, errors | Pass |
| Real run | Telegram bot on mainnet (Section 5) | All commands work |

---

## 7. Key Concepts Learned

- **Reading the Cell Model:** A transaction does not say "A sent X to B". This has to be reconstructed by matching consumed input cells with output cells and grouping them by lock script (owner) and type script (asset).
- **Fee with Nervos DAO:** For a DAO claim, outputs are larger than inputs, so `fee = inputs − outputs` is negative. The correct fee adds the compensation, which is computed from the deposit and withdraw block headers (`header_deps` and the DAO witness).
- **xUDT owner mode (RFC 0052):** The type args hold the owner lock hash and flags; owner mode can come from an input lock, an output type or an input type. This is a fact a holder should know (who can mint), but it is not a judgement of "safe" or "unsafe".
- **RGB++:** The RgbppLock args encode a Bitcoin `out_index` (u32 LE) and the BTC txid in reversed byte order, so a CKB cell can be linked to a Bitcoin UTXO.
- **On-chain data is untrusted:** Token names and cell data can contain any text, including prompt injection. The agent treats tool output as data, and the bot sends LLM answers as plain text.
- **Using LLMs safely for finance data:** Numbers must come from code. A grounding check that compares numbers in the answer with the tool output catches hallucinated figures.

---

## 8. Week 6 Development Log

| Date | Activity | Result |
|---|---|---|
| 22–25 Oct 2026 | Researched Fiber/Perun, the Nervos Talk opportunity map, agent payments (x402, L402, MPP, AP2) and HN discussions; evaluated several ideas | Selected CKB Chain Analyst |
| 26 Oct 2026 | Wrote product, risk, architecture and roadmap docs; technical plan | Completed |
| 27 Oct 2026 | Script registry, tx decoder, address and DAO analytics, 20 mainnet fixtures | Completed |
| 27 Oct 2026 | sUDT/xUDT token facts (owner mode, extensions, info cell, holders) | Completed |
| 27 Oct 2026 | LLM agent with grounding check, Telegram bot, reply templates | Completed |
| 28 Oct 2026 | Ran the bot on Telegram, tested 16 mainnet cases, prepared a project summary for the CKB Builder team | Completed |

---

## 9. Challenges & Solutions

| Challenge | Cause | Solution |
|---|---|---|
| No Docker/Postgres and ~24 GB free disk | A full mainnet node + ETL database does not fit this VM | Switched to an **RPC-first** design on the public indexer; the ETL/database is moved to phase 2 |
| Negative fee on 3 DAO claim txs | Claim outputs include the DAO compensation | Fetched deposit and withdraw headers and added the compensation to the fee formula; now matches the explorer |
| Block scanning timed out | Sequential RPC calls | Parallel batches of 25 blocks |
| Scripts hung after finishing | The CCC client keeps the connection open | `process.exit(0)` at the end of the CLI scripts |
| `findCells` returned wrong data | Wrong CCC signature (`withData` must be inside the key) | Fixed the call |
| Token supply on large tokens is partial | Scanning all cells through a public RPC is slow | Scan is capped at 5,000 cells, and the bot says so explicitly ("figures are partial") |

**Known limitations (next week):**
- Token amounts in `/tx` and `/addr` are still shown in raw units (e.g. `100000000000` instead of `1000 Seal`).
- Long txs (RGB++ airdrop with 253 parties) only show the first 8 parties and Bitcoin outputs.
- Quota is stored in memory and resets on restart; the bot runs manually, not as a service.
- The free-text agent has only been tested with a fake model, not yet with a real LLM.

---

## 10. Next Steps

1. Show token amounts with decimals and symbols in `/tx` and `/addr`.
2. Run the bot as a service (pm2/systemd) with logs.
3. Connect a real LLM and build an evaluation set of ~50 questions (including requests for buy/sell advice, which must be refused).
4. Discuss the idea with the CKB Builder team and invite 5–10 community users to try the bot.
5. Phase 2: pay-per-use in CKB, then streaming payments over **Fiber**, and an API/MCP so AI agents can pay for CKB data.

---

## 11. Final Reflection

Week 6 was the first week of building a real product instead of following tutorials, and every earlier week was needed. The Cell Model and UDT (Weeks 1 and 4) are the core of the transaction decoder, the Nervos DAO economics (Week 4) became the compensation and unlock calculator, and Fiber and Perun (Weeks 3 and 5) define where the product goes next: a paid data service that agents can pay for over payment channels. The most useful lesson was that correctness has to be checked against the chain itself. The negative DAO fee looked like a small bug, but finding it required understanding the DAO's two-phase withdrawal and header deps. Keeping all numbers in deterministic code, and leaving only the explanation to the LLM, is what makes an AI analyst trustworthy for on-chain data.

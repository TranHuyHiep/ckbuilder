# CKB Weekly Report - Week 7

**Reporting period:** 29 October - 4 November 2026  
**Publication date:** 4 November 2026  
**Participant:** Tran Huy Hiep  

---

## 1. Summary

A community member told me that wallets should explain, in plain language, what a user is about to sign: what goes in, what goes out, and whether a contract is standard or a fork. Last week's bot only explained transactions *after* they happened.

This week I added a **pre-sign check** (`/sign`). It reads an unsigned transaction and explains it before signing. It is read-only: no keys, nothing is sent. Every number and judgment comes from code, not from the LLM.

**What was added**
- **`/sign` command and `sign` CLI:** CKB and token changes per party, the fee, and the effect on the signer.
- **Contract check:** is each script official, and is its code **fixed** or **upgradable** (and by whom).
- **Risk findings:** spent inputs, unusual fee, burned or unrecognized output lock, upgradable unknown script, token burn, missing cell deps.
- **Testnet support** (`CKB_NETWORK=testnet`), plus a wallet generator and an unsigned-tx builder for safe testing.
- **Tests:** 356 → 374, all offline.

---

## 2. Evidence

All runs used CKB **testnet** with unsigned transactions. Nothing was signed or sent.

| # | Case | Result | Evidence |
|---|---|---|---|
| 1 | `/help` lists the new command | `/sign <tx hex or JSON> [your address]` | [`telegram/01_help_list_function.png`](./telegram/01_help_list_function.png) |
| 2 | Transfer of 100 CKB, all checks pass | 🟢 "No problems found". Signer pays 100.00000371 CKB, of which 0.00000371 is the fee. Lock is secp256k1, official, code fixed | [`telegram/03_sign_no_problems.png`](./telegram/03_sign_no_problems.png) |
| 3 | Transaction whose inputs were already spent | 🔴 three "already spent" findings, so the tx would be rejected | [`telegram/02_allert_not_valid_sign.png`](./telegram/02_allert_not_valid_sign.png) |
| 4 | 100 CKB sent to a burned lock (CLI) | `level: danger`, finding `burned_lock_output` | [`sign/2_burned_lock_cli.json`](./sign/2_burned_lock_cli.json) |
| 5 | Normal transfer (CLI) | `level: ok` | [`sign/1_transfer_ok_cli.json`](./sign/1_transfer_ok_cli.json) |

**Cross-check:** I took a mainnet Spore mint, removed its witnesses and ran `sign`. The fee (0.1 CKB) matched the existing `tx` command.

**Not verified yet:** Nervos DAO, token and Spore transactions before signing; the fee thresholds (warn > 0.01 CKB, danger > 1 CKB) are my own choice; the network-fetch code is tested by hand, not by automated tests. On testnet only core scripts (secp256k1) are recognized, because the registry holds mainnet code hashes.

---

## 3. Key Learnings

- **Feedback changed the product.** The need is at signing time, so the useful output is a library a wallet can call. The bot is now mainly a demo.
- **Upgradability depends on how a script is referenced.** Data hash means the code is fixed. Type hash means whoever controls the code cell's lock can replace the code. The tool names that party but does not judge whether they are trustworthy.
- **Not every script runs in a transaction.** Output-only locks are not executed, so the tool reports them as "not checked" instead of guessing.
- **I caught my own mistake.** My test builder produced a transaction with no cell deps (invalid on chain) and `sign` still said "ok". It now warns (`dep_missing`).
- **A correction.** I had told the community the bot could detect upgradable contracts. At that time the code did not. This week it does.

---

## 4. Challenges

| Challenge | Solution |
|---|---|
| Demo transactions from a busy address showed "already spent" (screenshot 02) | Built the passing demo (screenshot 03) from a funded address idle for about 26,000 blocks |
| Builder output had no cell deps | Builder now adds them; `sign` warns when they are missing |
| Test runner failed to start (npm optional-dependency bug) | Installed the native binding without changing `package.json` |

---

## 5. Open Questions and Next Steps

I asked whether something like this already exists on CKB (for example CoBuild) and whether CCC is the right home. The reply was an impression ("feels like a wallet plugin or CCC"), without links. I have not yet read CoBuild or the CCC signer code.

**Next**
1. Read CCC's signer flow and the CoBuild spec, then fix the draft proposal (`docs/product/04-ccc-integration.md`, not posted).
2. Check DAO and token transactions before signing; show token amounts with decimals.
3. Make the registry network-aware so tokens and Spore work on testnet.
4. Build the library to plain JavaScript so a wallet can install it.

---

## 6. Reflection

The most valuable input this week was one sentence from a user, which showed my product was aimed at the wrong moment. Testing on real chain data also exposed a bug that unit tests had missed, so I keep the rule from week 6: verify against the chain and report only what I have seen.

# CKB Weekly Report - Week 8

**Reporting period:** 5 - 11 November 2026  
**Publication date:** 11 November 2026  
**Participant:** Tran Huy Hiep  

---

## 1. Summary

Last week I added `/sign`, a read-only check of an unsigned transaction. This week I made it more useful and tested it against real chain data. It is still read-only: no keys, nothing is sent, and every number and judgment comes from code, not from the LLM.

**What was added**
- **Token decimals:** token amounts show with decimals and symbol (read from the token's info cell) instead of raw units.
- **Nervos DAO and token handling:** deposit, withdraw request and claim (with compensation) are explained; token mint and burn are flagged.
- **Spore, RGB++ and DID:** shown in `/sign`, with the Bitcoin output an RGB++ asset is bound to. Burning a Spore or destroying a DID gives a warning.
- **Large-share warning:** 🟡 when a transaction moves 90% or more of the signer's balance.
- **Network-aware registry:** testnet and mainnet each use their own list of known scripts.
- **Consistency check:** `npm run check:sign` compares `sign` with the existing `tx` command on 20 mainnet transactions.
- **Tests:** 374 → 383, all offline.

---

## 2. Evidence

| # | Case | Result | Evidence |
|---|---|---|---|
| 1 | `sign` vs `tx` on 20 mainnet transactions (witnesses removed) | **20/20 match** on type, fee, balances, tokens, DAO, Spore, RGB++, unknown scripts | [`sign/3_check_sign_20of20.txt`](./sign/3_check_sign_20of20.txt) |
| 2 | Testnet transfer of 100 CKB from an address holding 3,888.74 CKB | 🟢 `level: ok`, no finding except "not signed yet" | [`telegram/1_sign_no_problem.png`](./telegram/1_sign_no_problem.png), [`sign/1_transfer_100ckb_ok.json`](./sign/1_transfer_100ckb_ok.json) |
| 3 | Same address sends 3,700 CKB | 🟡 `level: warn`, finding `large_share`: "moves 95.14% of this address's balance" | [`telegram/2_sign_review.png`](./telegram/2_sign_review.png), [`sign/2_large_share_warn.json`](./sign/2_large_share_warn.json) |
| 4 | Mainnet Spore mint, witnesses removed | Fee 0.1 CKB, same as `tx`; minted Spores listed | run by hand |
| 5 | Mainnet DAO claim, witnesses removed | Fee 0.00000936 CKB and compensation 2,447.18366629 CKB, same as `tx` | run by hand |
| 6 | Mainnet xUDT transfer, witnesses removed | Fee 0.00002022 CKB, same as `tx`; info cell found (8 decimals) | run by hand |

Cases 2 and 3 were sent to the Telegram bot on testnet (screenshots above) and gave the same two levels as the CLI.

**What case 1 proves and does not prove.** The signed and unsigned paths share the same decoder, so 20/20 shows the unsigned path does not change the result. It does not prove the decoder is right. That rests on the explorer cross-check from week 6.

---

## 3. Key Learnings

- **Verify before claiming.** My first `check:sign` run reported "transaction not found" because my `.env` pointed at testnet while the fixtures are mainnet. The fix was one environment variable, but the result looked like a bug in the tool.
- **Defaults are not calibrated.** The fee thresholds and the 90% share are my own choices. I report them as defaults, not as facts about risk.
- **Output-only locks are not executed.** The tool says "code not checked" for them instead of guessing.
- **Different asset types, different checks.** DAO needs header data for compensation, tokens need an info cell, RGB++ needs the lock args decoded. One generic path would hide these differences.

---

## 4. Challenges

| Challenge | Solution |
|---|---|
| Fixtures are mainnet, `.env` is testnet | Run `CKB_NETWORK=mainnet npm run check:sign`; documented in the README |
| 43 tests failed when `CKB_NETWORK=testnet` was exported in the shell | Run tests with it unset |
| Demo transactions turned "already spent" on a busy address | Used an address idle for about 26,000 blocks |
| A test stub lacked the new Spore and RGB++ fields | Updated the stub |

---

## 5. Limits and Next Steps

**Not verified yet**
- No unsigned DAO or token transaction built by a real wallet; checks used mainnet transactions with witnesses removed, which can differ from a wallet's output.
- DID is covered by unit tests only; I did not find a real DID transaction to run.
- The free-text `/ask` path with an LLM has not been run (no LLM configured).
- The testnet registry has no RGB++ entries.
- I have not yet read the CCC signer code or the CoBuild spec, so the draft CCC proposal is still unposted.

**Next**
1. Test with unsigned transactions from a real wallet, including DAO and token cases.
2. Find a real DID transaction and run it through `sign`.
3. Read CCC's signer flow and CoBuild, then correct the draft proposal.
4. Run the bot as a persistent service and decide on a mainnet switch.

---

## 6. Reflection

This week was about trust in my own output. A 20/20 match is reassuring but only shows consistency, so I wrote down exactly what it does not prove. I will keep that habit: say what was checked, what was not, and which numbers are my own defaults.

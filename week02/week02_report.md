# CKB Weekly Report - Week 2

**Reporting period:** 24-28 September 2026  
**Publication date:** 28 September 2026  
**Participant:** Tran Huy Hiep

## 1. Week 2 Overview

The goal of Week 2 was to move from environment setup to application development and programmatic transaction handling on CKB.

The completed exercises and modules:
- **Module 0:** CKB Academy - Getting Started with NFTs
- **Module 1:** CCC (Common Chain Connector) & CCC Playground
- **Module 2:** Full-Stack Mini dApp Development (`my-ccc-app`)
- **Module 3:** L1 Developer Training Course & Labs:
  - CKB Node & CLI Management
  - Lab: Calculating Capacity Requirements (Lumos)
  - Lab: Automated Cell Collection (Lumos)
  - On-Chain Data Storage (CLI & Lumos SDK)

## 2. Development Environment

- **OS:** macOS (Apple Silicon arm64)
- **Node.js:** v25.2.1 | **pnpm:** 12.5.1
- **CKB-CLI:** v1.12.0 (Commit `278c7be`)
- **Frameworks:** CCC (`@ckb-ccc/ccc`, `@ckb-ccc/connector-react`), Lumos SDK, Next.js 16
- **Networks:** Local Devnet (`http://127.0.0.1:8114`) & Nervos Meepo Testnet

## 3. Evidence

| ID | Module / Task | Key Result | Evidence |
|---|---|---|---|
| 00 | CKB Academy: NFTs | Completed 6 modules (Spore vs CoTA) | [`0_ckb_academy/`](./0_ckb_academy/3_Getting_Started_With_NFTs.png) |
| 01 | CCC Playground | Transfer tx & Spore DOB minting | [`1_ccc/`](./1_ccc/) |
| 02 | Mini dApp (Next.js) | Testnet transfer confirmed (Tx: `0xdb7610...ea0f`) | [`2_mini_dapp/`](./2_mini_dapp/evidence/) |
| 03 | CKB-CLI & Node | Account setup & 100 CKB transfer | [`3_L1_Course: CLI`](./3_L1_Developer_Training_Course/2_ckb_cli.png) |
| 04 | Lab: Capacity Calculation | Exact fee calculation (Tx: `0x2fe6e3...4f8eb`) | [`Lab-Capacity`](./3_L1_Developer_Training_Course/Lab-Calculating-Capacity-Requirements-Exercise/) |
| 05 | Lab: Automated Cell Collection | Aggregated 3 inputs = 183 CKB (Tx: `0x24db28...a6df`) | [`Lab-Collection`](./3_L1_Developer_Training_Course/Lab-Implement-Automated-Cell-Collection-Exercise/) |
| 06 | On-Chain Data Storage | Stored 13B text in 74 CKB cell (Blake2b verified) | [`Data Storage`](./3_L1_Developer_Training_Course/10_storing_data_in_a_cell.png) |

All evidence files are stored in their respective directories.

## 4. Key Concepts Learned

- **61 CKB Minimum Floor:** A cell requires at least 61 CKB (8 bytes capacity + 53 bytes default lock script) to exist on-chain. Storing data increases required capacity by 1 CKB per byte.
- **Change Cell Rule:** When transferring CKB, any change output must also be $\ge 61$ CKB. If remaining funds are less than 61 CKB, the transaction must collect more input cells to form a valid change cell.
- **Spore vs. CoTA:** Spore stores DOB data directly in cell data backed by CKB capacity (burnable to reclaim CKB). CoTA aggregates thousands of NFTs into a single cell using Sparse Merkle Trees (SMT).
- **CCC vs. Lumos:** CCC provides high-level automated transaction building (`completeInputsByCapacity`, `completeFeeBy`) and unified wallet support (JoyID, UTXO Global), while Lumos handles low-level manual transaction skeletons.

## 5. Module 0 - CKB Academy: Getting Started with NFTs

### Objective
Understand CKB NFT standards (Spore Protocol and CoTA) and their architectural trade-offs.

### Procedure
- Studied Spore Protocol architecture, DOB lifecycle, and on-chain capacity backing.
- Explored CoTA aggregation using Sparse Merkle Trees for high-volume minting.
- Completed all 6 course chapters on `academy.ckb.dev`.

### Result
Completed the full curriculum with all badges verified.

### What I Learned
Spore is ideal for high-value on-chain digital objects with intrinsic capacity value, while CoTA is best for mass low-cost issuance.

### Evidence
See the [NFT Course Evidence](./0_ckb_academy/3_Getting_Started_With_NFTs.png).

## 6. Module 1 - CCC Playground

### Objective
Learn transaction construction, Spore minting, and chain querying with the CCC SDK.

### Procedure
- Built a 100 CKB transfer in CCC Playground (`live.ckbccc.com`).
- Used `completeInputsByCapacity` and `completeFeeBy` to automatically balance inputs and fees.
- Minted a Spore NFT with plain text data (`Hello, Spore!`).
- Queried tip block (`22559244`) and account balance (`11,444 CKB`).

### Result
- **Transfer Tx:** `0xffa0ee05abd57038751a5648ee8c2ac295408c8a9a61a617ffee62d530b0440c` (Inputs: 183 CKB, Change: 82.99 CKB).
- **Spore Tx:** `0x264810027e7537f5f2f65de8a75b9c5201a265ea2aedead2b302a716dfeda839` (173 CKB Spore cell).

### What I Learned
CCC significantly simplifies transaction assembly compared to manual skeleton building.

### Evidence
See the [CCC Playground Evidence folder](./1_ccc/).

## 7. Module 2 - Full-Stack Mini dApp (`my-ccc-app`)

### Objective
Build and deploy a Next.js dApp using CCC, connect a wallet, and send a transaction on testnet.

### Procedure
- Built a Next.js 16 app with `@ckb-ccc/connector-react` and Tailwind CSS.
- Connected JoyID / UTXO Global wallet on Nervos Meepo Testnet.
- Added transfer UI with input validation (minimum 61 CKB).
- Sent 101 CKB to a recipient address.

### Result
- Transaction broadcasted and confirmed on Nervos Meepo Testnet.
- **Tx Hash:** `0xdb7610d8d363c501b53e5ffa96fd99ca9a9e775dc78f19cc97dc0572fe40ea0f`
- **Block:** 22,564,734 | **Cycles:** 1,619,064 | **Fee:** 0.00015416 CKB.

### What I Learned
Learned how to integrate CCC connector hooks in React and handle testnet wallet signatures.

### Evidence
See the [Mini dApp Evidence folder](./2_mini_dapp/evidence/).

## 8. Module 3.1 - CKB-CLI & Node Operations

### Objective
Interact directly with the CKB node using `ckb-cli`.

### Procedure
- Verified genesis hash and system cell dep groups via `./ckb list-hashes`.
- Created a new testnet keypair with `ckb-cli account new`.
- Transferred 100 CKB via `ckb-cli wallet transfer`.
- Inspected transaction and live cell via RPC commands.

### Result
- **Tx Hash:** `0x48b62d1424afb4378e45ac05dc5827d5b73c6302185a132e57184716c6980a48`.
- RPC confirmed live cell with capacity `0x2540be400` (100 CKB).

### What I Learned
Learned direct CLI wallet operations and RPC debugging.

### Evidence
See [CLI Screenshots](./3_L1_Developer_Training_Course/2_ckb_cli.png).

## 9. Module 3.2 - Lab: Calculating Capacity Requirements

### Objective
Build a Lumos transaction with exact capacity, change calculation, and fixed fee.

### Procedure
- Constructed transaction skeleton in `index.js`.
- Consumed 1 input cell (998,999.9999 CKB).
- Created a 1,000 CKB target output and calculated change minus 10,000 Shannons fee.
- Validated and signed transaction.

### Result
- **Tx Hash:** `0x2fe6e3a92f6bb5bfecb2b93eb0f527770738e77efbf0121f6c27540c8443f8eb`.
- Outputs: 1,000 CKB + 997,999.9998 CKB change. Fee: exactly 0.0001 CKB. Lab passed.

### What I Learned
Understood exact Shannon math and zero-loss capacity conservation rules.

### Evidence
See [Lab 1 Evidence](./3_L1_Developer_Training_Course/8_lab_Calculating_Capacity_Requirements.png) and source files in [`Lab-Calculating-Capacity-Requirements-Exercise/`](./3_L1_Developer_Training_Course/Lab-Calculating-Capacity-Requirements-Exercise/).

## 10. Module 3.3 - Lab: Automated Cell Collection

### Objective
Implement automated coin selection to gather multiple cells when single cells are insufficient.

### Procedure
- Used `CellCollector` and `collectCapacity` to collect live cells matching sender lock.
- Accounted for target output (100 CKB) + minimum change cell (61 CKB) + fee.
- Aggregated 3 inputs of 61 CKB (183 CKB total).
- Generated 82.999 CKB change cell and broadcasted tx.

### Result
- **Tx Hash:** `0x24db285386f5f14b6d0da6f343be637feb611fd36998b22a2df36cafd1aea6df`.
- Lab validation passed successfully.

### What I Learned
Learned how UTXO aggregation works and why change output minimums must be factored into coin selection.

### Evidence
See [Lab 2 Evidence](./3_L1_Developer_Training_Course/9_Lab_Implement_Automated_Cell_Collection.png) and source files in [`Lab-Implement-Automated-Cell-Collection-Exercise/`](./3_L1_Developer_Training_Course/Lab-Implement-Automated-Cell-Collection-Exercise/).

## 11. Module 3.4 - On-Chain Data Storage

### Objective
Store raw text in cell data via CKB-CLI and Lumos SDK, and verify cryptographic hash.

### Procedure
- Stored `Hello Nervos!` (13 bytes) using `ckb-cli wallet transfer --to-data-path` (capacity: 74 CKB).
- Verified cell data and on-chain Blake2b hash matches `ckb-cli util blake2b`.
- Stored the same data programmatically using Lumos SDK (`Storing-Data-in-a-Cell-Example`).

### Result
- **CLI Tx:** `0xd30ed9e686e9058fa9b66644e04a8edca9191eb3a189d0fc82981e47563d1f98`.
- **Lumos Tx:** `0x7fb576c53e3f67e64fcd3082d2f6111837c69b229cdb23812f4ef3f8721e13fb`.
- Data hex `0x48656c6c6f204e6572766f7321` confirmed live on-chain.

### What I Learned
Learned how cell `data` field stores application state and how capacity scales with byte length.

### Evidence
See [Data Storage Evidence](./3_L1_Developer_Training_Course/10_storing_data_in_a_cell.png).

## 12. Development Log

| Date | Activity | Result | Evidence |
|---|---|---|---|
| 24 Sep 2026 | CKB Academy NFT Course | Completed 6 modules | [NFT Course](./0_ckb_academy/3_Getting_Started_With_NFTs.png) |
| 25 Sep 2026 | CKB-CLI & Node Operations | Account setup & transfer | [CLI](./3_L1_Developer_Training_Course/2_ckb_cli.png) |
| 26 Sep 2026 | Lumos Capacity & Cell Collection Labs | Both labs validated | [Labs](./3_L1_Developer_Training_Course/8_lab_Calculating_Capacity_Requirements.png) |
| 27 Sep 2026 | On-Chain Data Storage (CLI & Lumos) | Blake2b verified | [Data](./3_L1_Developer_Training_Course/10_storing_data_in_a_cell.png) |
| 28 Sep 2026 | CCC Playground & Mini dApp | Testnet transfer confirmed | [dApp](./2_mini_dapp/evidence/4_send_successfully.png) |

## 13. Challenges & Solutions

- **61 CKB Minimum Floor:** Creating cells with $<61$ CKB fails consensus. Solved by validating all outputs and change cells are $\ge 61$ CKB.
- **Change Cell Collection Deadlock:** Sending 100 CKB with 61 CKB inputs requires at least 3 cells ($183$ CKB), because 2 cells ($122$ CKB) leaves only 22 CKB change, violating the 61 CKB floor. Solved by including min change in `collectCapacity`.
- **Lumos Multi-input Witnesses:** Consuming multiple inputs requires correct witness table formatting. Solved using `addDefaultWitnessPlaceholders` and `prepareSigningEntries`.

## 14. Final Reflection

Week 2 successfully connected low-level transaction mechanics with modern dApp development. I now understand cell capacity calculation, coin selection algorithms, Spore vs CoTA standards, and how to build full-stack CKB dApps using CCC.
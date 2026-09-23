 # CKB Weekly Report - Week 1

 **Reporting period:** 21-23 September 2026  
 **Publication date:** 23 September 2026  
 **Participant:** Tran Huy Hiep

 ## 1. Week 1 Overview

 The goal of Week 1 was to set up a local CKB development environment, become familiar with the CKB Cell Model, and complete the beginner exercises.

 The completed exercises were:

 - Exercise 00: Getting Started
 - Exercise 01: Transfer CKB
 - Exercise 02: Store Data on Cell
 - Exercise 03: Create Fungible Token (xUDT)
 - Exercise 04: Create DOB
 - Exercise 05: Simple Lock

 ## 2. Development Environment

 - OS: macOS arm64
 - Node.js: v25.2.1
 - npm: 11.6.2
 - pnpm: 12.5.1
 - OffCKB: 0.4.13
 - CKB Debugger: 1.1.1
 - Repository: `/Users/tranhuyhiep/Documents/CKB`
 - Devnet RPC: `http://127.0.0.1:8114`
 - Devnet RPC Proxy: `http://127.0.0.1:28114`

 ## 3. Evidence

 | ID | Exercise | Evidence |
 |---|---|---|
 | 00 | Getting Started | [`1_setup/`](./1_setup/) |
 | 01 | Transfer CKB | [`2_transfer_CKB/`](./2_transfer_CKB/) |
 | 02 | Store Data on Cell | [`3_store_Data_On_Cells/`](./3_store_Data_On_Cells/) |
 | 03 | Create Fungible Token | [`4_create_Fungible_Token/`](./4_create_Fungible_Token/) |
 | 04 | Create DOB | [`5_create_DOB/`](./5_create_DOB/) |
 | 05 | Simple Lock | [`6_build_simple_Lock/`](./6_build_simple_Lock/) |

 All evidence screenshots are stored in their respective exercise directories.

 ## 4. CKB Fundamentals

 ### Cell Model

 CKB uses a Cell Model similar to the UTXO model. Blockchain state is represented by cells. Live cells can be consumed by transactions and replaced by newly created cells.

 ### Capacity

 Capacity represents the amount of CKBytes held by a cell and also pays for the cell's storage. In CKB, `1 CKB = 100,000,000 shannons`, and approximately one CKB represents one byte of storage capacity.

 ### Transaction

 A transaction consumes existing input cells and creates new output cells, representing a state transition. The total capacity of outputs must not exceed the total capacity of inputs; the difference is used as the transaction fee.

 ### Input and Output

 Inputs reference existing live cells that are consumed. Outputs define the new cells created by the transaction.

 ### Lock Script

 A Lock Script defines the conditions required to unlock and consume a cell. It commonly represents ownership or authorization.

 ### Type Script

 A Type Script is optional and validates application-specific rules for a cell's data and state transitions. Lock scripts protect ownership, while type scripts enforce asset or application logic.

 ### Data

 Cell data stores arbitrary application-specific bytes. It can represent a message, token data, or Digital Object content.

 ### Cell Dependency

 A script's code is stored in a dependency cell. The script's `code_hash` and `hash_type` identify the code that should be loaded during transaction verification.

 ### RPC

 RPC (Remote Procedure Call) allows applications and development tools to query CKB state and submit transactions to a node.

 ## 5. Exercise 00 - Getting Started

 ### Objective

 Install the required tooling and start a local CKB Devnet.

 ### Procedure

 - Installed the OffCKB CLI globally with `npm install -g @offckb/cli`.
 - Started the local Devnet with `offckb node`.
 - Listed the pre-funded development accounts with `offckb accounts`.

 ### Result

 - The CKB Devnet became available at `http://127.0.0.1:8114`.
 - The Devnet RPC Proxy became available at `http://127.0.0.1:28114`.
 - The pre-funded development accounts were listed and available for later exercises.

 ### What I Learned

 I learned how to install OffCKB, launch a local CKB Devnet, and inspect the accounts provisioned for local testing.

 ### Evidence

 See the [Exercise 00 evidence folder](./1_setup/).

 ## 6. Exercise 01 - Transfer CKB

 ### Objective

 Complete a CKB transfer between development accounts and understand how transactions affect account balances.

 ### Procedure

 - Ran the Simple Transfer dApp locally.
 - Connected the dApp to the local Devnet.
 - Used the dApp to transfer CKB between Devnet addresses.

 ### Result

 The transfer transaction was submitted successfully to the local CKB Devnet and confirmed.

 ### What I Learned

 I learned how a frontend dApp connects to CKB, selects cells, submits a transaction, and displays the resulting balance and transaction state.

 ### Evidence

 See the [Exercise 01 evidence folder](./2_transfer_CKB/).

 ## 7. Exercise 02 - Store Data on Cell

 ### Objective

 Learn how application data can be written to and read from a CKB cell.

 ### Procedure

 - Ran the Store Data on Cell dApp locally and connected it to the Devnet.
 - Wrote the message `hello common knowledge base!` into a cell.
 - Read the message back from the cell.

 ### Result

 - The write transaction completed successfully.
 - The Read action returned: `Message: hello common knowledge base!`.

 ### What I Learned

 I learned that a cell's `data` field can persist arbitrary application data on-chain and that a dApp can read that data back from the cell.

 ### Evidence

 See the [Exercise 02 evidence folder](./3_store_Data_On_Cells/).

 ## 8. Exercise 03 - Create Fungible Token (xUDT)

 ### Objective

 Issue a custom fungible token on the local Devnet, query it, and transfer part of its balance.

 ### Procedure

 - Ran the xUDT Scripts dApp locally.
 - Issued a custom token with an amount of `50`.
 - Queried the token using its xUDT args.
 - Transferred `1` token unit to another Devnet address.

 ### Result

 - The custom token was issued successfully.
 - The token query showed a cell holding an amount of `50`.
 - The transfer of `1` token unit was submitted successfully.

 ### What I Learned

 I learned how xUDT represents fungible tokens using CKB cells, how type script args identify a token, and how token balances are associated with cells and lock scripts.

 ### Evidence

 See the [Exercise 03 evidence folder](./4_create_Fungible_Token/).

 ## 9. Exercise 04 - Create DOB

 ### Objective

 Create an on-chain Digital Object (DOB/Spore) by storing image data in a CKB cell and verify the stored content.

 ### Procedure

 - Ran the Create DOB dApp locally and connected it to the Devnet.
 - Uploaded a JPEG image named `download.jfif`.
 - Created the Spore/DOB cell.
 - Checked the Spore content from the dApp.

 ### Result

 - The DOB was created successfully.
 - The content check confirmed `contentType: image/jpeg`.
 - The uploaded image was rendered from the on-chain cell data.

 ### What I Learned

 I learned that a Digital Object's content can be stored directly in CKB cell data and that the required capacity depends on the size of the stored content. I also learned that the selected network must match the network where the account is funded.

 ### Evidence

 See the [Exercise 04 evidence folder](./5_create_DOB/).

 ## 10. Exercise 05 - Simple Lock

 ### Objective

 Build and deploy a custom hash-lock contract to the CKB Devnet, run its frontend, and unlock a cell with the correct preimage.

 ### Procedure

 - Installed the required JavaScript dependencies with pnpm.
 - Installed `protobuf` and `ckb-debugger` for contract compilation.
 - Built and deployed the JavaScript hash-lock contract with `pnpm run deploy -- --network devnet`.
 - Started the Next.js frontend locally at `http://localhost:3000`.
 - Generated a hash-lock address from the preimage `Hello World`.
 - Deposited `123 CKB` into the generated hash-lock address.
 - Revealed the preimage and transferred `99 CKB` from the hash-lock cell.

 ### Result

 - `hash-lock.bc` was built and deployed successfully.
 - The frontend deployment health indicator showed `DEVNET · READY`.
 - The deposit was confirmed with `123 CKB` of live capacity.
 - The reveal-and-transfer transaction was committed successfully.
 - Transfer transaction hash: `0x5f1a9ab191d7ea68b8a5f4968b6edf3b6d5772e7663e5c5d04b94471f9d03ec6`

 ### What I Learned

 I gained practical experience compiling and deploying a custom JavaScript script, interacting with it through a Next.js frontend, and unlocking a cell by revealing its correct preimage. I also learned that returning change to the same hash-lock address is suitable only for this educational example; production transactions should use a signature-protected change address.

 ### Evidence

 See the [Exercise 05 evidence folder](./6_build_simple_Lock/).

 ## 11. Week 1 Development Log

 | Date | Activity | Result | Evidence |
 |---|---|---|---|
 | 21 Sep 2026 | OffCKB installation and Devnet startup | Completed | [Exercise 00 evidence](./1_setup/) |
 | 21 Sep 2026 | CKB transfer dApp setup and execution | Completed | [Exercise 01 evidence](./2_transfer_CKB/) |
 | 22 Sep 2026 | Store Data on Cell: write and read | Completed | [Exercise 02 evidence](./3_store_Data_On_Cells/) |
 | 22 Sep 2026 | Issue, view, and transfer a custom xUDT token | Completed | [Exercise 03 evidence](./4_create_Fungible_Token/) |
 | 22 Sep 2026 | Create DOB and verify Spore content | Completed | [Exercise 04 evidence](./5_create_DOB/) |
 | 23 Sep 2026 | Simple Lock contract build, deploy, deposit, and unlock | Completed | [Exercise 05 evidence](./6_build_simple_Lock/) |

 ## 12. Challenges

 The main challenges during Week 1 were:

 - **pnpm build approval:** pnpm 12 blocked native postinstall scripts for `esbuild`, `secp256k1`, `sharp`, and `unrs-resolver`. The packages were explicitly approved in `pnpm-workspace.yaml`, after which both normal and frozen-lockfile installs succeeded.

 - **Missing CKB Debugger:** The Simple Lock build initially failed because `ckb-debugger` was not installed. Running `cargo install ckb-debugger` exposed a second missing prerequisite: `protoc`.

 - **Missing protoc:** The Rust build for `ckb-debugger` failed because the Protocol Buffers compiler was unavailable. Installing `protobuf` with Homebrew fixed the toolchain, and `ckb-debugger 1.1.1` was installed successfully.

 - **Devnet versus public network configuration:** The Create DOB workflow required the dApp network configuration to match the network where the development account was funded. Using the wrong network can appear as an insufficient-capacity error.

 - **CKB Cell Model:** Understanding the difference between account balances and CKB cells required some adjustment. The exercises helped connect the Cell Model with real transactions, stored data, tokens, and scripts.

 All challenges were resolved, and the Week 1 exercises were completed successfully.

 ## 13. Final Reflection

 Week 1 gave me a practical introduction to developing on Nervos CKB. I set up a local development environment, started a Devnet, worked with development accounts, and completed CKB workflows involving transfers, cell data, xUDT tokens, Digital Objects, and a custom hash-lock contract.

 The practical exercises helped me connect the Cell Model with real transaction workflows and understand how cells, capacity, lock scripts, type scripts, data, and transaction fees work together. I am now prepared to continue with CKB application and script development.

 ## 14. Week 2 Goals - Building Applications on CKB

 The main goal for Week 2 is to move from beginner exercises to hands-on application development.

 I will focus on:

 - Exploring the CCC App and CCC Playground.
 - Learning the core CCC APIs and transaction workflow.
 - Building simple JavaScript/TypeScript CKB application flows.
 - Understanding how frontend applications interact with cells and transactions.
 - Setting up the Rust environment for CKB Script development.
 - Exploring CKB Rust SDK examples, CKB-CLI, and CKB Debugger.
 - Building and testing a simple CKB Script.

 By the end of Week 2, I aim to understand the basic workflow of building a CKB application with CCC, build and test simple application flows, understand the architecture and execution model of CKB Scripts, and have a working Rust environment for further development.

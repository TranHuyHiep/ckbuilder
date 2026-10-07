# CKB Weekly Report - Week 3

**Reporting period:** 29 September - 07 October 2026  
**Publication date:** 07 October 2026  
**Participant:** Tran Huy Hiep  

---

## 1. Week 3 Overview

The goal of Week 3 was to research and deploy the **Fiber Network** (CKB Layer 2 Lightning/Payment Channel Network), following the [Basic Transfer Guide](https://www.fiber.world/docs/quick-start/basic-transfer).

Key accomplishments:
- Installed and configured **Fiber Network Node (FNN) v0.9.1** on Ubuntu aarch64 (VMware).
- Set up 2 independent nodes (`node1` and `node2`) with separate keypairs and non-conflicting ports.
- Connected peers via P2P and opened a payment channel funded with **500 CKB** on Meepo Testnet.
- Generated a **100 CKB** testnet invoice on `node2` and executed an instant off-chain transfer from `node1`.
- Verified off-chain state updates and balance reallocation (Node 1: 301 CKB, Node 2: 100 CKB).

---

## 2. Development Environment

- **OS:** Ubuntu Linux aarch64 (VMware) / macOS Host
- **Fiber Node:** FNN v0.9.1
- **CKB Tools:** CKB-CLI v1.12.0, Python 3.12 (balance parser)
- **Networks:** Nervos Meepo Testnet (L1) & Fiber Testnet `Fibt` (L2)
- **Endpoints:**
  - **Node 1:** RPC: `http://127.0.0.1:8227` | P2P: `8228`
  - **Node 2:** RPC: `http://127.0.0.1:8237` | P2P: `8238`

---

## 3. Evidence

| ID | Task / Operation | Key Result | Evidence |
|---|---|---|---|
| 00 | Install Fiber Node | Bootstrap script completed; binary placed in `~/.fiber` | [`0_install_fiber.png`](./evidence/0_install_fiber.png) |
| 01 | Start Node Daemon | Node initialized with `store_actor` & watchtower running | [`1_node_running.png`](./evidence/1_node_running.png) |
| 02 | P2P Peer Connection | Node 1 connected to Node 2 (`/ip4/127.0.0.1/tcp/8238`) | [`2_connect_peer.png`](./evidence/2_connect_peer.png) |
| 03 | Open Channel | Initiated 500 CKB funding; Temp ID: `0x6f2e6f...98f` | [`3_open_channel.png`](./evidence/3_open_channel.png) |
| 04 | Funding Commitment | Channel negotiation passed; pending on-chain tx | [`4_funding_tx_committed.png`](./evidence/4_funding_tx_committed.png) |
| 05 | Channel Finalization | On-chain tx committed; `state_name: ChannelReady` | [`5_channel_ready.png`](./evidence/5_channel_ready.png) |
| 06 | Generate Invoice | 100 CKB invoice created on Node 2 (`fibt1...`) | [`6_create_invoice.png`](./evidence/6_create_invoice.png) |
| 07 | Send Payment | Instant transfer completed; `status: Success` | [`7_send_payment_success.png`](./evidence/7_send_payment_success.png) |
| 08 | Verify Balance (Node 1) | Local balance updated: 401.00 CKB -> **301.00 CKB** | [`9_balance_after_payment-1.png`](./evidence/9_balance_after_payment-1.png) |
| 09 | Verify Balance (Node 2) | Local balance credited: 0.00 CKB -> **100.00 CKB** | [`9_balance_after_payment-2.png`](./evidence/9_balance_after_payment-2.png) |

---

## 4. Key Concepts Learned

- **Layer 2 Payment Channels:** Only channel opening and closing require L1 transactions. Intermediary payments happen off-chain instantly with zero L1 gas fees.
- **HTLC / TLC Contracts:** Payments are routed using cryptographic preimages and hashlocks with expiration timeouts, ensuring atomic and trustless transfers.
- **FundingLock & CommitmentLock:** CKB scripts that enforce channel safety on L1 and penalize any party attempting to broadcast outdated channel states.
- **Invoice Protocol:** Bech32m-encoded payment requests (`fibt...` on testnet) that embed payee pubkey, amount in Shannons, expiry, and payment hash.

---

## 5. Practical Operations Walkthrough

### Step 1: Install FNN & Setup Keypairs
- **Objective:** Install Fiber node binaries and provision dedicated CKB accounts.
- **Procedure:** Installed FNN v0.9.1 via the official installer, generated accounts with `ckb-cli account new`, and exported raw 64-hex private keys to `./ckb/key` (`chmod 600`).
- **Result:** Node binaries and keypairs ready for both `node1` and `node2`.
- **Evidence:** [`0_install_fiber.png`](./evidence/0_install_fiber.png)

### Step 2: Configure Ports & Run Nodes
- **Objective:** Run two isolated node instances on a single machine without conflicts.
- **Procedure:** Assigned `node1` to RPC `8227` / P2P `8228` and `node2` to RPC `8237` / P2P `8238`. Started both daemons with `RUST_LOG=info ./fnn -c config.yml -d .`.
- **Result:** Both nodes started their actor subsystems (`store_actor`, `watchtower`).
- **Evidence:** [`1_node_running.png`](./evidence/1_node_running.png)

### Step 3: Peer Discovery & P2P Handshake
- **Objective:** Connect `node1` to `node2` over the local P2P network.
- **Procedure:** Retrieved Node 2's pubkey (`02d982a839...`) and ran `./fnn-cli peer connect_peer` from Node 1.
- **Result:** P2P connection established successfully.
- **Evidence:** [`2_connect_peer.png`](./evidence/2_connect_peer.png)

### Step 4: Open & Fund Payment Channel
- **Objective:** Lock 500 CKB on CKB Meepo Testnet to establish a payment channel.
- **Procedure:** Executed `./fnn-cli channel open_channel --funding-amount 50000000000 --public true` from Node 1. Monitored until on-chain confirmation.
- **Result:** Channel `0x9f89856444...` transitioned from `CollaboratingFundingTx` to `ChannelReady`.
- **Evidence:** [`3_open_channel.png`](./evidence/3_open_channel.png), [`4_funding_tx_committed.png`](./evidence/4_funding_tx_committed.png), [`5_channel_ready.png`](./evidence/5_channel_ready.png)

### Step 5: Invoice Generation & Payment
- **Objective:** Create a payment invoice on Node 2 and pay it from Node 1.
- **Procedure:**
  1. Created a 100 CKB invoice (`fibt1...`) on Node 2 via `invoice new_invoice`.
  2. Sent payment from Node 1 using `payment send_payment --invoice ...`.
  3. Polled status via `payment get_payment --payment-hash ...`.
- **Result:** Payment confirmed immediately with `status: Success` and preimage revealed.
- **Evidence:** [`6_create_invoice.png`](./evidence/6_create_invoice.png), [`7_send_payment_success.png`](./evidence/7_send_payment_success.png)

### Step 6: Verify Balance Updates
- **Objective:** Confirm channel balance redistribution across both nodes.
- **Procedure:** Queried `channel list_channels` on both nodes with a Python Shannon-to-CKB converter.
- **Result:**
  - **Node 1:** Local balance reduced from 401.00 CKB to **301.00 CKB**; Remote balance = **100.00 CKB**.
  - **Node 2:** Local balance increased from 0.00 CKB to **100.00 CKB**; Remote balance = **301.00 CKB**.
- **Evidence:** [`9_balance_after_payment-1.png`](./evidence/9_balance_after_payment-1.png), [`9_balance_after_payment-2.png`](./evidence/9_balance_after_payment-2.png)

---

## 6. Week 3 Development Log

| Date | Activity | Result | Evidence |
|---|---|---|---|
| 29 Sep 2026 | Fiber Network architecture research | Completed | Fiber Documentation |
| 01 Oct 2026 | FNN v0.9.1 installation & environment setup | Completed | [`0_install_fiber.png`](./evidence/0_install_fiber.png) |
| 03 Oct 2026 | Dual-node configuration & port separation | Completed | [`1_node_running.png`](./evidence/1_node_running.png) |
| 04 Oct 2026 | Faucet funding & P2P handshake | Completed | [`2_connect_peer.png`](./evidence/2_connect_peer.png) |
| 05 Oct 2026 | Channel funding (500 CKB) -> `ChannelReady` | Completed | [`5_channel_ready.png`](./evidence/5_channel_ready.png) |
| 06 Oct 2026 | Invoice creation (100 CKB) & payment execution | Completed | [`7_send_payment_success.png`](./evidence/7_send_payment_success.png) |
| 07 Oct 2026 | Balance verification & weekly documentation | Completed | [`9_balance_after_payment-1.png`](./evidence/9_balance_after_payment-1.png) |

---

## 7. Challenges & Solutions

- **Port Collisions on Single Host:** Running two nodes locally required custom port mapping. Resolved by configuring `node2` with RPC `8237` and P2P `8238`.
- **FNN Key Formatting:** `ckb-cli account export` includes headers and `0x`. Stripped `0x` with `sed '1s/^0x//'` and set `chmod 600` for FNN compatibility.
- **Asynchronous Payment Check:** `send_payment` initially returns `Created`. Verified completion by polling `get_payment` until status changed to `Success`.
- **Hex Shannon Conversion:** Balances are reported in hex Shannons (e.g. `0x702198d00`). Used a lightweight Python script to convert hex Shannons to readable CKB amounts.

---

## 8. Final Reflection

Week 3 provided practical experience with CKB Layer 2 payment channels. Setting up dual Fiber nodes, establishing on-chain funding, and executing instant off-chain payments demonstrated how Fiber resolves L1 throughput limitations while preserving CKB's security model.

---

## 9. Week 4 Goals

- Test **multi-hop routing** across 3 nodes (`A -> B -> C`).
- Experiment with **UDT / Stablecoin transfers** over Fiber channels.
- Integrate Fiber payment APIs with a frontend/dApp using Fiber SDK.

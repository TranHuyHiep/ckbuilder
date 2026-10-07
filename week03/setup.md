# Fiber Network (FNN) Dual-Node Setup & Basic Transfer Guide

> Reference: [Fiber World Quick Start - Basic Transfer](https://www.fiber.world/docs/quick-start/basic-transfer)  
> Environment: Ubuntu Linux (ARM64 on VMware) / macOS Host  
> Fiber Version: `v0.9.1` | Network: Nervos Meepo Testnet (`Fibt`)

---

## 1. Overview & Architecture

Fiber Network is an off-chain Layer 2 Payment Channel Network built on top of Nervos CKB, providing instant, low-cost, and scalable payments (analogous to the Bitcoin Lightning Network).

In this lab, we set up **two independent Fiber Network Nodes (`node1` and `node2`)** on the same machine, connect them via P2P, establish a bidirectional payment channel funded with CKB on Testnet, and execute an off-chain payment of **100 CKB** using an invoice.

```text
+-------------------------+                           +-------------------------+
|      Node 1 (FNN)       |                           |      Node 2 (FNN)       |
|  RPC: 127.0.0.1:8227    |                           |  RPC: 127.0.0.1:8237    |
|  P2P: 0.0.0.0:8228      |                           |  P2P: 127.0.0.1:8238    |
+------------+------------+                           +------------+------------+
             |                                                     |
             |============= P2P Handshake (connect_peer) ==========|
             |                                                     |
             |------ Open Channel (500 CKB Funding on CKB L1) ---->|
             |       Channel ID: 0x9f89856444c...                  |
             |                                                     |
             |<----- Generate Invoice (100 CKB / fibt1...) --------|
             |                                                     |
             |------ Send Payment (HTLC / Off-chain routing) ----->|
             |       Payment Hash: 0xc0f8de... -> Success          |
             |                                                     |
    Balance: 301.00 CKB                                   Balance: 100.00 CKB
```

---

## 2. Prerequisites & Toolchain

- **Operating System:** Linux (Ubuntu aarch64 on VMware) or macOS
- **Network Access:** Nervos Meepo Testnet (`https://testnet.ckbapp.dev/`)
- **Dependencies:** `curl`, `bash`, `sed`, `chmod`, `python3` (for balance formatting)
- **CKB CLI:** `ckb-cli` v1.12.0+

If encountering HTTP proxy issues (such as HTTP 503), bypass proxy for local loopback:
```bash
export NO_PROXY=127.0.0.1,localhost
```

---

## 3. Step 1: Install Fiber Network Node (v0.9.1)

### 3.1 Automated Bootstrap
Run the official bootstrap script to download and install the `v0.9.1` release bundle and `ckb-cli`:

```bash
curl -sSfL https://raw.githubusercontent.com/nervosnetwork/fiber/v0.9.1/tools/install/install.sh \
  | INSTALL_REF=v0.9.1 FNN_VERSION=0.9.1 bash
```

This step places `ckb-cli` in `/usr/local/bin/ckb-cli`, downloads the release bundle into `~/.fiber`, and prepares the Testnet templates.

### 3.2 Guided Installer Setup
Run the guided installer to initialize the root testnet configuration:

```bash
~/.fiber/tools/install/install.sh ~/.fiber testnet
```

Follow the prompts to create or import a CKB key and set a wallet password.

> **Evidence:** See [`0_install_fiber.png`](./evidence/0_install_fiber.png)

---

## 4. Step 2: Configure Dual Nodes (`node1` & `node2`)

To simulate payment channels on a single virtual machine, we create two separate directory trees with their own binaries, configurations, and keypairs.

### 4.1 Create Node Directories
```bash
cd ~/.fiber

# Prepare Node 1
mkdir -p node1/ckb
cp fnn fnn-cli node1/
cp config/testnet/config.yml node1/config.yml

# Prepare Node 2
mkdir -p node2/ckb
cp fnn fnn-cli node2/
cp config/testnet/config.yml node2/config.yml
```

### 4.2 Generate and Export CKB Keys
Generate dedicated CKB keypairs for each node:

```bash
# Generate account for Node 1
ckb-cli account new

# Generate account for Node 2
ckb-cli account new
```

Export and strip the `0x` prefix for the raw 64-character private keys:

```bash
# In node1
ckb-cli account export --lock-arg <node1_lock_arg> --extended-privkey-path ./node1/ckb/exported-key
sed '1s/^0x//' ./node1/ckb/exported-key > ./node1/ckb/key
chmod 600 ./node1/ckb/key

# In node2
ckb-cli account export --lock-arg <node2_lock_arg> --extended-privkey-path ./node2/ckb/exported-key
sed '1s/^0x//' ./node2/ckb/exported-key > ./node2/ckb/key
chmod 600 ./node2/ckb/key
```

### 4.3 Fund Testnet Accounts
Get Testnet CKB from the official faucet:
- Faucet URL: [https://faucet.nervos.org/](https://faucet.nervos.org/)
- Provide the CKB address generated for `node1` and `node2`. Wait for on-chain confirmation before channel funding.

---

## 5. Step 3: Network & Port Configuration

Modify `node2/config.yml` to prevent port collisions with `node1`:

- **Node 1:**
  - RPC Port: `8227` (`127.0.0.1:8227`)
  - P2P Listening Addr: `/ip4/0.0.0.0/tcp/8228`
- **Node 2:**
  - RPC Port: `8237` (`127.0.0.1:8237`)
  - P2P Listening Addr: `/ip4/127.0.0.1/tcp/8238`

Apply changes to `node2/config.yml`:
```bash
sed -i.bak 's|/ip4/0.0.0.0/tcp/8228|/ip4/127.0.0.1/tcp/8238|' node2/config.yml
sed -i.bak 's|127.0.0.1:8227|127.0.0.1:8237|' node2/config.yml
```

Verify CKB RPC endpoint in both configs:
```yaml
ckb:
  rpc_url: "https://testnet.ckbapp.dev/"
```

---

## 6. Step 4: Start the Fiber Nodes

Open two separate terminals (or tmux sessions):

### Terminal 1 — Start Node 1
```bash
cd ~/.fiber/node1
FIBER_SECRET_KEY_PASSWORD='<password1>' RUST_LOG=info ./fnn -c config.yml -d .
# Or run ~/.fiber/start-node.sh
```

### Terminal 2 — Start Node 2
```bash
cd ~/.fiber/node2
FIBER_SECRET_KEY_PASSWORD='<password2>' RUST_LOG=info ./fnn -c config.yml -d .
```

Verify that both nodes start their actor systems (`store_actor`, `watchtower::actor`, `cch::actor`) and start logging periodic checks.

> **Evidence:** See [`1_node_running.png`](./evidence/1_node_running.png)

---

## 7. Step 5: Connect Peers (P2P Handshake)

### 7.1 Retrieve Node 2 Public Key
```bash
cd ~/.fiber/node2
./fnn-cli --url http://127.0.0.1:8237 info | grep pubkey
```
*Node 2 Pubkey:* `02d982a839b8d0f95aff01168766d3de2e44c0522627b4bddfbd317c94bbbe6fd8`

### 7.2 Connect from Node 1
```bash
cd ~/.fiber/node1
./fnn-cli peer connect_peer \
  --pubkey 02d982a839b8d0f95aff01168766d3de2e44c0522627b4bddfbd317c94bbbe6fd8 \
  --address "/ip4/127.0.0.1/tcp/8238" \
  --save false
```

> **Evidence:** See [`2_connect_peer.png`](./evidence/2_connect_peer.png)

---

## 8. Step 6: Open Payment Channel

### 8.1 Submit Channel Opening Request
Fund the channel from Node 1 with **500 CKB** (`50,000,000,000` shannons):

```bash
cd ~/.fiber/node1
./fnn-cli channel open_channel \
  --pubkey 02d982a839b8d0f95aff01168766d3de2e44c0522627b4bddfbd317c94bbbe6fd8 \
  --funding-amount 50000000000 \
  --public true
```

*Output:*
- `temporary_channel_id`: `0x6f2e6f8b967cf0638c37e5838f21981a3e829a0be914ebdbf39efef6c94be98f`

> **Evidence:** See [`3_open_channel.png`](./evidence/3_open_channel.png)

### 8.2 Monitor Funding Transaction Commitment
Check pending channel status:
```bash
cd ~/.fiber/node1
./fnn-cli channel list_channels --only-pending true
```

The state progresses through:
1. `CollaboratingFundingTx` (`AWAITING_REMOTE_TX_COLLABORATION_MSG`)
2. On-chain funding transaction committed on CKB Testnet
3. Channel finalized: `state_name: ChannelReady`
4. Generated `channel_id`: `0x9f89856444c303b811d5c7cf8941bb8c639f365e48e029d53ac95fddfab647fc`

> **Evidence:** See [`4_funding_tx_committed.png`](./evidence/4_funding_tx_committed.png) & [`5_channel_ready.png`](./evidence/5_channel_ready.png)

---

## 9. Step 7: Create Invoice on Node 2

Create an invoice on the recipient node (Node 2) for **100 CKB** (`10,000,000,000` shannons = `0x2540be400`):

```bash
cd ~/.fiber/node2
./fnn-cli --url http://127.0.0.1:8237 invoice new_invoice \
  --amount 10000000000 \
  --currency Fibt \
  --description "test invoice"
```

*Output details:*
- **Payment Hash:** `0xc0f8de49f0a886584b4f91d0c2abba21a52e46a743e067aeb4816bfd917f21d0`
- **Payee Public Key:** `02d982a839b8d0f95aff01168766d3de2e44c0522627b4bddfbd317c94bbbe6fd8`
- **Invoice Address:**  
  `fibt100000000001p4pt6sknukgfw0pnm6vkwk2qpx87v8savarzvp5y0r7h30esdd7e5v9243myx8tft4qw8rplu0xcx7ltjyncaux944axcv3kmned9x4f0lttxtw3uhv05pex5kmp2devx687v6pt30psr027n3ucssgw6g3am8d7zc87e8m4gh2u832c3f8tp7za8l5vy3hxe7lg7zxs6d7u2gk3xs30h2qqmg0cart9qceh4g67uyrufflaj8qgjhjqqemfl30eu7stpcfcqqpke9dq5w965a3uvyyth8p9y0ntp2592vafkuphs9rfvkndppyl77qtgpm7f9vu`

> **Evidence:** See [`6_create_invoice.png`](./evidence/6_create_invoice.png)

---

## 10. Step 8: Send Payment from Node 1

### 10.1 Execute Payment
From Node 1, pay the invoice:

```bash
cd ~/.fiber/node1
./fnn-cli payment send_payment \
  --invoice fibt100000000001p4pt6sknukgfw0pnm6vkwk2qpx87v8savarzvp5y0r7h30esdd7e5v9243myx8tft4qw8rplu0xcx7ltjyncaux944axcv3kmned9x4f0lttxtw3uhv05pex5kmp2devx687v6pt30psr027n3ucssgw6g3am8d7zc87e8m4gh2u832c3f8tp7za8l5vy3hxe7lg7zxs6d7u2gk3xs30h2qqmg0cart9qceh4g67uyrufflaj8qgjhjqqemfl30eu7stpcfcqqpke9dq5w965a3uvyyth8p9y0ntp2592vafkuphs9rfvkndppyl77qtgpm7f9vu
```

Initial response returns `status: Created`.

### 10.2 Poll Payment Status
Poll until settlement is verified:

```bash
cd ~/.fiber/node1
./fnn-cli payment get_payment \
  --payment-hash 0xc0f8de49f0a886584b4f91d0c2abba21a52e46a743e067aeb4816bfd917f21d0
```

*Result:*
- **Status:** `Success`
- **Fee:** `0x0`
- **Payment Preimage:** `0xead799376d24b445db2e3979568e6810703d500d6aedf1653236c1029c108169`

> **Evidence:** See [`7_send_payment_success.png`](./evidence/7_send_payment_success.png)

---

## 11. Step 9: Verify Balance Updates

Verify channel balance redistribution on both nodes using the balance conversion script:

```bash
# Helper parser
python_parser="import sys, re
for line in sys.stdin:
    m = re.search(r'(\w+_balance):\s*.*0x([0-9a-fA-F]+)', line)
    if m:
        ckb = int(m.group(2), 16) / 10**8
        print(f'{line.rstrip()} 👉 ({ckb:,.2f} CKB)')
    else:
        print(line, end='')"
```

### Check Node 1 Balances:
```bash
cd ~/.fiber/node1
./fnn-cli channel list_channels | python3 -c "$python_parser"
```
*Output:*
- `local_balance: '0x702198d00'` 👉 **(301.00 CKB)** (decreased by 100 CKB)
- `remote_balance: '0x2540be400'` 👉 **(100.00 CKB)** (increased by 100 CKB)

> **Evidence:** See [`9_balance_after_payment-1.png`](./evidence/9_balance_after_payment-1.png)

### Check Node 2 Balances:
```bash
cd ~/.fiber/node2
./fnn-cli --url http://127.0.0.1:8237 channel list_channels | python3 -c "$python_parser"
```
*Output:*
- `local_balance: '0x2540be400'` 👉 **(100.00 CKB)**
- `remote_balance: '0x702198d00'` 👉 **(301.00 CKB)**

> **Evidence:** See [`9_balance_after_payment-2.png`](./evidence/9_balance_after_payment-2.png)

---

## 12. Summary of Evidence Files

| Index | Step | Evidence Filename | Description |
|---|---|---|---|
| 0 | Setup & Install | [`0_install_fiber.png`](./evidence/0_install_fiber.png) | Automated installer executing bootstrap & guided setup on Ubuntu aarch64 |
| 1 | Node Daemon | [`1_node_running.png`](./evidence/1_node_running.png) | Node running with `store_actor` and periodic watchtower checks |
| 2 | P2P Peer | [`2_connect_peer.png`](./evidence/2_connect_peer.png) | P2P handshake established between Node 1 and Node 2 |
| 3 | Channel Funding | [`3_open_channel.png`](./evidence/3_open_channel.png) | Channel open call with `funding-amount 50000000000` |
| 4 | Commitment Check | [`4_funding_tx_committed.png`](./evidence/4_funding_tx_committed.png) | Funding tx collaboration and pending state check |
| 5 | Channel Ready | [`5_channel_ready.png`](./evidence/5_channel_ready.png) | Verification of `state_name: ChannelReady` on Node 1 |
| 6 | Create Invoice | [`6_create_invoice.png`](./evidence/6_create_invoice.png) | Generated 100 CKB Testnet invoice (`fibt1...`) on Node 2 |
| 7 | Payment Execution | [`7_send_payment_success.png`](./evidence/7_send_payment_success.png) | Instant payment confirmation with preimage revelation (`Success`) |
| 8 | Node 1 Balance | [`9_balance_after_payment-1.png`](./evidence/9_balance_after_payment-1.png) | Node 1 local balance updated: 401 CKB -> 301 CKB |
| 9 | Node 2 Balance | [`9_balance_after_payment-2.png`](./evidence/9_balance_after_payment-2.png) | Node 2 local balance updated: 0 CKB -> 100 CKB |

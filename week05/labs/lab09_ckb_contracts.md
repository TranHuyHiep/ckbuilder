# Lab 09 — Đọc Perun contract trên CKB

## Mục tiêu
- Hiểu Perun được hiện thực thế nào trên **mô hình Cell (UTXO)** của CKB, so với mô hình account của Ethereum.
- Nối kiến thức Week 04 (lock script, type script, witness, `since`) với một hệ thống thật.

Lab này đọc code Rust. Không cần build.

## Chuẩn bị
```sh
cd ~/Documents/CKB/ckbuilder/week05/workspace
git clone --depth 1 https://github.com/perun-network/perun-ckb-contract.git
cd perun-ckb-contract
ls contracts crates/perun-common
```

## Ba script chính
| Script | Gắn vào | Vai trò |
|---|---|---|
| `perun-channel-lockscript` (PCLS) | Lock của **channel cell** | Chỉ cho phép giao dịch có input thuộc về một trong hai bên |
| `perun-channel-typescript` (PCTS) | Type của **channel cell** | State machine: kiểm tra Start, Progress và Close |
| `perun-funds-lockscript` (PFLS) | Lock của các **funds cell** | Tiền chỉ được mở khóa khi có channel cell (PCTS) trong cùng giao dịch |

Ngoài ra còn `perun-vchannel-*` cho virtual channel và `sample-udt` cho SUDT dùng khi test.

## Các bước

### 1. Đọc kiểu dữ liệu Molecule
Mở `crates/perun-common/types.mol` và tìm các kiểu sau:
- `ChannelParameters`, `Participant`
- `ChannelConstants`: được lưu ở đâu? Gợi ý: đọc `main()` của PCTS.
- `ChannelState`, `ChannelStatus`
- `ChannelWitness`: union gồm `Fund`, `Abort`, `Dispute`, `VCDispute`, `Close`, `ForceClose`.
- `ChannelToken`: chứa một `OutPoint`.

### 2. Đọc entry point của PCTS
Mở `contracts/perun-channel-typescript/src/lib.rs`, hàm `main`. Vẽ sơ đồ nhánh:
```
get_channel_action()
 ├─ Start    → check_valid_start
 ├─ Progress → load_witness → check_valid_progress  (Fund | Dispute)
 └─ Close    → load_witness → check_valid_close     (Abort | Close | ForceClose)
```
Xác nhận sơ đồ bằng cách đọc `PChannelAction` trong `crates/perun-common/src/channels.rs`.

### 3. Đọc lock script
Đọc hai file ngắn:
- `contracts/perun-channel-lockscript/src/lib.rs`
- `contracts/perun-funds-lockscript/src/lib.rs`

## Bài tập

### 9.1. Ánh xạ Ethereum ↔ CKB
| Khái niệm | Ethereum (perun-eth-contracts) | CKB (perun-ckb-contract) |
|---|---|---|
| Nơi giữ tiền của kênh | Contract `AssetHolder` | ? |
| Nơi lưu state đã register | Storage của `Adjudicator` | ? |
| Hành động register và dispute | Gọi hàm `register` | ? |
| Hành động conclude | Gọi hàm `conclude` | ? |
| Chống tạo trùng kênh | Channel ID là hash của params | ? |
| Đo challenge duration | `block.timestamp` | ? |

### 9.2. Thread token
`ChannelToken` chứa một `OutPoint`. Vì sao dùng một OutPoint lại bảo đảm mỗi kênh là **duy nhất**? Liên hệ với cách hoạt động của Type ID trên CKB.

### 9.3. Vòng đời cell
Vẽ các cell trong từng giao dịch:
1. **Start:** Alice tạo channel cell và funds cell của mình.
2. **Fund:** Bob thêm funds cell. Status chuyển sang `funded = true`.
3. **Close:** cả hai ký state final. Channel cell và mọi funds cell bị tiêu, tạo cell trả tiền cho từng bên.

Với mỗi giao dịch, ghi rõ input, output, và witness thuộc loại nào trong `ChannelWitness`.

### 9.4. ForceClose
Tìm `check_normal_force_close` trong PCTS. Hàm này kiểm tra điều kiện thời gian bằng cách nào? Gợi ý: tìm `verify_time_lock_expired` và nhớ lại trường `since` ở Week 04.

## Câu hỏi kiểm tra
1. Vì sao trên CKB cần **hai loại lock script khác nhau**, một cho channel cell và một cho funds cell?
2. State on-chain của kênh nằm ở đâu: data của cell, args, hay witness?
3. README nói contract CKB dùng "Ethereum-style binary encoding" cho state đã ký. Điều này mở ra khả năng gì?

<details>
<summary>Gợi ý đáp án</summary>

- **Bài 9.1:**
  - Tiền nằm trong các **funds cell** bị khóa bởi PFLS.
  - State đã register nằm trong **cell data** của channel cell (`ChannelStatus`).
  - Register và dispute là một giao dịch tiêu channel cell cũ và tạo channel cell mới với witness `Dispute`.
  - Conclude là giao dịch với witness `Close` hoặc `ForceClose`.
  - Tính duy nhất của kênh đến từ thread token (OutPoint).
  - Thời gian dựa vào `since` và header của block.
- **Bài 9.2:** một OutPoint chỉ có thể bị tiêu **một lần**. Kênh gắn với việc tiêu OutPoint đó nên không thể tạo lại kênh trùng. Type ID dùng đúng ý tưởng này.
1. Channel cell cần được cập nhật bởi một trong hai bên. Funds cell thì phải đi cùng channel cell, để tiền không bị rút riêng lẻ.
2. `ChannelStatus` nằm trong **data** của channel cell. `ChannelConstants` nằm trong **args** của type script. Chữ ký và loại hành động nằm trong **witness**.
3. Cùng một state đã ký có thể được kiểm tra trên cả Ethereum và CKB. Đây là nền tảng cho channel **xuyên chuỗi** CKB ↔ ETH, xem ví dụ `payment-channel-ckb-eth`.
</details>

## Checklist
- [ ] Vẽ được sơ đồ nhánh của PCTS
- [ ] Hoàn thành bảng ánh xạ ETH ↔ CKB
- [ ] Vẽ được vòng đời cell qua 3 giao dịch

# Lab 10 — Payment channel trên CKB devnet

## Mục tiêu
- Chạy demo Alice ↔ Bob mở payment channel bằng **CKByte** trên CKB devnet, qua `perun-ckb-backend`.
- Quan sát channel cell và funds cell thật bằng `ckb-cli`.
- Mở rộng demo với **SUDT**.

Đây là lab khó nhất vì cần dựng toàn bộ devnet và build contract RISC-V.

## Yêu cầu môi trường

| Công cụ | Yêu cầu | Máy hiện tại |
|---|---|---|
| `ckb` | ≥ 0.201.0 | 0.119.0, **cần nâng cấp** |
| `ckb-cli` | ≥ 1.13.0 | 1.12.0, **cần nâng cấp** |
| `tmux`, `tmuxp` | có | **chưa cài** |
| `jq`, `expect`, `sed`, `awk`, `make` | có | ✅ |
| Rust 1.85 + target `riscv64imac-unknown-none-elf` | để build contract | cần kiểm tra |
| `clang-18`, `riscv64-unknown-elf-gcc` | để build contract | cần kiểm tra |
| Go | ≥ 1.23 | đã cài ở Lab 01 |

### ⚠️ Lưu ý cho macOS
Script `devnet/setup-devnet.sh` **build lại contract** bằng `contracts/setup_env.sh`. Script này dùng các đường dẫn của Linux như `/usr/riscv64-linux-gnu/include` và `clang-18`, cài bằng `apt`. Trên macOS các đường dẫn này không tồn tại, nên bước build gần như chắc chắn lỗi.

**Khuyến nghị:** chạy lab này trong một máy Linux, ví dụ Ubuntu 22.04/24.04 qua OrbStack, UTM, Multipass hoặc Docker. Trong máy Linux đó, cài theo README của `perun-ckb-contract`:
```sh
sudo apt install gcc-riscv64-unknown-elf binutils-riscv64-unknown-elf \
  libc6-dev-riscv64-cross libc6-riscv64-cross linux-libc-dev-riscv64-cross \
  tmux jq expect make
pip install tmuxp
wget https://apt.llvm.org/llvm.sh && chmod +x llvm.sh && sudo ./llvm.sh 18 && rm llvm.sh
rustup toolchain install 1.85.0 && rustup default 1.85.0
rustup target add riscv64imac-unknown-none-elf
cargo install cargo-generate
```
Tải `ckb` v0.201.0+ và `ckb-cli` v1.13.0+ từ trang Releases của `nervosnetwork/ckb` và `nervosnetwork/ckb-cli`. Chọn đúng bản cho kiến trúc CPU của máy ảo.

## Các bước

### 1. Lấy submodule contract
Submodule `devnet/contracts` dùng URL SSH. Nếu chưa có SSH key cho GitHub:
```sh
cd perun-examples/payment-channel-ckb
git -c url."https://github.com/".insteadOf=git@github.com: submodule update --init --recursive
ls devnet/contracts     # phải thấy mã nguồn perun-ckb-contract
```

### 2. Terminal 1: khởi động devnet
```sh
cd devnet
make dev
```
Lệnh này dùng `tmuxp` mở nhiều pane:
1. `setup-devnet.sh`: build contract, tạo tài khoản, sửa cấu hình dev chain.
2. `ckb run` và `ckb miner`.
3. `fund_accounts.expect`: chuyển CKB cho Alice và Bob.
4. `deploy_contracts.sh`: deploy PCLS, PCTS, PFLS và SUDT.
5. `sudt_helper.sh fund`: cấp SUDT.

Chờ khoảng 15–30 giây đến khi pane cuối in ra số dư SUDT.

### 3. Terminal 2: chạy demo
```sh
cd perun-examples/payment-channel-ckb
go run .
```
Kịch bản trong `main.go`:
- Alice mở kênh với **100 CKB**.
- Alice trả Bob **10 CKB**.
- Bob trả lại Alice **5 CKB**.
- Alice settle kênh.

## Bài tập

### 10.1. Đối chiếu số dư
Tính số dư trong kênh sau mỗi bước, rồi so với output của `printBalances`. Vì sao số dư ví cuối cùng của Alice **không** đúng bằng "ban đầu − 5"? Gợi ý: phí giao dịch, và capacity tối thiểu cần để tạo cell.

### 10.2. Soi cell trên chain
Trong lúc kênh đang mở, chèn `time.Sleep(60 * time.Second)` trước `Settle`. Dùng `ckb-cli` hoặc RPC `get_cells` để tìm:
- channel cell: lock là PCLS, type là PCTS.
- funds cell: lock là PFLS.

Đọc `code_hash` của script trong `devnet/contracts` hoặc trong output của `deploy_contracts.sh`. Ghi lại capacity của mỗi cell.

### 10.3. Thanh toán bằng SUDT
`OpenChannel` và `SendPayment` nhận `map[channel.Asset]float64`. Đọc `setup` trong `util.go` để tìm asset SUDT, rồi:
1. Mở kênh với cả CKB và SUDT.
2. Alice trả Bob một ít SUDT.
3. Kiểm tra số dư SUDT trước và sau.

### 10.4. Tranh chấp trên CKB
Áp dụng ý tưởng của Lab 06: settle mà **không** finalize. Quan sát:
- Có giao dịch nào mang witness `Dispute` không?
- Khoảng thời gian chờ trước giao dịch `ForceClose` là bao lâu?

## Câu hỏi kiểm tra
1. Vì sao một kênh trên CKB cần **nhiều cell** thay vì một contract như trên Ethereum?
2. Ai trả capacity cho channel cell, và capacity đó được trả lại cho ai khi đóng kênh?
3. Muốn dùng kênh Perun trong một dApp CKB thật, bạn cần thay `devnet` bằng gì?

<details>
<summary>Gợi ý đáp án</summary>

1. CKB không có storage chung của contract. State phải nằm trong cell. Mỗi bên cần funds cell riêng để nạp tiền độc lập, và channel cell giữ trạng thái chung.
2. Người mở kênh tạo channel cell nên trả capacity. Capacity được hoàn lại khi đóng. Hãy đọc `payment_min_capacity` trong `Participant` và logic của `check_valid_close` để biết chính xác.
3. Thay bằng testnet hoặc mainnet và dùng contract đã deploy ở đó. Cấu hình `deployment` cũng phải trỏ đúng `code_hash` và `cell_dep`.
</details>

## Checklist
- [ ] Devnet chạy, contract đã deploy
- [ ] `go run .` hoàn tất
- [ ] Tìm được channel cell và funds cell trên chain
- [ ] Thanh toán SUDT thành công
- [ ] (Nâng cao) Quan sát được đường tranh chấp

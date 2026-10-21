# Lab 01 — Chuẩn bị môi trường

## Mục tiêu
- Cài đủ công cụ để chạy các ví dụ Perun trên Ethereum local.
- Clone `perun-examples` và hiểu bố cục repo.

## Hiện trạng máy (đã kiểm tra ngày 2026-10-21)

| Công cụ | Trạng thái | Cần cho |
|---|---|---|
| `go` | ❌ chưa cài | Lab 01–10 |
| `node` | ✅ có (Homebrew) | Lab 02–07 |
| `jq`, `expect` | ✅ có | Lab 10 |
| `tmux`, `tmuxp` | ❌ chưa cài | Lab 10 |
| `ckb` | 0.119.0, cần ≥ 0.201.0 | Lab 10 |
| `ckb-cli` | 1.12.0, cần ≥ 1.13.0 | Lab 10 |

Lab 01 chỉ cần Go và Node. Phần CKB để dành cho Lab 10.

## Các bước

### 1. Cài Go (≥ 1.23)
```sh
brew install go
go version        # mong đợi: go1.23.x trở lên
```

### 2. Kiểm tra Node (≥ 20)
```sh
node --version
npm --version
```

### 3. Clone repo ví dụ
```sh
mkdir -p ~/Documents/CKB/ckbuilder/week05/workspace
cd ~/Documents/CKB/ckbuilder/week05/workspace
git clone https://github.com/perun-network/perun-examples.git
cd perun-examples
```

Một số submodule dùng URL SSH (`git@github.com:...`). Nếu bạn chưa cấu hình SSH key cho GitHub, ép dùng HTTPS cho riêng lệnh này:
```sh
git -c url."https://github.com/".insteadOf=git@github.com: submodule update --init --recursive
```

### 4. Tải dependency cho payment-channel
```sh
cd payment-channel
npm install          # cài Hardhat
go mod download      # tải go-perun, perun-eth-backend, go-ethereum
go build ./...       # build thử, không có lỗi là được
```

### 5. Khám phá bố cục repo
```sh
cd ..
ls
```

## Bài tập
1. Mở `payment-channel/go.mod`. Ghi lại phiên bản của `go-perun`, `perun-eth-backend` và `go-ethereum`.
2. Mở `payment-channel-ckb/go.mod`. So sánh phiên bản `go-perun` với ví dụ Ethereum. Dòng `replace` ở cuối file dùng để làm gì?
3. Liệt kê các thư mục `payment-channel-*`. Mỗi thư mục ứng với blockchain nào?

## Câu hỏi kiểm tra
1. Vì sao go-perun là một module riêng, còn phần Ethereum lại nằm trong `perun-eth-backend`?
2. `hardhat.config.js` tạo sẵn 3 tài khoản. Vai trò của từng tài khoản là gì?

<details>
<summary>Gợi ý đáp án</summary>

1. go-perun được thiết kế **blockchain-agnostic**. Lõi giao thức không biết gì về Ethereum. Mỗi chain có một backend riêng hiện thực các interface `Funder`, `Adjudicator` và `Wallet`, ví dụ perun-eth-backend hay perun-ckb-backend.
2. Deployer triển khai contract Adjudicator và AssetHolder. Alice và Bob là hai bên của kênh. Mỗi tài khoản có 10 ETH.
3. Bài tập 2: `replace` trỏ `ckb-sdk-go` sang bản fork của perun-network, vì backend cần các thay đổi chưa có ở bản gốc.
</details>

## Checklist
- [ ] `go version` ≥ 1.23
- [ ] `go build ./...` trong `payment-channel` không lỗi
- [ ] Đã trả lời 3 bài tập

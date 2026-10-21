# Lab 02 — Chạy payment channel đầu tiên

## Mục tiêu
- Chạy demo Alice ↔ Bob trên Ethereum local.
- Nhìn thấy vòng đời kênh: **open → fund → update off-chain → finalize → settle**.
- Tự kiểm chứng số dư trước và sau.

## Kiến thức liên quan
- Chỉ **mở kênh** và **settle** mới chạm blockchain. Các lần thanh toán ở giữa chỉ là trao đổi chữ ký giữa hai bên.
- Mỗi state có `version` tăng dần. State có version cao nhất mà cả hai đã ký là state hợp lệ.

## Các bước

### 1. Terminal 1: chạy chain local
```sh
cd ~/Documents/CKB/ckbuilder/week05/workspace/perun-examples/payment-channel
npx hardhat node
```
Chờ đến khi thấy dòng báo đang lắng nghe ở cổng 8545. Kiểm tra danh sách account in ra có đúng 3 tài khoản 10 ETH.

### 2. Terminal 2: chạy demo
```sh
cd ~/Documents/CKB/ckbuilder/week05/workspace/perun-examples/payment-channel
go run .
```

### 3. Output mong đợi (gần đúng)
```
Deploying contracts.
Setting up clients.
Client balances (ETH): [10 10]
Opening channel and depositing funds.
Sending payments...
Settling channel.
Adjudicator event: type = *channel.ConcludedEvent, client = ...
Adjudicator event: type = *channel.ConcludedEvent, client = ...
Client balances (ETH): [6.99... 13]
```
Số dư của Alice hơi nhỏ hơn 7 vì Alice trả phí gas cho giao dịch deposit và conclude.

## Bài tập

### 2.1. Đối chiếu kịch bản với code
Mở `main.go`. Điền bảng sau, trong đó số dư là số dư **trong kênh**:

| Bước | Hành động | Alice (kênh) | Bob (kênh) | Version |
|---|---|---|---|---|
| 0 | `OpenChannel(bob, 5)` | 5 | 0 | 0 |
| 1 | `chAlice.SendPayment(3)` | ? | ? | ? |
| 2 | `chBob.SendPayment(1)` | ? | ? | ? |
| 3 | `chAlice.SendPayment(1)` | ? | ? | ? |
| 4 | Settle (`IsFinal = true`) | ? | ? | ? |

### 2.2. Đếm giao dịch on-chain
Nhìn log của Hardhat ở terminal 1. Đếm số giao dịch được gửi lên chain trong suốt demo, sau đó phân loại: deploy, deposit, conclude và withdraw.

So sánh: nếu Alice và Bob trả nhau 3 lần trực tiếp on-chain thì cần mấy giao dịch?

### 2.3. Thêm log số dư trong kênh
Thêm method sau vào `client/channel.go`:
```go
// Balances trả về số dư hiện tại trong kênh của [Alice, Bob] theo ETH.
func (c PaymentChannel) Balances() []*big.Float {
	s := c.ch.State()
	return []*big.Float{
		WeiToEth(s.Allocation.Balance(0, c.currency)),
		WeiToEth(s.Allocation.Balance(1, c.currency)),
	}
}
```
Gọi `log.Println("In-channel:", chAlice.Balances(), "version", ...)` sau mỗi lần `SendPayment` trong `main.go`. Tự tìm cách lấy `version` từ `c.ch.State()`.

So sánh output với bảng ở bài 2.1.

## Câu hỏi kiểm tra
1. Vì sao Bob chỉ cần gọi `AcceptedChannel()` mà không cần gọi `OpenChannel`?
2. Vì sao cả Alice và Bob đều gọi `Settle()`?
3. Nếu bạn chạy `go run .` lần thứ hai mà không khởi động lại Hardhat, điều gì xảy ra với số dư ban đầu?

<details>
<summary>Gợi ý đáp án</summary>

- Bài 2.1: sau bước 1 là 2/3 (v1), bước 2 là 3/2 (v2), bước 3 là 2/3 (v3), bước 4 vẫn là 2/3 nhưng `IsFinal=true` (v4).
1. Alice là người đề xuất (proposer). Bob nhận đề xuất qua callback `HandleProposal`, chấp nhận, rồi đẩy kênh vào Go channel `c.channels`. `AcceptedChannel()` chỉ lấy kênh đó ra.
2. Alice gọi trước để **conclude** kênh trên chain và rút phần của mình. Bob gọi sau để **rút phần của Bob**. Mỗi bên tự withdraw tiền của mình.
3. Số dư không còn là 10/10 mà là số dư cuối của lần chạy trước. Contract cũng được deploy lại ở địa chỉ mới.
</details>

## Checklist
- [ ] Demo chạy xong, in ra hai `ConcludedEvent`
- [ ] Điền xong bảng 2.1
- [ ] Đếm được số giao dịch on-chain
- [ ] Log số dư trong kênh khớp với bảng

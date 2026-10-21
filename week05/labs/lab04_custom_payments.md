# Lab 04 — Tùy biến thanh toán

## Mục tiêu
- Tự viết kịch bản thanh toán mới.
- Đo được lợi ích về tốc độ và số giao dịch của off-chain.
- Cho phép **cả hai bên cùng nạp tiền** khi mở kênh.

## Chuẩn bị
- Lab 02 đã chạy được.
- Tạo branch để giữ bản gốc:
```sh
cd ~/Documents/CKB/ckbuilder/week05/workspace/perun-examples
git checkout -b my-labs
```

## Phần A — 100 micro-payment

Thay đoạn "Sending payments..." trong `main.go` bằng:
```go
log.Println("Sending 100 micro-payments...")
start := time.Now()
for i := 0; i < 100; i++ {
	chAlice.SendPayment(0.01)
}
elapsed := time.Since(start)
log.Printf("100 payments took %v (%.2f ms/payment)", elapsed, float64(elapsed.Milliseconds())/100)
```
Chạy lại Hardhat và `go run .`.

### Bài tập A
1. Ghi lại thời gian trung bình mỗi lần thanh toán.
2. Hardhat mine một block mỗi 5 giây. Nếu mỗi thanh toán là một giao dịch on-chain thì 100 thanh toán cần ít nhất bao lâu?
3. Dự đoán số dư cuối cùng của Alice và Bob trước khi chạy, rồi kiểm chứng.

## Phần B — Nạp tiền hai phía

Hiện tại Alice nạp 5 ETH, Bob nạp 0. Hãy để Bob nạp 2 ETH.

### B1. Sửa `OpenChannel`
Thêm tham số `peerAmount float64` và đặt số dư ban đầu thành:
```go
initAlloc.SetAssetBalances(c.currency, []channel.Bal{
	EthToWei(big.NewFloat(amount)),
	EthToWei(big.NewFloat(peerAmount)),
})
```

### B2. Chạy thử, quan sát lỗi
Gọi `alice.OpenChannel(bob.WireAddress(), 5, 2)` trong `main.go` và chạy. Kênh sẽ **không mở được**.

Tìm trong `client/handle.go` dòng nào làm Bob từ chối, và vì sao.

### B3. Sửa chính sách của Bob
Đổi điều kiện từ chối thành: Bob chỉ chấp nhận nếu phần mình phải nạp **không quá 3 ETH**.
```go
maxPeerFunding := EthToWei(big.NewFloat(3))
if lcp.FundingAgreement[assetIdx][peerIdx].Cmp(maxPeerFunding) > 0 {
	return nil, fmt.Errorf("peer funding too high")
}
```
Đồng thời thêm `return` ngay sau `r.Reject(...)` trong `HandleProposal`. Lý do được giải thích ở Lab 05.

### Bài tập B
1. Với kịch bản gốc (Alice trả 3, Bob trả 1, Alice trả 1) và nạp 5/2, tính số dư **trong kênh** và **trong ví** cuối cùng.
2. So sánh với kết quả Lab 02. Vì sao số dư ví cuối cùng gần như không đổi dù Bob đã nạp thêm?
3. Thử `OpenChannel(bob, 5, 4)`. Kết quả có đúng như chính sách mới không?

## Câu hỏi kiểm tra
1. `FundingAgreement` khác `InitBals` ở điểm nào?
2. Vì sao Bob phải tự kiểm tra proposal thay vì tin Alice?

<details>
<summary>Gợi ý đáp án</summary>

- **Phần A:** 100 × 0,01 = 1 ETH chuyển sang Bob. Trong kênh còn 4/1. Ví cuối khoảng 9/11, trừ phí gas của Alice. 100 thanh toán on-chain cần ít nhất 100 block × 5 giây, khoảng 8 phút.
- **Phần B, bài 1:** trong kênh là 5−3+1−1 = 2 cho Alice và 2+3−1+1 = 5 cho Bob. Ví cuối là 5+2 = 7 cho Alice và 8+5 = 13 cho Bob.
- **Phần B, bài 2:** Bob nạp 2 rồi rút lại 2 khi settle. Tiền nạp chỉ bị khóa tạm thời, không mất đi.
1. `InitBals` là số dư ban đầu trong kênh. `FundingAgreement` là số tiền mỗi bên **cam kết nạp** on-chain. Trong ví dụ này hai giá trị bằng nhau.
2. Proposal là message off-chain do Alice tạo. Nếu Bob không kiểm tra, Alice có thể buộc Bob khóa một khoản tiền lớn hoặc dùng asset lạ.
</details>

## Checklist
- [ ] Đo được thời gian 100 micro-payment
- [ ] Mở được kênh nạp 5/2
- [ ] Kênh 5/4 bị từ chối đúng chính sách

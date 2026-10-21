# Lab 05 — Kiểm tra hợp lệ và săn bug

## Mục tiêu
- Hiểu **ai kiểm tra cái gì** khi có một update: state machine của go-perun hay logic của ứng dụng.
- Tự tấn công kênh của mình để thấy các lớp bảo vệ hoạt động.
- Tìm và sửa một bug có thật trong code ví dụ.

## Kiến thức liên quan
go-perun tự kiểm tra một số quy tắc trong `channel/machine.go`:
- Tổng tiền phải được bảo toàn qua `EqualSum`.
- State đã `IsFinal` thì không được update nữa.
- `version` phải tăng.

Logic ứng dụng, ví dụ "chỉ chấp nhận update làm tăng số dư của tôi", nằm trong `HandleUpdate` do bạn viết.

## Thí nghiệm 1 — In tiền từ hư không
Thêm vào `client/channel.go`:
```go
// MintForMe thử tăng số dư của chính mình mà không trừ của ai.
func (c PaymentChannel) MintForMe(amount float64) error {
	return c.ch.Update(context.TODO(), func(state *channel.State) {
		me := c.ch.Idx()
		bal := state.Allocation.Balance(me, c.currency)
		newBal := new(big.Int).Add(bal, EthToWei(big.NewFloat(amount)))
		state.Allocation.SetBalance(me, c.currency, newBal)
	})
}
```
Gọi trong `main.go` sau khi mở kênh:
```go
if err := chAlice.MintForMe(100); err != nil {
	log.Println("Mint rejected:", err)
}
```
**Câu hỏi:** lỗi đến từ phía Alice hay phía Bob? Tìm câu thông báo lỗi trong go-perun bằng `grep`.

## Thí nghiệm 2 — Lấy tiền của Bob
```go
// Steal chuyển tiền từ peer sang mình.
func (c PaymentChannel) Steal(amount float64) error {
	return c.ch.Update(context.TODO(), func(state *channel.State) {
		me := c.ch.Idx()
		state.Allocation.TransferBalance(1-me, me, c.currency, EthToWei(big.NewFloat(amount)))
	})
}
```
Sau khi Alice đã trả Bob 3 ETH, cho Alice gọi `chAlice.Steal(2)`.

Tổng tiền vẫn được bảo toàn, nên state machine không chặn. Người chặn phải là `HandleUpdate` của Bob.

**Quan sát kỹ:** chương trình có chạy đúng như bạn nghĩ không?

## Thí nghiệm 3 — Săn bug
Mở `client/handle.go` và đọc đoạn cuối của `HandleUpdate`:
```go
if err != nil {
	r.Reject(context.TODO(), err.Error())
}

// Send the acceptance message.
err = r.Accept(context.TODO())
```
1. Khi update không hợp lệ, sau khi `Reject` thì code làm gì tiếp?
2. Đoạn tương tự trong `HandleProposal` có cùng lỗi không? Khi proposal không hợp lệ, biến `lcp` có giá trị gì, và dòng `lcp.Accept(...)` sẽ gây ra chuyện gì?
3. Sửa cả hai chỗ bằng cách thêm `return` sau `Reject`.
4. Chạy lại Thí nghiệm 2. Lần này Alice nhận được lỗi gì?

## Bài tập
1. **Xử lý lỗi gọn gàng:** đổi `SendPayment` để trả về `error` thay vì `panic`, và cập nhật `main.go` cho phù hợp.
2. **Thêm luật mới cho Bob:** Bob từ chối mọi khoản thanh toán **nhỏ hơn 0,001 ETH**, coi đó là spam. Viết luật trong `HandleUpdate` và thử bằng `SendPayment(0.0001)`.
3. **Thí nghiệm sau final:** sau khi `Settle` đã đặt `IsFinal = true`, thử gọi thêm `SendPayment`. Ghi lại lỗi.

## Câu hỏi kiểm tra
1. Phân biệt hai lớp kiểm tra: lớp nào do go-perun bảo đảm, lớp nào do ứng dụng tự viết?
2. Nếu Bob có bug và vô tình ký một state bất lợi cho mình, hậu quả có sửa được không? Vì sao?

<details>
<summary>Gợi ý đáp án</summary>

- **Thí nghiệm 1:** lỗi xảy ra ngay phía Alice, trước khi gửi cho Bob. Máy trạng thái kiểm tra `EqualSum` trong `channel/machine.go`.
- **Thí nghiệm 3:** đây là bug thật trong repo ví dụ. Sau `Reject`, code vẫn chạy tiếp tới `Accept`.
  - Trong `HandleUpdate`, responder bị gọi hai lần.
  - Trong `HandleProposal`, `lcp` là `nil` nên `lcp.Accept(...)` gây **nil pointer panic**.
  - Hành vi cụ thể phụ thuộc vào cách go-perun xử lý lần gọi thứ hai. Hãy ghi lại những gì bạn thấy.
1. go-perun bảo đảm các bất biến chung: bảo toàn tổng, version tăng, không sửa state final, chữ ký hợp lệ. Ứng dụng tự quyết định chính sách kinh doanh, ví dụ ai được trả cho ai và mức tối thiểu bao nhiêu.
2. Không sửa được. Chữ ký của Bob trên state đó là hợp lệ, và state có version cao nhất sẽ thắng khi tranh chấp. Vì thế logic kiểm tra trong `HandleUpdate` là tuyến phòng thủ quan trọng nhất của Bob.
</details>

## Checklist
- [ ] Giải thích được nguồn gốc lỗi của Thí nghiệm 1
- [ ] Tìm ra và sửa bug thiếu `return` ở cả hai handler
- [ ] Hoàn thành 3 bài tập

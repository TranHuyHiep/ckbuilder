# Lab 06 — Tranh chấp on-chain (Dispute)

## Mục tiêu
- Thấy con đường settle **không hợp tác**: register state lên chain, chờ challenge duration rồi conclude.
- Phân biệt các event của adjudicator: `RegisteredEvent`, `ProgressedEvent`, `ConcludedEvent`.

## Kiến thức liên quan
Trong `go-perun/client/adjudicate.go`, `Settle` hoạt động như sau:
- Nếu state **đã final**: conclude ngay, gọi là đường nhanh (fast path).
- Nếu state **chưa final**: gọi `ensureRegistered` để đưa state mới nhất lên chain. Sau đó phải **chờ hết challenge duration** rồi mới conclude được.

Trong thời gian chờ, watcher của phía kia có thể nộp state có version cao hơn nếu state vừa register đã cũ.

## Các bước

### 1. Mô phỏng "Bob không hợp tác"
Tạo method settle mà không finalize trong `client/channel.go`:
```go
// ForceSettle settle mà không cần peer đồng ý finalize.
func (c PaymentChannel) ForceSettle() {
	start := time.Now()
	if err := c.ch.Settle(context.TODO(), false); err != nil {
		panic(err)
	}
	log.Printf("Force settle took %v", time.Since(start))
	c.ch.Close()
}
```
Trong `main.go`, thay `chAlice.Settle()` bằng `chAlice.ForceSettle()`. Giữ nguyên `chBob.Settle()` nhưng chuyển nó sang dạng **settle phụ** để chỉ rút tiền:
```go
chAlice.ForceSettle()
chBob.Settle()
```
Nếu `chBob.Settle()` cố finalize trên kênh đã đóng và báo lỗi, hãy tự viết một hàm chỉ gọi `c.ch.Settle(ctx, true)` cho Bob. Tham số `secondary = true` nghĩa là "người khác đã conclude, tôi chỉ rút tiền".

### 2. Log chi tiết event
Sửa `HandleAdjudicatorEvent` trong `client/handle.go`:
```go
func (c *PaymentClient) HandleAdjudicatorEvent(e channel.AdjudicatorEvent) {
	switch ev := e.(type) {
	case *channel.RegisteredEvent:
		log.Printf("[%v] REGISTERED version=%d timeout=%v", c.account, ev.State.Version, ev.Timeout())
	case *channel.ProgressedEvent:
		log.Printf("[%v] PROGRESSED version=%d by idx=%d", c.account, ev.State.Version, ev.Idx)
	case *channel.ConcludedEvent:
		log.Printf("[%v] CONCLUDED", c.account)
	default:
		log.Printf("[%v] event %T", c.account, e)
	}
}
```

### 3. Chạy và quan sát
Khởi động lại Hardhat, rồi `go run .`.

## Bài tập
1. Ghi lại thứ tự các event cho **mỗi** client. So sánh với Lab 02, nơi chỉ có `ConcludedEvent`.
2. Ghi lại thời gian `Force settle took ...`. Thời gian này liên quan thế nào đến `challengeDuration = 10` và block time 5 giây?
3. Đổi `challengeDuration` thành 30 và chạy lại. Thời gian settle thay đổi ra sao?
4. Đếm số giao dịch on-chain trong đường tranh chấp và so với đường nhanh ở Lab 02.

## Câu hỏi kiểm tra
1. Vì sao đường tranh chấp vẫn **an toàn** cho Bob dù Bob không làm gì?
2. Nếu Alice register một state cũ, ví dụ version 1 khi Alice còn nhiều tiền, ai sẽ phản ứng và bằng cách nào?
3. Tại sao các ứng dụng thực tế luôn cố gắng đi đường nhanh?

<details>
<summary>Gợi ý đáp án</summary>

- **Bài 1 và 2:** đường tranh chấp có `RegisteredEvent` trước, sau đó mới có `ConcludedEvent`. Thời gian settle ít nhất bằng challenge duration cộng thời gian mine vài block.
1. Bob giữ chữ ký trên state mới nhất. Watcher của Bob theo dõi chain, và state được register phải có chữ ký của cả hai bên.
2. Watcher của Bob thấy `RegisteredEvent` với version cũ hơn state Bob đang giữ. Nó tự động register lại state mới hơn trong thời gian challenge. State có version cao nhất thắng.
3. Đường tranh chấp chậm vì phải chờ challenge duration. Nó cũng tốn nhiều giao dịch và phí gas hơn.
</details>

## Checklist
- [ ] Thấy `RegisteredEvent` trong log
- [ ] Đo được thời gian force settle với hai giá trị challenge duration
- [ ] Giải thích được vai trò của watcher

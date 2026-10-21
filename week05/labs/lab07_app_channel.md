# Lab 07 — App channel: Tic-Tac-Toe

## Mục tiêu
- Chạy một **state channel tổng quát**, nơi state chứa dữ liệu ứng dụng chứ không chỉ số dư.
- Hiểu interface `App` của go-perun: `ValidInit`, `ValidTransition`, `DecodeData`.
- Dùng `ForceUpdate` để tiến hành nước đi trên chain khi đối thủ không hợp tác.

## Các bước

### 1. Chạy demo
```sh
cd ~/Documents/CKB/ckbuilder/week05/workspace/perun-examples/app-channel
npm install
npx hardhat node          # terminal 1
go run .                  # terminal 2
```
Alice và Bob mỗi người đặt cược 5 ETH. Sau 6 nước đi off-chain, Alice đi nước thắng **on-chain** bằng `ForceSet(1, 2)`. Alice nhận toàn bộ 10 ETH.

### 2. Đọc code theo thứ tự
| File | Cần hiểu |
|---|---|
| `app/data.go` | `TicTacToeAppData`: `NextActor`, `Grid[9]`, encode và decode |
| `app/app.go` | `ValidInit`, `ValidTransition`, `Set` |
| `app/util.go` | `CheckFinal`, `computeFinalBalances` |
| `client/channel.go` | `Set` dùng `ch.Update`, `ForceSet` dùng `ch.ForceUpdate` |
| `contracts/TicTacToeApp.sol` | Bản Solidity của cùng luật chơi |

### 3. Liệt kê các luật trong `ValidTransition`
Viết ra **tất cả** các điều kiện mà một nước đi phải thỏa. Gợi ý: có ít nhất 7 điều kiện.

## Bài tập

### 7.1. Phạm luật có chủ đích
Lần lượt thử từng kịch bản trong `main.go`, mỗi lần một kịch bản. Ghi lại thông báo lỗi:

| Kịch bản | Code | Lỗi mong đợi trong `ValidTransition` |
|---|---|---|
| Đi hai lần liên tiếp | `appAlice.Set(2,0); appAlice.Set(0,0)` | ? |
| Ghi đè ô đã có | `appAlice.Set(2,0); appBob.Set(2,0)` | ? |
| Đi khi chưa đến lượt | Bob đi trước | ? |

Câu hỏi: lỗi xảy ra ở máy của **người đi sai** hay của **đối thủ**? Vì sao?

### 7.2. Hòa cờ và bug thứ hai
Viết kịch bản 9 nước dẫn đến hòa, ví dụ:
```
X O X
X O O
O X X
```
Trước khi chạy, đọc kỹ đoạn "Check all set" trong `CheckFinal` ở `app/util.go`:
```go
for _, v := range d.Grid {
	if v != notSet {
		return false, nil
	}
}
return true, nil
```
1. Đoạn này trả về "chưa kết thúc" khi nào? Điều kiện có đúng với ý "mọi ô đã được đánh" không?
2. Chạy kịch bản hòa. Sau nước thứ 9, `IsFinal` có được bật không? Thử gọi `Settle()` và ghi lại chuyện gì xảy ra.
3. Sửa bug trong Go.
4. Mở `checkFinal` trong `contracts/TicTacToeApp.sol`. Bản Solidity có cùng bug không?
5. **Suy nghĩ:** nếu chỉ sửa bản Go mà không sửa Solidity, một ván hòa kết thúc off-chain có còn tranh chấp on-chain được không? Liên hệ với câu hỏi kiểm tra số 3.

### 7.3. So sánh Go và Solidity
Mở `contracts/TicTacToeApp.sol` và tìm hàm tương ứng với `ValidTransition`. Vì sao phải viết **cùng một luật hai lần**, một bằng Go và một bằng Solidity?

## Câu hỏi kiểm tra
1. App channel khác payment channel ở chỗ nào trong `channel.State`?
2. `ForceUpdate` khác `Update` thế nào? Khi nào cần dùng `ForceUpdate`?
3. Nếu luật trong Go và luật trong Solidity không khớp nhau, chuyện gì có thể xảy ra?

<details>
<summary>Gợi ý đáp án</summary>

- **Bài 7.1:** máy của người đi sai đã kiểm tra `ValidTransition` trước khi gửi. Vì vậy lỗi xuất hiện ngay ở đó. Đối thủ cũng kiểm tra lại khi nhận update.
- **Bài 7.2:** điều kiện bị ngược. Code trả về "chưa kết thúc" ngay khi gặp **một ô đã đánh**, nên bàn đầy mà không có người thắng thì không bao giờ được coi là final. Bản sửa đúng là `if v == notSet`. Bản Solidity có cùng lỗi. Nếu chỉ sửa Go thì off-chain chấp nhận state hòa có `IsFinal = true`, còn contract tính ra "chưa final". Đó chính là sự lệch luật nói ở câu hỏi 3.
- **Bài 7.3:** Go dùng để kiểm tra off-chain. Solidity dùng để contract kiểm tra khi tranh chấp.
1. App channel có `state.App` khác `NoApp` và `state.Data` chứa dữ liệu ứng dụng, ở đây là bàn cờ và người đi tiếp.
2. `Update` cần chữ ký của mọi bên và chạy off-chain. `ForceUpdate` register state lên chain, rồi tiến state bằng một bên ký với contract App kiểm tra tính hợp lệ. Dùng khi đối thủ không phản hồi hoặc cố tình bỏ cuộc lúc sắp thua.
3. Một nước đi hợp lệ off-chain có thể bị contract từ chối, hoặc ngược lại. Kết quả tranh chấp khi đó sẽ khác với những gì hai bên đã đồng ý off-chain.
</details>

## Checklist
- [ ] Demo chạy xong, Alice nhận 10 ETH
- [ ] Liệt kê đủ luật trong `ValidTransition`
- [ ] Thử 3 kịch bản phạm luật
- [ ] Tìm và sửa bug hòa cờ ở cả Go và Solidity

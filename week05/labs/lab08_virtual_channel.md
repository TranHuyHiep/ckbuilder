# Lab 08 — Virtual channel

## Mục tiêu
- Hiểu cách hai người **không có kênh trực tiếp** giao dịch qua một trung gian (Ingrid).
- Chạy test virtual channel của go-perun trên **simulated backend**, không cần blockchain thật.
- Tự tính được số dư của từng kênh cha.

## Kiến thức liên quan
```
   Alice ══(ledger ch)══ Ingrid ══(ledger ch)══ Bob
     ╲                                          ╱
      ╲────────── virtual channel ─────────────╱
```
- Kênh ảo được **nạp tiền off-chain** từ hai kênh cha Alice–Ingrid và Bob–Ingrid.
- Sau khi mở, Alice và Bob giao dịch trực tiếp. Ingrid không tham gia các update.
- Khi đóng theo đường lạc quan, tiền quay về hai kênh cha. **Không có giao dịch on-chain nào** cho việc mở hay đóng kênh ảo.

> Bài giới thiệu ở `document/` ghi rằng virtual channel "chưa được triển khai". Đó là tài liệu cũ. go-perun hiện đã có `client/virtual_channel.go`.

## Các bước

### 1. Clone go-perun
```sh
cd ~/Documents/CKB/ckbuilder/week05/workspace
git clone --branch v0.15.0 --depth 1 https://github.com/hyperledger-labs/go-perun.git
cd go-perun
```

### 2. Chạy test
```sh
go test ./client -run 'TestVirtualChannels' -v -count=1
```
Kỳ vọng hai test `TestVirtualChannelsOptimistic` và `TestVirtualChannelsDispute` đều `PASS`.

### 3. Đọc code theo thứ tự
| File | Nội dung |
|---|---|
| `client/virtual_channel_test.go` | Số dư ban đầu và kỳ vọng |
| `client/test/virtualchannel.go` | Kịch bản: setup, mở kênh ảo, update, settle |
| `client/virtual_channel.go` | `NewVirtualChannelProposal`, cách kênh cha nạp tiền cho kênh ảo |
| `client/virtual_channel_settlement.go` | Cách tiền quay về kênh cha |

## Bài tập

### 8.1. Tính số dư
Từ `makeVirtualChannelSetup` trong `client/virtual_channel_test.go`:

| Kênh | Ban đầu | Sau khi khóa cho kênh ảo | Sau khi kênh ảo đóng |
|---|---|---|---|
| Alice–Ingrid `[Alice, Ingrid]` | 10 / 10 | ? / ? | 7 / 13 |
| Bob–Ingrid `[Bob, Ingrid]` | 10 / 10 | ? / ? | 13 / 7 |
| Alice–Bob (ảo) `[Alice, Bob]` | — | 5 / 5 | 2 / 8 |

Giải thích vì sao Ingrid có 13 ở một kênh và 7 ở kênh còn lại. Ingrid lãi hay lỗ?

### 8.2. Index map
Trong `NewVirtualChannelProposal`, Alice dùng `indexMapAlice = {0, 1}` và Bob dùng `indexMapBob = {1, 0}`. Giải thích bằng lời: index map trả lời câu hỏi gì?

### 8.3. Thay đổi kịch bản
Đổi `VirtualBalsUpdated` thành `{0, 10}`, nghĩa là Alice trả hết cho Bob. Tự tính lại `FinalBalsAlice` và `FinalBalsBob`, sửa test cho đúng và chạy lại đến khi `PASS`.

### 8.4. Đường tranh chấp
Đọc `TestVirtualChannelDispute` trong `client/test/virtualchannel.go`. Khi có tranh chấp, kênh nào phải lên chain: kênh ảo, kênh cha, hay cả hai?

## Câu hỏi kiểm tra
1. Ingrid phải khóa bao nhiêu tiền của chính mình để kênh ảo hoạt động?
2. Vì sao Ingrid không thể ăn cắp tiền trong kênh ảo dù là trung gian?
3. So sánh với thanh toán nhiều bước (multi-hop) của Lightning: điểm khác biệt chính là gì?

<details>
<summary>Gợi ý đáp án</summary>

- **Bài 8.1:** sau khi khóa, mỗi kênh cha còn 5/5 tự do và khóa 5+5 cho kênh ảo.
  - Ở kênh Alice–Ingrid, Ingrid đứng phía Bob: Ingrid nhận phần cuối của Bob là 8 và trả phần cuối của Alice là 2.
  - Kết quả: 5+2 = 7 cho Alice, 5+8 = 13 cho Ingrid. Kênh Bob–Ingrid đối xứng ngược lại.
  - Tổng của Ingrid vẫn là 20, không lãi không lỗ.
- **Bài 8.2:** index map cho biết mỗi người trong kênh ảo tương ứng với vị trí nào trong kênh cha.
- **Bài 8.3:** kết quả là Alice–Ingrid 5/15 và Bob–Ingrid 15/5.
1. Ingrid khóa ở mỗi kênh cha đúng bằng số dư ban đầu của người ở đầu kia kênh ảo. Ở ví dụ là 5 + 5.
2. Mọi update của kênh ảo đều cần chữ ký của Alice và Bob. Kênh cha cũng ràng buộc rằng tiền khóa chỉ được chia theo state cuối của kênh ảo.
3. Lightning định tuyến **từng khoản** thanh toán qua trung gian. Mỗi lần trả, trung gian đều phải tham gia và có thể thấy giao dịch. Kênh ảo của Perun chỉ cần trung gian **lúc mở và lúc đóng**, còn các update ở giữa là trực tiếp.
</details>

## Checklist
- [ ] Hai test virtual channel `PASS`
- [ ] Bảng số dư 8.1 đúng
- [ ] Sửa kịch bản 8.3 và test vẫn `PASS`

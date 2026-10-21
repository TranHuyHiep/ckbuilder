# Perun Framework — Giới thiệu (Ghi chú học tập, Tiếng Việt)

> **Nguồn:** [Perun Framework: Introduction](https://hyperledger-labs.github.io/perun-doc/introduction.html)
> **Bản tiếng Anh:** [perun_introduction_en.md](./perun_introduction_en.md)

---

## 1. Perun là gì?

Perun là một **dự án mã nguồn mở giúp tăng thông lượng giao dịch (throughput) của blockchain**.

- Perun chỉ dùng **một vài giao dịch on-chain** để khởi tạo, sau đó chuyển phần lớn hoạt động sang một **mạng ngang hàng (P2P) off-chain**.
- Sau khi các bên mở một số kênh (channel) cơ bản, **bất kỳ ai cũng có thể giao dịch với bất kỳ ai khác qua kênh ảo (virtual channel)**, mà không cần thiết lập thêm gì trên chuỗi.
- **Mục tiêu:** đưa blockchain đến mức sử dụng đại trà bằng cách giải quyết ba vấn đề:
  - **phí** giao dịch cao,
  - **độ trễ** cao (xác nhận chậm),
  - **thông lượng** thấp (ít giao dịch mỗi giây).

### Nguồn gốc
- Bắt đầu là dự án mã nguồn mở độc lập, hợp tác giữa:
  - **Robert Bosch GmbH**, qua dự án "Economy of Things", và
  - **nhóm Perun tại Đại học Kỹ thuật Darmstadt (TUDa)**.
- Sau đó gia nhập **Hyperledger Foundation** dưới dạng một dự án **Hyperledger Labs**.

---

## 2. Giao thức Perun (Perun Protocol)

Framework được xây dựng trên **giao thức Perun**. Đây là một tập các giao thức mật mã được các nhà nghiên cứu ở TUDa và Đại học Warsaw thiết kế và **phân tích hình thức (formally analyzed)**.

### Cách hoạt động
1. **Mở kênh.** Người dùng khóa (lock) coin vào một smart contract trên blockchain.
2. **Cập nhật off-chain.** Người dùng trao đổi trực tiếp với nhau các trạng thái mới đã ký. Mỗi lần cập nhật không cần giao dịch on-chain.
3. **Cưỡng chế on-chain khi cần.** Blockchain đảm bảo trạng thái off-chain mới nhất luôn có thể được thực thi trên chuỗi, ví dụ khi một bên gian lận hoặc mất kết nối.
4. **Đóng / tất toán (settle).** Số dư cuối cùng được ghi lại lên blockchain.

Perun hỗ trợ cả **payment channel** (chỉ theo dõi số dư) và **state channel** (trạng thái ứng dụng tùy ý và logic smart contract).

### Điểm nổi bật so với các hệ thống kênh khác như Lightning Network

| Tính năng | Ý nghĩa |
|---|---|
| **Ảo hóa state channel** | Các kênh có sẵn có thể ghép lại thành **kênh ảo**. Hai người không có kênh trực tiếp vẫn giao dịch được qua bên trung gian. Kênh ảo được mở và đóng **hoàn toàn off-chain**. *Trang perun-doc ghi là chưa triển khai, nhưng go-perun v0.15.0 đã có (`client/virtual_channel.go`).* |
| **Không phụ thuộc blockchain** | Thiết kế dạng module cho phép tích hợp với **bất kỳ blockchain nào**, kể cả sổ cái truyền thống. |
| **Khả năng liên thông** | Thiết kế không phụ thuộc chuỗi cộng với ảo hóa cho phép giao dịch và chạy smart contract **xuyên chuỗi (cross-chain)**. Các phần logging, messaging, persistence có thể tùy biến. |
| **Bảo mật cao** | Đặc tả giao thức đã được **chứng minh toán học** là an toàn bằng các phương pháp nghiên cứu hiện đại. |

### Ứng dụng
Tài chính/FinTech, giao thông (mobility), năng lượng, thương mại điện tử, viễn thông, và mọi tình huống cần **vi giao dịch (micro-transaction) trực tiếp**.

### Bài báo khoa học
- IEEE: https://ieeexplore.ieee.org/document/8835315
- ACM (CCS 2018): https://dl.acm.org/doi/10.1145/3243734.3243856
- Springer, *Multi-party Virtual State Channels*: https://www.springerprofessional.de/en/multi-party-virtual-state-channels/16720256

---

## 3. Các thành phần của Framework

Perun Framework gồm **ba thành phần phần mềm chính**:

```
                ┌───────────────────────────┐
 người dùng ─▶  │        perun-node         │  node đa người dùng (gRPC API)
                │  khóa, danh tính, mạng,   │  ◀── perun-nodecli (CLI)
                │  mở / giao dịch / settle  │
                └────────────┬──────────────┘
                             │ sử dụng
                ┌────────────▼──────────────┐
                │         go-perun          │  SDK: giao thức đề xuất kênh,
                │ (không phụ thuộc chuỗi)   │  state machine + lưu trữ,
                └────────────┬──────────────┘  watcher, connector
                             │ giao tiếp với
                ┌────────────▼──────────────┐
                │   perun-eth-contracts     │  smart contract on-chain
                │       (Ethereum)          │
                └───────────────────────────┘
```

1. **perun-eth-contracts:** các smart contract Ethereum cần thiết để chạy giao thức Perun trên chuỗi.
2. **go-perun:** SDK viết bằng Go, hiện thực phần lõi của giao thức:
   - **giao thức đề xuất state channel**, xử lý việc đề xuất và mở kênh,
   - **state machine có lưu trữ (persistence)**, theo dõi trạng thái kênh và lưu lại để không mất khi khởi động lại,
   - **watcher**, giám sát chuỗi và phản ứng nếu đối phương cố đăng ký một trạng thái cũ.
   - Có sẵn connector cho Ethereum nhưng được thiết kế **không phụ thuộc blockchain**.
3. **perun-node:** một **node đa người dùng** xây trên go-perun. Nó cho phép người dùng:
   - quản lý **khóa và danh tính**,
   - xử lý **kết nối mạng off-chain** với các peer,
   - **mở, giao dịch và tất toán** kênh.

### Trạng thái hiện tại
Bảng dưới là nội dung trang perun-doc (đã cũ):

| Thành phần | Trạng thái theo perun-doc |
|---|---|
| go-perun SDK | **State channel tổng quát đầy đủ trên Ethereum**, có persistence |
| perun-node | **Payment channel hai bên**, ID provider dạng file cục bộ, giao diện **gRPC** |
| perun-nodecli | Client tham chiếu cho gRPC API của perun-node, CLI tương tác có tự hoàn thành lệnh |

Đối chiếu với mã nguồn đang dùng trên server (`/root/perun-examples`, kiểm tra ngày 2026-10-21):

| Hạng mục | Thực tế |
|---|---|
| Phiên bản go-perun | **v0.15.0** (ví dụ Ethereum). Ví dụ CKB dùng bản pseudo-version `v0.13.1-0.20250528…` |
| Kênh ảo | **Đã có** trong go-perun (`client/virtual_channel.go`, có test) |
| Backend blockchain | Ethereum (`perun-eth-backend` v0.6.0), **CKB** (`perun-ckb-backend`), Stellar, Polkadot, ICP, và ví dụ xuyên chuỗi CKB ↔ ETH (`payment-channel-ckb-eth`) |
| Contract Ethereum | `perun-eth-backend` nhúng sẵn binding của `perun-eth-contracts`. Ví dụ payment-channel gọi `DeployAdjudicator` và `DeployETHAssetholder`, không cần clone repo contract |
| Contract CKB | `perun-ckb-contract` (Rust, RISC-V): PCLS, PCTS, PFLS, vchannel, sample-udt |
| Chain local | **Hardhat 2.29.1** thay cho `ganache-cli` trong tài liệu cũ |
| API | `ch.Update(...)` thay cho `UpdateBy`, địa chỉ dạng `map[wallet.BackendID]...` |

### Lộ trình (Roadmap) theo perun-doc
1. Kênh ảo (virtual channels): **đã xong** trong go-perun
2. Tích hợp SSI (danh tính tự chủ) với Hyperledger Aries
3. Thêm các blockchain backend khác: **đã có** CKB, Stellar, Polkadot, ICP
4. Kênh xuyên chuỗi (cross-chain channels): có ví dụ `payment-channel-ckb-eth` và `multiledger-channel`

---

## 4. Ý chính cần nhớ

- Perun là một **framework kênh Layer-2**: khóa tiền on-chain một lần, giao dịch off-chain nhiều lần, tất toán on-chain một lần.
- An toàn nhờ việc **trạng thái off-chain luôn có thể được cưỡng chế on-chain**.
- Điểm khác biệt chính là **kênh ảo** và **không phụ thuộc blockchain**.
- Trên CKB, kênh là một nhóm cell: **channel cell** (PCLS + PCTS) giữ trạng thái, các **funds cell** (PFLS) giữ tiền.
- Kiến trúc phần mềm đi từ chuỗi lên người dùng: **contracts → go-perun SDK → perun-node**.

## 5. Thuật ngữ

| Thuật ngữ | Giải thích |
|---|---|
| Off-chain | Hoạt động diễn ra giữa người dùng mà không ghi lên blockchain |
| Payment channel | Kênh chỉ theo dõi số dư giữa các bên |
| State channel | Kênh theo dõi trạng thái ứng dụng tùy ý, không chỉ số dư |
| Virtual channel | Kênh dựng trên các kênh có sẵn qua bên trung gian, mở và đóng mà không cần giao dịch on-chain |
| Watcher | Thành phần theo dõi chuỗi và thay người dùng phản đối các trạng thái cũ |
| Settle | Kết thúc kênh và chi trả số dư cuối cùng trên chuỗi |

## 6. Đọc thêm
- **Thực hành:** [5 bài lab Perun trên server](../labs/README.md)
- Trang tiếp theo trong tài liệu: [Concepts: State Channels](https://hyperledger-labs.github.io/perun-doc/concepts/state-channel.html)
- [Protocols](https://hyperledger-labs.github.io/perun-doc/concepts/protocols.html)
- GitHub: [go-perun](https://github.com/hyperledger-labs/go-perun), [perun-node](https://github.com/hyperledger-labs/perun-node), [perun-eth-contracts](https://github.com/hyperledger-labs/perun-eth-contracts), [perun-examples](https://github.com/perun-network/perun-examples), [perun-ckb-contract](https://github.com/perun-network/perun-ckb-contract), [perun-ckb-backend](https://github.com/perun-network/perun-ckb-backend)

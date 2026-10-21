# Perun Labs — Thực hành Perun theo từng bước

Bộ lab này đi tiếp sau bài đọc [Perun Introduction](../document/perun_introduction_vi.md).
Mỗi lab nhỏ, làm trong khoảng 30–90 phút, và lab sau dùng kết quả của lab trước.

## Lộ trình

| # | Lab | Chủ đề | Chain | Độ khó |
|---|---|---|---|---|
| 01 | [Chuẩn bị môi trường](./lab01_setup.md) | Go, Node, Hardhat, clone repo | Ethereum local | ⭐ |
| 02 | [Chạy payment channel đầu tiên](./lab02_first_payment_channel.md) | Vòng đời kênh: open → update → settle | Ethereum local | ⭐ |
| 03 | [Mổ xẻ kiến trúc client](./lab03_client_architecture.md) | Funder, Adjudicator, Watcher, Wire Bus | Đọc code | ⭐⭐ |
| 04 | [Tùy biến thanh toán](./lab04_custom_payments.md) | Micro-payment, nạp tiền hai phía | Ethereum local | ⭐⭐ |
| 05 | [Kiểm tra hợp lệ & săn bug](./lab05_validation_and_bugs.md) | HandleUpdate, bảo toàn tổng, bug thiếu `return` | Ethereum local | ⭐⭐⭐ |
| 06 | [Tranh chấp on-chain](./lab06_dispute.md) | Register, challenge duration, adjudicator event | Ethereum local | ⭐⭐⭐ |
| 07 | [App channel: Tic-Tac-Toe](./lab07_app_channel.md) | `ValidTransition`, `ForceUpdate` | Ethereum local | ⭐⭐⭐ |
| 08 | [Virtual channel](./lab08_virtual_channel.md) | Kênh ảo qua trung gian Ingrid | Simulated backend | ⭐⭐⭐ |
| 09 | [Đọc Perun contract trên CKB](./lab09_ckb_contracts.md) | Lock/Type script, witness, cell model | Đọc code Rust | ⭐⭐⭐ |
| 10 | [Payment channel trên CKB devnet](./lab10_ckb_devnet.md) | perun-ckb-backend, CKB + SUDT | CKB devnet | ⭐⭐⭐⭐ |

## Cấu trúc mỗi lab
- **Mục tiêu:** bạn sẽ hiểu hoặc làm được gì.
- **Các bước:** lệnh và thay đổi code cụ thể.
- **Bài tập:** việc tự làm, không có lời giải sẵn ngay bên cạnh.
- **Câu hỏi kiểm tra:** gợi ý đáp án nằm trong khối `<details>` ở cuối, hãy tự trả lời trước khi mở.
- **Checklist hoàn thành.**

## Thư mục làm việc đề xuất
```
week05/
├── document/        # ghi chú lý thuyết
├── labs/            # đề bài (thư mục này)
└── workspace/       # nơi clone code và làm bài (tạo ở Lab 01)
```

## Phiên bản tham chiếu
Các lab được viết dựa trên mã nguồn tại thời điểm 2026-10-21:

| Thành phần | Phiên bản |
|---|---|
| `perun-examples/payment-channel` | go-perun v0.15.0, perun-eth-backend v0.6.0, Go 1.23 |
| `perun-examples/payment-channel-ckb` | perun-ckb-backend v0.2.x, ckb-sdk-go v2.2 |
| CKB devnet | `ckb` ≥ 0.201.0, `ckb-cli` ≥ 1.13.0 |

> **Lưu ý:** Tài liệu trên perun-doc vẫn hướng dẫn dùng `ganache-cli` và API cũ như `UpdateBy`. Code hiện tại dùng **Hardhat** và `ch.Update(...)`, địa chỉ dạng `map[wallet.BackendID]...`. Các lab theo code hiện tại.

> **Chưa chạy thử:** Các lệnh trong lab được viết từ việc đọc mã nguồn, chưa được chạy trên máy này vì máy chưa cài Go. Nếu output khác với mô tả, hãy ghi lại sự khác biệt. Đó cũng là một phần của việc học.

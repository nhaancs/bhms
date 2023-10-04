# bhms
Hệ thống quản lý nhà trọ được phát triển với Go, Angular, Flutter, và Kubernetes.

## Tính năng dành cho chủ trọ (MVP)
- Tạo khu trọ và quản lý danh sách phòng
- Quản lý điện, nước, và các dịch vụ khác
- Tạo hợp đồng cho người đến thuê
- Tạo và chia sẻ hóa đơn cho từng phòng

## Thiết kế cơ sở dữ liệu

### Bảng `users`

| Tên cột         | Kiểu dữ liệu | Bắt buộc | Mô tả            |
|-----------------|--------------|----------|------------------|
| `id`            | uuid         | - [x]    | Primary key      |
| `fullname`      | text         | - [x]    |
| `phone`         | text         | - [x]    | Used for login   |
| `password_hash` | text         | - [x]    |                  |
| `roles`         | text[]       | - [x]    | USER, ADMIN      |
| `status`        | text         | - [x]    | ACTIVE, DISABLED |
| `created_at`    | timestamp    | - [x]    |                  |
| `updated_at`    | timestamp    |          |                  |

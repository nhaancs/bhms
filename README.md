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
|-----------------|--------------|:--------:|------------------|
| `id`            | uuid         | &check;  | Primary key      |
| `fullname`      | text         | &check;  |                  |
| `phone`         | text         | &check;  | Used for login   |
| `password_hash` | text         | &check;  |                  |
| `roles`         | text[]       | &check;  | USER, ADMIN      |
| `status`        | text         | &check;  | ACTIVE, DISABLED |
| `created_at`    | timestamp    | &check;  |                  |
| `updated_at`    | timestamp    | &cross;  |                  |

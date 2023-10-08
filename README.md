# bhms
Hệ thống quản lý nhà trọ được phát triển với Go, Flutter, Angular, Postgres, và Kubernetes.

## Bản thử nghiệm (MVP)

### Kế hoạch phát triển
- [x] ~~Xác định các tính năng cho bản thử nghiệm và thiết kế cơ sở dữ liệu~~
- [ ] **Thiết kế và phát triển API cho ứng dụng di động &larr;**
- [ ] Thiết kế và phát triển ứng dụng Android, iOS
- [ ] Thiết kế và phát triển API cho backoffice
- [ ] Thiết kế và phát triển backoffice cho đội ngũ vận hành

### Tính năng dành cho chủ trọ
- Tạo khu trọ và quản lý danh sách phòng
- Tạo hợp đồng cho người đến thuê
- Quản lý điện, nước, và các dịch vụ khác
- Quản lý thông tin thanh toán (tài khoản ngân hàng)
- Tạo và chia sẻ hóa đơn cho từng phòng


### Thiết kế cơ sở dữ liệu

#### Bảng `users`
Bảng `users` lưu các thông tin cần thiết của người dùng.

| Tên cột         | Kiểu dữ liệu | Bắt buộc | Mô tả                                   |
|-----------------|--------------|:--------:|-----------------------------------------|
| `id`            | uuid         | &check;  | Khóa chính                              |
| `first_name`    | text         | &check;  | Tên                                     |
| `last_name`     | text         | &cross;  | Họ và chữ lót                           |
| `phone`         | text         | &check;  | Dùng để đăng nhập                       |
| `password_hash` | text         | &check;  | Mật khẩu đã được mã hóa                 |
| `roles`         | text[]       | &check;  | Phân quyền: USER, ADMIN                 |
| `status`        | text         | &check;  | Trạng thái người dùng: ACTIVE, DISABLED |
| `created_at`    | timestamp    | &check;  | Thời gian khởi tạo                      |
| `updated_at`    | timestamp    | &cross;  | Thời gian cập nhật                      |

#### Bảng `divisions`
Bảng `divisions` lưu thông tin các đơn vị hành chính của Việt Nam.

| Tên cột      | Kiểu dữ liệu | Bắt buộc | Mô tả                                                                     |
|--------------|--------------|:--------:|---------------------------------------------------------------------------|
| `id`         | serial       | &check;  | Khóa chính, tự tăng                                                       |
| `name`       | text         | &check;  | Tên đơn vị hành chính                                                     |
| `code`       | integer      | &check;  | Mã đơn vị hành chính                                                      |
| `level`      | smallint     | &check;  | Cấp đơn vị hành chính - 1: Tỉnh, thành phố, 2: quận, huyện, 3: xã, phường |
| `parent_id`  | uuid         | &cross;  | ID của cấp cao hơn                                                        |
| `created_at` | timestamp    | &check;  | Thời gian khởi tạo                                                        |
| `updated_at` | timestamp    | &cross;  | Thời gian cập nhật                                                        |

#### Bảng `properties`
Bảng `properties` lưu các thông tin cần thiết của khu trọ.

| Tên cột             | Kiểu dữ liệu | Bắt buộc | Mô tả                                |
|---------------------|--------------|:--------:|--------------------------------------|
| `id`                | uuid         | &check;  | Khóa chính                           |
| `name`              | text         | &check;  | Tên khu trọ                          |
| `address_level1_id` | serial       | &check;  | ID tỉnh, thành phố                   |
| `address_level2_id` | serial       | &check;  | ID quận, huyện                       |
| `address_level3_id` | serial       | &check;  | ID xã, phường                        |
| `street`            | text         | &check;  | Số nhà và tên đường                  |
| `manager_id`        | uuid         | &check;  | ID của chủ trọ                       |
| `status`            | text         | &check;  | Trạng thái khu trọ: ACTIVE, DISABLED |
| `created_at`        | timestamp    | &check;  | Thời gian khởi tạo                   |
| `updated_at`        | timestamp    | &cross;  | Thời gian cập nhật                   |

#### Bảng `blocks`
Bảng `blocks` lưu thông tin các dãy của khu trọ.

| Tên cột       | Kiểu dữ liệu | Bắt buộc | Mô tả              |
|---------------|--------------|:--------:|--------------------|
| `id`          | uuid         | &check;  | Khóa chính         |
| `name`        | text         | &check;  | Tên dãy            |
| `property_id` | uuid         | &check;  | ID của khu trọ     |
| `created_at`  | timestamp    | &check;  | Thời gian khởi tạo |
| `updated_at`  | timestamp    | &cross;  | Thời gian cập nhật |

#### Bảng `floors`
Bảng `floors` lưu thông tin các tầng của khu trọ.

| Tên cột      | Kiểu dữ liệu | Bắt buộc | Mô tả              |
|--------------|--------------|:--------:|--------------------|
| `id`         | uuid         | &check;  | Khóa chính         |
| `name`       | text         | &check;  | Tên tầng           |
| `block_id`   | uuid         | &check;  | ID của dãy         |
| `created_at` | timestamp    | &check;  | Thời gian khởi tạo |
| `updated_at` | timestamp    | &cross;  | Thời gian cập nhật |

#### Bảng `units`
Bảng `units` lưu thông tin các phòng của khu trọ.

| Tên cột       | Kiểu dữ liệu | Bắt buộc | Mô tả              |
|---------------|--------------|:--------:|--------------------|
| `id`          | uuid         | &check;  | Khóa chính         |
| `name`        | text         | &check;  | Tên phòng          |
| `property_id` | uuid         | &check;  | ID của khu trọ     |
| `block_id`    | uuid         | &check;  | ID của dãy         |
| `floor_id`    | uuid         | &check;  | ID của tầng        |
| `created_at`  | timestamp    | &check;  | Thời gian khởi tạo |
| `updated_at`  | timestamp    | &cross;  | Thời gian cập nhật |

#### Bảng `payment_methods`
Bảng `payment_methods` lưu thông tin thanh toán (tài khoản ngân hàng) của các chủ trọ.

| Tên cột          | Kiểu dữ liệu | Bắt buộc | Mô tả                                         |
|------------------|--------------|:--------:|-----------------------------------------------|
| `id`             | bigserial    | &check;  | Khóa chính                                    |
| `name`           | text         | &check;  | Tên ngân hàng                                 |
| `account_name`   | text         | &check;  | Tên chủ tài khoản                             |
| `account_number` | text         | &check;  | Số tài khoản                                  |
| `note`           | text         | &cross;  | Ghi chú                                       |
| `enabled`        | boolean      | &check;  | Trạng thái bật/tắt của phương thức thanh toán |
| `property_id`    | uuid         | &check;  | ID của khu trọ                                |
| `created_at`     | timestamp    | &check;  | Thời gian khởi tạo                            |
| `updated_at`     | timestamp    | &cross;  | Thời gian cập nhật                            |

#### Bảng `sessions`
Bảng `sessions` lưu thông tin thuê phòng theo thời gian trên hợp đồng.

| Tên cột            | Kiểu dữ liệu | Bắt buộc | Mô tả                            |
|--------------------|--------------|:--------:|----------------------------------|
| `id`               | uuid         | &check;  | Khóa chính                       |
| `unit_id`          | uuid         | &check;  | ID của phòng                     |
| `start_at`         | timestamp    | &check;  | Thời gian bắt đầu vào ở          |
| `duration_in_days` | smallint     | &check;  | Thời hạn hợp đồng tính theo ngày |
| `rental_fee`       | bigint       | &check;  | Giá thuê (đ)                     |
| `num_of_members`   | smallint     | &check;  | Số lượng thành viên              |
| `renew_times`      | smallint     | &check;  | Số lần gia hạn hợp đồng          |
| `note`             | text         | &cross;  | Ghi chú                          |
| `enabled`          | boolean      | &check;  | Trạng thái bật/tắt của session   |
| `created_at`       | timestamp    | &check;  | Thời gian khởi tạo               |
| `updated_at`       | timestamp    | &cross;  | Thời gian cập nhật               |

#### Bảng `services`
Bảng `services` lưu thông tin các dịch vụ của khu trọ.

| Tên cột        | Kiểu dữ liệu | Bắt buộc | Mô tả                                                             |
|----------------|--------------|:--------:|-------------------------------------------------------------------|
| `id`           | uuid         | &check;  | Khóa chính                                                        |
| `property_id`  | uuid         | &check;  | ID của khu trọ                                                    |
| `name`         | text         | &check;  | Tên dịch vụ                                                       |
| `service_fee`  | bigint       | &check;  | Giá dịch vụ                                                       |
| `unit`         | text         | &check;  | Đơn vị: kwh, m3, room, member, piece, time                        |
| `invoice_type` | text         | &check;  | Đơn vị tính: PER_USAGE, PER_ROOM, PER_MEMBER, PER_PIECE, PER_TIME | 
| `enabled`      | boolean      | &check;  | Trạng thái bật/tắt dịch vụ                                        |
| `created_at`   | timestamp    | &check;  | Thời gian khởi tạo                                                |
| `updated_at`   | timestamp    | &cross;  | Thời gian cập nhật                                                |


#### Bảng `session_services`
Bảng `session_services` lưu thông tin các dịch vụ của từng session.

| Tên cột      | Kiểu dữ liệu | Bắt buộc | Mô tả                   |
|--------------|--------------|:--------:|-------------------------|
| `id`         | bigserial    | &check;  | ID của session-service  |
| `session_id` | uuid         | &check;  | ID của session          |
| `service_id` | uuid         | &check;  | ID của dịch vụ          |
 
#### Bảng `session_invoices`
Bảng `session_invoices` lưu thông tin hóa đơn của từng session.

| Tên cột                 | Kiểu dữ liệu | Bắt buộc | Mô tả                                 |
|-------------------------|--------------|:--------:|---------------------------------------|
| `id`                    | uuid         | &check;  | Khóa chính                            |
| `session_id`            | uuid         | &check;  | ID của session                        |
| `start_at`              | timestamp    | &check;  | Thời gian bắt đầu tính phí            |
| `rental_fee`            | bigint       | &check;  | Giá thuê (đ)                          |
| `additional_fee`        | bigint       | &cross;  | Phí phát sinh (đ)                     |
| `discount`              | bigint       | &cross;  | Giảm giá (đ)                          |
| `additional_fee_reason` | text         | &cross;  | Lý do phát sinh phí                   |
| `discount_reason`       | text         | &cross;  | Lý do giảm giá                        |
| `note`                  | text         | &cross;  | Ghi chú                               |
| `finished`              | boolean      | &check;  | Trạng thái chưa hoàn thành/hoàn thành |
| `enabled`               | boolean      | &check;  | Trạng thái bật/tắt của hóa đơn        |
| `created_at`            | timestamp    | &check;  | Thời gian khởi tạo                    |
| `updated_at`            | timestamp    | &cross;  | Thời gian cập nhật                    |

#### Bảng `session_service_invoices`
Bảng `session_service_invoices` lưu thông tin hóa đơn dịch vụ của từng session.

| Tên cột                | Kiểu dữ liệu | Bắt buộc | Mô tả                                                             |
|------------------------|--------------|:--------:|-------------------------------------------------------------------|
| `id`                   | uuid         | &check;  | Khóa chính                                                        |
| `invoice_id`           | uuid         | &check;  | ID của hóa đơn                                                    |
| `session_service_id`   | bigserial    | &check;  | ID của session-service                                            |
| `service_name`         | text         | &check;  | Tên dịch vụ                                                       |
| `service_unit`         | text         | &check;  | Đơn vị: kwh, m3, room, member, piece, time                        |
| `service_invoice_type` | text         | &check;  | Đơn vị tính: PER_USAGE, PER_ROOM, PER_MEMBER, PER_PIECE, PER_TIME | 
| `service_fee`          | bigint       | &check;  | Giá dịch vụ (đ)                                                   |
| `latest_index`         | integer      | &check;  | Chỉ số ghi nhận lần trước                                         |
| `current_index`        | integer      | &check;  | Chỉ số ghi nhận hiện tại                                          |
| `quantity`             | integer      | &check;  | Số lượng                                                          |
| `note`                 | integer      | &cross;  | Ghi chú                                                           |
| `total`                | bigint       | &check;  | Thành tiền                                                        |
| `created_at`           | timestamp    | &check;  | Thời gian khởi tạo                                                |
| `updated_at`           | timestamp    | &cross;  | Thời gian cập nhật                                                |
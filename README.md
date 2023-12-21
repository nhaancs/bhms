# bhms

Boarding house management system developed with Go, Flutter, Postgres, and Kubernetes.

## Minimum Viable Product (MVP)

### Development Plan
- [x] ~~Define features for the MVP and design the database~~
- [ ] **Design and develop APIs for the mobile application &larr;**
- [ ] Design and develop Android, iOS applications

### Features for Landlords
- Create boarding areas and manage room lists
- Create contracts for tenants
- Manage utilities such as electricity, water, and other services
- Manage payment information (bank accounts)
- Create and share invoices for each room.

### Database Design

#### Table `users`
The `users` table stores essential information about users.

| Column Name    | Data Type | Required | Description                           |
|----------------|-----------|:--------:|---------------------------------------|
| `id`           | uuid      | &check;  | Primary key                           |
| `first_name`   | text      | &check;  | First name                            |
| `last_name`    | text      | &cross;  | Last name                             |
| `phone`        | text      | &check;  | Used for login, unique                |
| `password_hash`| text      | &check;  | Encrypted password                    |
| `roles`        | text[]    | &check;  | Roles: USER, ADMIN                    |
| `status`       | text      | &check;  | User status: CREATED, ACTIVE, DELETED |
| `created_at`   | timestamp | &check;  | Creation timestamp                    |
| `updated_at`   | timestamp | &check;  | Update timestamp                      |

#### Table `divisions`
The `divisions` table stores information about administrative divisions in Vietnam.

| Column Name | Data Type | Required | Description                                                          |
|-------------|-----------|:--------:|----------------------------------------------------------------------|
| `id`        | serial    | &check;  | Primary key, auto-increment                                          |
| `name`      | text      | &check;  | Name of the administrative division                                  |
| `code`      | integer   | &check;  | Administrative division code                                         |
| `level`     | smallint  | &check;  | Level of administrative division - 1: Province, 2: District, 3: Ward |
| `parent_id` | serial    | &cross;  | ID of the higher-level division                                      |
| `created_at`| timestamp | &check;  | Creation timestamp                                                   |
| `updated_at`| timestamp | &check;  | Update timestamp                                                     |

#### Table `properties`
The `properties` table stores essential information about boarding areas.

| Column Name         | Data Type | Required | Description                                    |
|---------------------|-----------|:--------:|------------------------------------------------|
| `id`                | uuid      | &check;  | Primary key                                    |
| `name`              | text      | &check;  | Boarding area name                             |
| `address_level_1_id`| serial    | &check;  | ID of province                                 |
| `address_level_2_id`| serial    | &check;  | ID of district                                 |
| `address_level_3_id`| serial    | &check;  | ID of ward                                     |
| `street`            | text      | &check;  | House number and street name                   |
| `manager_id`        | uuid      | &check;  | ID of the landlord                             |
| `status`            | text      | &check;  | Boarding area status: CREATED, ACTIVE, DELETED |
| `created_at`        | timestamp | &check;  | Creation timestamp                             |
| `updated_at`        | timestamp | &check;  | Update timestamp                               |

#### Table `blocks`
The `blocks` table stores information about the blocks within the boarding area.

| Column Name   | Data Type | Required | Description                           |
|---------------|-----------|:--------:|---------------------------------------|
| `id`          | uuid      | &check;  | Primary key                           |
| `name`        | text      | &check;  | Block name                            |
| `property_id` | uuid      | &check;  | ID of the boarding area               |
| `status`      | text      | &check;  | Block status: ACTIVE, DELETED         |
| `created_at`  | timestamp | &check;  | Creation timestamp                    |
| `updated_at`  | timestamp | &check;  | Update timestamp                      |

#### Table `floors`
The `floors` table stores information about the floors within the boarding area.

| Column Name    | Data Type | Required | Description                          |
|----------------|-----------|:--------:|--------------------------------------|
| `id`           | uuid      | &check;  | Primary key                          |
| `name`         | text      | &check;  | Floor name                           |
| `block_id`     | uuid      | &check;  | ID of the block                      |
| `property_id`  | uuid      | &check;  | ID of the boarding area              |
| `status`       | text      | &check;  | Floor status: ACTIVE, DELETED        |
| `created_at`   | timestamp | &check;  | Creation timestamp                   |
| `updated_at`   | timestamp | &check;  | Update timestamp                     |

#### Table `units`
The `units` table stores information about the rooms within the boarding area.

| Column Name    | Data Type | Required | Description                       |
|----------------|-----------|:--------:|-----------------------------------|
| `id`           | uuid      | &check;  | Primary key                       |
| `name`         | text      | &check;  | Room name                         |
| `property_id`  | uuid      | &check;  | ID of the boarding area           |
| `block_id`     | uuid      | &check;  | ID of the block                   |
| `floor_id`     | uuid      | &check;  | ID of the floor                   |
| `status`       | text      | &check;  | Room status: ACTIVE, DELETED      |
| `created_at`   | timestamp | &check;  | Creation timestamp                |
| `updated_at`   | timestamp | &check;  | Update timestamp                  |

#### Table `payment_methods`
The `payment_methods` table stores payment information (bank accounts) of landlords.

| Column Name      | Data Type | Required | Description                                 |
|------------------|-----------|:--------:|---------------------------------------------|
| `id`             | bigserial | &check;  | Primary key                                 |
| `name`           | text      | &check;  | Bank name                                   |
| `account_name`   | text      | &check;  | Account holder's name                       |
| `account_number` | text      | &check;  | Account number                              |
| `note`           | text      | &cross;  | Note                                        |
| `enabled`        | boolean   | &check;  | Enable/disable status of the payment method |
| `property_id`    | uuid      | &check;  | ID of the boarding area                     |
| `status`         | text      | &check;  | Status: ACTIVE, DELETED                     |
| `created_at`     | timestamp | &check;  | Creation timestamp                          |
| `updated_at`     | timestamp | &check;  | Update timestamp                            |

#### Table `sessions`
The `sessions` table stores information about room rentals over time in contracts.

| Column Name          | Data Type | Required | Description                    |
|----------------------|-----------|:--------:|--------------------------------|
| `id`                 | uuid      | &check;  | Primary key                    |
| `unit_id`            | uuid      | &check;  | ID of the room                 |
| `start_at`           | timestamp | &check;  | Start time of occupancy        |
| `duration_in_days`   | smallint  | &check;  | Contract duration in days      |
| `rental_fee`         | bigint    | &check;  | Rental fee (VND)               |
| `num_of_members`     | smallint  | &check;  | Number of members              |
| `renew_times`        | smallint  | &check;  | Number of contract renewals    |
| `note`               | text      | &cross;  | Note                           |
| `status`             | text      | &check;  | Status: ACTIVE, ENDED, DELETED |
| `created_at`         | timestamp | &check;  | Creation timestamp             |
| `updated_at`         | timestamp | &check;  | Update timestamp               |

#### Table `services`
The `services` table stores information about the services provided in the boarding area.

| Column Name    | Data Type | Required | Description                                                        |
|----------------|-----------|:--------:|--------------------------------------------------------------------|
| `id`           | uuid      | &check;  | Primary key                                                        |
| `property_id`  | uuid      | &check;  | ID of the boarding area                                            |
| `name`         | text      | &check;  | Service name                                                       |
| `service_fee`  | bigint    | &check;  | Service fee (VND)                                                  |
| `unit`         | text      | &check;  | Unit: kwh, m3, room, member, piece, time                           |
| `invoice_type` | text      | &check;  | Billing unit: PER_USAGE, PER_ROOM, PER_MEMBER, PER_PIECE, PER_TIME |
| `status`       | text      | &check;  | Status: ACTIVE, DELETED                                            |
| `enabled`      | boolean   | &check;  | Enable/disable status of the service                               |
| `created_at`   | timestamp | &check;  | Creation timestamp                                                 |
| `updated_at`   | timestamp | &check;  | Update timestamp                                                   |

#### Table `session_services`
The `session_services` table stores information about the services registered for each session.

| Column Name      | Data Type | Required | Description                   |
|------------------|-----------|:--------:|-------------------------------|
| `id`             | bigserial | &check;  | ID of the session-service     |
| `session_id`     | uuid      | &check;  | ID of the session             |
| `service_id`     | uuid      | &check;  | ID of the service             |
| `status`         | text      | &check;  | Status: ACTIVE, DELETED       |
| `created_at`     | timestamp | &check;  | Creation timestamp            |
| `updated_at`     | timestamp | &check;  | Update timestamp              |

#### Table `session_invoices`
The `session_invoices` table stores information about invoices for each session.

| Column Name             | Data Type | Required | Description                        |
|-------------------------|-----------|:--------:|------------------------------------|
| `id`                    | uuid      | &check;  | Primary key                        |
| `session_id`            | uuid      | &check;  | ID of the session                  |
| `start_at`              | timestamp | &check;  | Start time for calculating charges |
| `rental_fee`            | bigint    | &check;  | Rental fee (VND)                   |
| `additional_fee`        | bigint    | &cross;  | Additional charges (VND)           |
| `discount`              | bigint    | &cross;  | Discount (VND)                     |
| `additional_fee_reason` | text      | &cross;  | Reason for additional charges      |
| `discount_reason`       | text      | &cross;  | Reason for discount                |
| `note`                  | text      | &cross;  | Note                               |
| `status`                | text      | &check;  | Status: ACTIVE, PAID, DELETED      |
| `created_at`            | timestamp | &check;  | Creation timestamp                 |
| `updated_at`            | timestamp | &check;  | Update timestamp                   |

#### Table `session_service_invoices`
The `session_service_invoices` table stores information about service invoices for each session.

| Column Name            | Data Type | Required | Description                                                        |
|------------------------|-----------|:--------:|--------------------------------------------------------------------|
| `id`                   | uuid      | &check;  | Primary key                                                        |
| `invoice_id`           | uuid      | &check;  | ID of the invoice                                                  |
| `session_service_id`   | bigserial | &check;  | ID of the session-service                                          |
| `service_name`         | text      | &check;  | Service name                                                       |
| `service_unit`         | text      | &check;  | Unit: kwh, m3, room, member, piece, time                           |
| `service_invoice_type` | text      | &check;  | Billing unit: PER_USAGE, PER_ROOM, PER_MEMBER, PER_PIECE, PER_TIME |
| `service_fee`          | bigint    | &check;  | Service fee (VND)                                                  |
| `latest_index`         | integer   | &check;  | Previous recorded index                                            |
| `current_index`        | integer   | &check;  | Current recorded index                                             |
| `quantity`             | integer   | &check;  | Quantity                                                           |
| `note`                 | integer   | &cross;  | Note                                                               |
| `total`                | bigint    | &check;  | Total amount (VND)                                                 |
| `status`               | text      | &check;  | Status: ACTIVE, DELETED                                            |
| `created_at`           | timestamp | &check;  | Creation timestamp                                                 |
| `updated_at`           | timestamp | &check;  | Update timestamp                                                   |

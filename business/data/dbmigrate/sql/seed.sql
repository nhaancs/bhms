INSERT INTO users (id, first_name, last_name, phone, roles, password_hash, status, created_at, updated_at) VALUES
	('5cf37266-3473-4006-984f-9325122678b7', 'Admin', 'Gopher', '0984250065', '{ADMIN,USER}', '$2a$10$1ggfMVZV6Js0ybvJufLRUOWHS5f6KneuP0XwwHpJ8L8ipdry9f2/a', 'ACTIVE', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('45b5fbd3-755f-4379-8f07-a58d4a30fa2f', 'User', 'Gopher', '0984250066', '{USER}', '$2a$10$9/XASPKBbJKVfCAZKDH.UuhsuALDr5vVm6VrYA9VFR8rccK86C1hW', 'ACTIVE', '2019-03-24 00:00:00', '2019-03-24 00:00:00')
ON CONFLICT DO NOTHING;

INSERT INTO users (id, first_name, last_name, phone, roles, password_hash, status, created_at, updated_at) VALUES
	('5cf37266-3473-4006-984f-9325122678b7', 'Admin', 'Gopher', '0984250066', '{ADMIN,USER}', '$2y$08$STUxn3m68IyESNYGYvIdcOrThrj./e8H/RXmMOaw56VesxSgThZTi', 'CREATED', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('45b5fbd3-755f-4379-8f07-a58d4a30fa2f', 'User', 'Gopher', '0984250067', '{USER}', '$2y$08$j6bZfxi8xVhFK7H2mdbeUO5taLmc98RXqusrfS7C75F3jrU6yLNxK', 'CREATED', '2019-03-24 00:00:00', '2019-03-24 00:00:00')
ON CONFLICT DO NOTHING;

INSERT INTO roles (id, name) VALUES
(1, 'admin'),
(2, 'client'),
(3, 'assessor')
ON CONFLICT (id) DO NOTHING;
INSERT INTO users (id, email) VALUES
(1, 'admin@mail.com'),
(2, 'client@mail.com'),
(3, 'assessor1@mail.com'),
(4, 'assessor2@mail.com'),
ON CONFLICT (id) DO NOTHING;
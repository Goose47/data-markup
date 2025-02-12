INSERT INTO assessment_types (id, name) VALUES
(1, 'radio'),
(2, 'checkbox'),
(3, 'select'),
(4, 'multiselect'),
(5, 'text')
ON CONFLICT (id) DO NOTHING;
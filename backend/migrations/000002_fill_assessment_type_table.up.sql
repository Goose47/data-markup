INSERT INTO assessment_types (id, name) VALUES
(1, 'radio'),
(2, 'checkbox'),
(3, 'select')
ON CONFLICT (unique_column) DO NOTHING;
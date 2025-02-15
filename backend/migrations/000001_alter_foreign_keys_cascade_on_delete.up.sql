ALTER TABLE markup_type_fields DROP CONSTRAINT fk_markup_types_fields;

ALTER TABLE markup_type_fields
    ADD CONSTRAINT fk_markup_types_fields
        FOREIGN KEY (markup_type_id) REFERENCES markup_types(id)
            ON DELETE CASCADE;

ALTER TABLE assessment_fields DROP CONSTRAINT fk_assessments_fields;

ALTER TABLE assessment_fields
    ADD CONSTRAINT fk_assessments_fields
        FOREIGN KEY (assessment_id) REFERENCES assessments(id)
            ON DELETE CASCADE;
ALTER TABLE markup_type_fields DROP CONSTRAINT fk_markup_types_fields;

ALTER TABLE markup_type_fields
    ADD CONSTRAINT fk_markup_types_fields
        FOREIGN KEY (markup_type_id) REFERENCES markup_types(id)
            ON DELETE CASCADE;
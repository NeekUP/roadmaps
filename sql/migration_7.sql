ALTER TABLE steps
    ADD COLUMN title character varying(255) COLLATE pg_catalog."default";

UPDATE steps SET title = 'change me! max 256chr';

ALTER TABLE steps
    ALTER COLUMN title SET NOT NULL;
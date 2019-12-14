ALTER TABLE topics
    ADD COLUMN tags character varying(512)[];

CREATE INDEX ix_topics_tags
    ON topics USING gin
        (tags COLLATE pg_catalog."default")
    TABLESPACE pg_default;
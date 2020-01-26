CREATE TABLE comments
(
    id bigserial NOT NULL,
    entitytype integer NOT NULL,
    entityid bigint NOT NULL,
    date timestamp without time zone NOT NULL,
    parentid bigint,
    threadid bigint,
    userid character varying(36) NOT NULL,
    text text NOT NULL,
	title character varying(256),
	deleted boolean NOT NULL,
    points integer NOT NULL DEFAULT 0,
    PRIMARY KEY (id)
)
WITH (
    OIDS = FALSE
);

ALTER TABLE comments
    OWNER to postgres;

ALTER TABLE comments
    ADD CONSTRAINT fk_comments_userid FOREIGN KEY (userid)
    REFERENCES users (id) MATCH SIMPLE
    ON UPDATE CASCADE
    ON DELETE NO ACTION
    NOT VALID;

ALTER TABLE comments
    ADD CONSTRAINT fk_parentid_id FOREIGN KEY (parentid)
    REFERENCES comments (id) MATCH SIMPLE
    ON UPDATE CASCADE
    ON DELETE CASCADE
    NOT VALID;


ALTER TABLE comments
    ADD CONSTRAINT fk_comments_threadid FOREIGN KEY (threadid)
    REFERENCES comments (id) MATCH SIMPLE
    ON UPDATE CASCADE
    ON DELETE CASCADE
    NOT VALID;
    
CREATE INDEX fki_fk_comments_threadid
    ON comments(threadid);

CREATE INDEX ix_comments_entitytype_entityid
    ON public.comments USING btree
    (entitytype ASC NULLS LAST, entityid ASC NULLS LAST)
    TABLESPACE pg_default;
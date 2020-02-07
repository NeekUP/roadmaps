CREATE TABLE changelog
(
    id bigserial NOT NULL,
    action integer NOT NULL,
    userid character varying(36) NOT NULL,
    entitytype integer NOT NULL,
    entityid bigint NOT NULL,
    diff text,
    points integer NOT NULL DEFAULT 0,
    date timestamp without time zone NOT NULL,
    PRIMARY KEY (id)
)
WITH (
    OIDS = FALSE
);
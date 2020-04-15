CREATE TABLE users_oauth
(
    userid character varying(36) NOT NULL,
    provider character varying(36) NOT NULL,
    id character varying(36) NOT NULL,
    date timestamp without time zone NOT NULL
)
WITH (
    OIDS = FALSE
);

ALTER TABLE users_oauth
    ADD CONSTRAINT fk_users_oauth_userid FOREIGN KEY (userid)
        REFERENCES users (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
        NOT VALID;
CREATE INDEX fki_fk_users_oauth_userid
    ON users_oauth(userid);

CREATE UNIQUE INDEX ux_users_oauth_provider_id
    ON users_oauth USING btree
        (provider ASC NULLS LAST, id ASC NULLS LAST);
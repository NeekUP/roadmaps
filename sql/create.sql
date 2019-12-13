CREATE TABLE plans (
    id integer NOT NULL,
    title character varying(256) NOT NULL,
    topic character varying(128) NOT NULL,
    owner character varying(36) NOT NULL,
    points integer DEFAULT 0 NOT NULL
);

CREATE SEQUENCE plans_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE plans_id_seq OWNED BY plans.id;

CREATE TABLE sources (
    id bigint NOT NULL,
    title character varying(256) NOT NULL,
    identifier character varying(256) NOT NULL,
    normalizedidentifier character varying(256) NOT NULL,
    type character varying(24) NOT NULL,
    properties character varying(1024),
    img character varying(256),
    description character varying(4096)
);

CREATE SEQUENCE sources_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE sources_id_seq OWNED BY sources.id;

CREATE TABLE steps (
    id bigint NOT NULL,
    planid integer NOT NULL,
    referenceid bigint NOT NULL,
    referencetype character varying(24) NOT NULL,
    "position" integer NOT NULL
);

CREATE SEQUENCE steps_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE steps_id_seq OWNED BY steps.id;

CREATE TABLE steps_sources (
    stepid bigint NOT NULL,
    sourceid bigint NOT NULL
);

CREATE TABLE topics (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    title character varying(100) NOT NULL,
    description character varying(1024),
    creator character varying(64) NOT NULL
);

CREATE SEQUENCE topics_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE topics_id_seq OWNED BY topics.id;

CREATE TABLE users (
    id character varying(36) NOT NULL,
    name character varying(64) NOT NULL,
    normalizedname character varying(100) NOT NULL,
    email character varying(128) NOT NULL,
    emailconfirmed boolean NOT NULL,
    emailconfirmation character varying(64),
    img character varying(64),
    tokens character varying(4096),
    rights integer NOT NULL,
    password bytea NOT NULL,
    salt bytea NOT NULL
);

CREATE TABLE usersplans (
    userid character varying(36) NOT NULL,
    topic character varying(100) NOT NULL,
    planid integer NOT NULL
);

ALTER TABLE ONLY plans ALTER COLUMN id SET DEFAULT nextval('plans_id_seq'::regclass);
ALTER TABLE ONLY sources ALTER COLUMN id SET DEFAULT nextval('sources_id_seq'::regclass);
ALTER TABLE ONLY steps ALTER COLUMN id SET DEFAULT nextval('steps_id_seq'::regclass);
ALTER TABLE ONLY topics ALTER COLUMN id SET DEFAULT nextval('topics_id_seq'::regclass);
ALTER TABLE ONLY plans
    ADD CONSTRAINT plans_pkey PRIMARY KEY (id);
ALTER TABLE ONLY sources
    ADD CONSTRAINT sources_pkey PRIMARY KEY (id);
ALTER TABLE ONLY steps
    ADD CONSTRAINT steps_pkey PRIMARY KEY (id);
ALTER TABLE ONLY steps_sources
    ADD CONSTRAINT steps_sources_pkey PRIMARY KEY (stepid, sourceid);
ALTER TABLE ONLY topics
    ADD CONSTRAINT topics_pkey PRIMARY KEY (id);
ALTER TABLE ONLY sources
    ADD CONSTRAINT u_sources_identifier UNIQUE (identifier);
ALTER TABLE ONLY sources
    ADD CONSTRAINT u_sources_normalizedidentifier UNIQUE (normalizedidentifier);
ALTER TABLE ONLY topics
    ADD CONSTRAINT u_topics_name UNIQUE (name);
ALTER TABLE ONLY users
    ADD CONSTRAINT u_users_email UNIQUE (email);
ALTER TABLE ONLY users
    ADD CONSTRAINT u_users_name UNIQUE (normalizedname);
ALTER TABLE ONLY users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);

CREATE INDEX fki_fk_plans_owner ON plans USING btree (owner);
CREATE INDEX fki_fk_plans_topic ON plans USING btree (topic);
CREATE INDEX fki_fk_steps_planid ON steps USING btree (planid);
CREATE INDEX fki_fk_steps_sources_sourceid ON steps_sources USING btree (sourceid);
CREATE INDEX fki_fk_steps_sources_stepid ON steps_sources USING btree (stepid);
CREATE INDEX fki_fk_topics_creator ON topics USING btree (creator);
CREATE INDEX fki_fk_userplans_planid ON usersplans USING btree (planid);
CREATE INDEX fki_fk_userplans_topic ON usersplans USING btree (topic);
CREATE INDEX fki_fk_usersplans_userid ON usersplans USING btree (userid);
CREATE UNIQUE INDEX ix_usersplans_userid_topic ON usersplans USING btree (userid, topic);

ALTER TABLE ONLY plans
    ADD CONSTRAINT fk_plans_owner FOREIGN KEY (owner) REFERENCES users(id) ON UPDATE CASCADE NOT VALID;
ALTER TABLE ONLY plans
    ADD CONSTRAINT fk_plans_topic FOREIGN KEY (topic) REFERENCES topics(name) ON UPDATE CASCADE NOT VALID;
ALTER TABLE ONLY steps
    ADD CONSTRAINT fk_steps_planid FOREIGN KEY (planid) REFERENCES plans(id) ON UPDATE CASCADE NOT VALID;
ALTER TABLE ONLY steps_sources
    ADD CONSTRAINT fk_steps_sources_sourceid FOREIGN KEY (sourceid) REFERENCES sources(id) ON UPDATE CASCADE NOT VALID;
ALTER TABLE ONLY steps_sources
    ADD CONSTRAINT fk_steps_sources_stepid FOREIGN KEY (stepid) REFERENCES steps(id) ON UPDATE CASCADE NOT VALID;
ALTER TABLE ONLY topics
    ADD CONSTRAINT fk_topics_creator FOREIGN KEY (creator) REFERENCES users(id) ON UPDATE CASCADE NOT VALID;
ALTER TABLE ONLY usersplans
    ADD CONSTRAINT fk_userplans_planid FOREIGN KEY (planid) REFERENCES plans(id) ON UPDATE CASCADE NOT VALID;
ALTER TABLE ONLY usersplans
    ADD CONSTRAINT fk_userplans_topic FOREIGN KEY (topic) REFERENCES topics(name) ON UPDATE CASCADE NOT VALID;
ALTER TABLE ONLY usersplans
    ADD CONSTRAINT fk_usersplans_userid FOREIGN KEY (userid) REFERENCES users(id) NOT VALID;

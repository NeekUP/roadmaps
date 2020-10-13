
CREATE SEQUENCE public.groups_id_seq
    INCREMENT 1
    START 1
    MINVALUE 1
    MAXVALUE 9223372036854775807
    CACHE 1;

ALTER SEQUENCE public.groups_id_seq
    OWNER TO postgres;

CREATE TABLE public.groups
(
    id bigint NOT NULL DEFAULT nextval('groups_id_seq'::regclass),
    name character varying(64) COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT groups_pkey PRIMARY KEY (id)
)

WITH (
    OIDS = FALSE
);

CREATE TABLE public.permissions_groups_entity
(
    groupid bigint NOT NULL,
    entitytype integer NOT NULL,
    permissions bigint NOT NULL,
    CONSTRAINT permissions_group_entity_pkey PRIMARY KEY (groupid),
    CONSTRAINT permissions_group_entity_groupid FOREIGN KEY (groupid)
        REFERENCES public.groups (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
        NOT VALID
)
WITH (
    OIDS = FALSE
);

CREATE TABLE public.permissions_base
(
    entityid bigint NOT NULL,
    permissions bigint,
    userid character varying(36) COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT permissions_base_userid FOREIGN KEY (userid)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
        NOT VALID
)
WITH (
    OIDS = FALSE
);

CREATE INDEX fki_permissions_base_userid
    ON public.permissions_base USING btree
        (userid COLLATE pg_catalog."default" ASC NULLS LAST)
    TABLESPACE pg_default;

CREATE TABLE public.permissions_plans
(
    -- Inherited from table public.permissions_base: entityid bigint NOT NULL,
    -- Inherited from table public.permissions_base: permissions bigint,
    -- Inherited from table public.permissions_base: userid character varying(36) COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT permissions_plans_plan FOREIGN KEY (entityid)
        REFERENCES public.plans (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
        NOT VALID,
    CONSTRAINT permissions_plans_userid FOREIGN KEY (userid)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
        NOT VALID
)
INHERITS (public.permissions_base)
WITH (
    OIDS = FALSE
);

CREATE INDEX fki_permissions_plans_plan
    ON public.permissions_plans USING btree
        (entityid ASC NULLS LAST)
    TABLESPACE pg_default;

CREATE TABLE public.permissions_topics
(
    -- Inherited from table public.permissions_base: entityid bigint NOT NULL,
    -- Inherited from table public.permissions_base: permissions bigint,
    -- Inherited from table public.permissions_base: userid character varying(36) COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT permissions_topics_topic FOREIGN KEY (entityid)
        REFERENCES public.topics (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
        NOT VALID,
    CONSTRAINT permissions_topics_userid FOREIGN KEY (userid)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
        NOT VALID
)
    INHERITS (public.permissions_base)
    WITH (
        OIDS = FALSE
    )
    TABLESPACE pg_default;

CREATE INDEX fki_permissions_topics_topic
    ON public.permissions_topics USING btree
        (entityid ASC NULLS LAST)
    TABLESPACE pg_default;

CREATE TABLE public.user_groups
(
    groupid bigint NOT NULL,
    userid character varying(36) COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT user_groups_pkey PRIMARY KEY (userid),
    CONSTRAINT user_groups_groupid_userid UNIQUE (groupid, userid),
    CONSTRAINT user_groups_groupid FOREIGN KEY (groupid)
        REFERENCES public.groups (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID,
    CONSTRAINT user_groups_userid FOREIGN KEY (userid)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
        NOT VALID
)
    WITH (
        OIDS = FALSE
    )
    TABLESPACE pg_default;


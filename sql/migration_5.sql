
-- points_base
CREATE TABLE points_base
(
    entityid bigint NOT NULL,
    userid character varying(36) NOT NULL,
    date timestamp(0) without time zone NOT NULL,
    value integer NOT NULL
)
WITH (
    OIDS = FALSE
);

CREATE INDEX ix_points_base
    ON points_base USING btree
    (entityid ASC NULLS LAST)
    TABLESPACE pg_default;

CREATE INDEX ix_points_base_userid
    ON points_base USING btree
    (userid ASC NULLS LAST)

ALTER TABLE points_base
    ADD CONSTRAINT u_points_base_entityid_userid UNIQUE (entityid, userid);

-- points_aggregated_base
CREATE TABLE points_aggregated_base
(
    entityid bigint NOT NULL,
    updatedate timestamp without time zone NOT NULL,
    count integer NOT NULL DEFAULT 0,
    value integer NOT NULL DEFAULT 0,
    avg double precision NOT NULL DEFAULT 0;
    PRIMARY KEY (entityid)
)
WITH (
    OIDS = FALSE
);


-- drop unnecessary columns
ALTER TABLE plans DROP COLUMN points;
ALTER TABLE comments DROP COLUMN points;

-- concrete points tables





-- PLANS
CREATE TABLE points_plans (LIKE points_base INCLUDING ALL) INHERITS (points_base);
CREATE TABLE points_aggregated_plans (LIKE points_aggregated_base INCLUDING ALL) INHERITS (points_aggregated_base);

CREATE OR REPLACE FUNCTION aggregate_plans_points()
  RETURNS trigger AS
$BODY$
DECLARE
   _avg double precision := 0;
   _count INTEGER := 0;
BEGIN
    -- minVote = 10
    -- avgConst = 7
	SELECT avg, count into _avg, _count + 1
	FROM points_aggregated_plans WHERE entityid = NEW.entityid LIMIT 1 FOR UPDATE;

	UPDATE points_aggregated_plans
	SET count = _count,
		avg = _avg + (NEW.value - _avg) / _count,
		value = ((_avg + (NEW.value - _avg) / _count) * _count + 7 * 10) /  (_count + 10)
	WHERE entityid = NEW.entityid;
	RETURN NULL;
END;
$BODY$
LANGUAGE plpgsql VOLATILE

CREATE OR REPLACE FUNCTION create_aggregate_plans_points()
  RETURNS trigger AS
$BODY$
BEGIN
	INSERT INTO points_aggregated_plans (entityid,updatedate,count,value,avg)
	VALUES (NEW.id, now(), 0, 0, 0 )
	ON CONFLICT (entityid)
	DO NOTHING;
	RETURN NULL;
END;
$BODY$
LANGUAGE plpgsql VOLATILE

CREATE TRIGGER create_aggregate_points_plans_trigger
  AFTER INSERT
  ON plans
  FOR EACH ROW
  EXECUTE PROCEDURE create_aggregate_plans_points();

CREATE TRIGGER aggregate_points_plans_trigger
  AFTER INSERT
  ON points_plans
  FOR EACH ROW
  EXECUTE PROCEDURE aggregate_plans_points();

ALTER TABLE points_plans
    ADD CONSTRAINT points_plans_entityid_plans_id FOREIGN KEY (entityid)
    REFERENCES plans (id) MATCH SIMPLE
    ON UPDATE CASCADE
    ON DELETE CASCADE
    NOT VALID;

-- COMMENTS
CREATE TABLE points_comments (LIKE points_base INCLUDING ALL) INHERITS (points_base);
CREATE TABLE points_aggregated_comments (LIKE points_aggregated_base INCLUDING ALL) INHERITS (points_aggregated_base);

CREATE OR REPLACE FUNCTION create_aggregate_comments_points()
  RETURNS trigger AS
$BODY$
BEGIN
	INSERT INTO comments_aggregated_comments (entityid,updatedate,count,value,avg)
	VALUES (NEW.id, now(), 0, 0, 0 )
	ON CONFLICT (entityid)
	DO NOTHING;
	RETURN NULL;
END;
$BODY$
LANGUAGE plpgsql VOLATILE;

CREATE OR REPLACE FUNCTION aggregate_comments_points()
  RETURNS trigger AS
$BODY$
DECLARE
   _avg double precision := 0;
   _count INTEGER := 0;
BEGIN
    -- minVote = 10
    -- avgConst = 7
	SELECT avg, count into _avg, _count + 1
	FROM points_aggregated_comments WHERE entityid = NEW.entityid LIMIT 1 FOR UPDATE;

	UPDATE points_aggregated_comments
	SET count = _count,
		avg = _avg + (NEW.value - _avg) / _count,
		value = ((_avg + (NEW.value - _avg) / _count) * _count + 7 * 10) /  (_count + 10)
	WHERE entityid = NEW.entityid;
	RETURN NULL;
END;
$BODY$
LANGUAGE plpgsql VOLATILE;

CREATE TRIGGER create_aggregate_points_comments_trigger
  AFTER INSERT
  ON comments
  FOR EACH ROW
  EXECUTE PROCEDURE create_aggregate_comments_points();

CREATE TRIGGER aggregate_points_comments_trigger
  AFTER INSERT
  ON points_comments
  FOR EACH ROW
  EXECUTE PROCEDURE aggregate_comments_points();

ALTER TABLE points_comments
    ADD CONSTRAINT points_comments_entityid_commants_id FOREIGN KEY (entityid)
    REFERENCES comments (id) MATCH SIMPLE
    ON UPDATE CASCADE
    ON DELETE CASCADE
    NOT VALID;

-- PROJECTS
CREATE TABLE points_projects (LIKE points_base INCLUDING ALL) INHERITS (points_base);
CREATE TABLE points_aggregated_projects (LIKE points_aggregated_base INCLUDING ALL) INHERITS (points_aggregated_base);

CREATE OR REPLACE FUNCTION create_aggregate_projects_points()
    RETURNS trigger AS
$BODY$
BEGIN
    INSERT INTO points_aggregated_projects (entityid,updatedate,count,value,avg)
    VALUES (NEW.id, now(), 0, 0, 0 )
    ON CONFLICT (entityid)
        DO NOTHING;
    RETURN NULL;
END;
$BODY$
    LANGUAGE plpgsql VOLATILE;

CREATE OR REPLACE FUNCTION aggregate_projects_points()
    RETURNS trigger AS
$BODY$
DECLARE
    _avg double precision := 0;
    _count INTEGER := 0;
BEGIN
    -- minVote = 10
    -- avgConst = 7
    SELECT avg, count into _avg, _count + 1
    FROM points_aggregated_projects WHERE entityid = NEW.entityid LIMIT 1 FOR UPDATE;

    UPDATE points_aggregated_projects
    SET count = _count,
        avg = _avg + (NEW.value - _avg) / _count,
        value = ((_avg + (NEW.value - _avg) / _count) * _count + 7 * 10) /  (_count + 10)
    WHERE entityid = NEW.entityid;
    RETURN NULL;
END;
$BODY$
    LANGUAGE plpgsql VOLATILE;

CREATE TRIGGER create_aggregate_points_projects_trigger
    AFTER INSERT
    ON projects
    FOR EACH ROW
EXECUTE PROCEDURE create_aggregate_projects_points();

CREATE TRIGGER aggregate_points_projects_trigger
    AFTER INSERT
    ON points_projects
    FOR EACH ROW
EXECUTE PROCEDURE aggregate_projects_points();

ALTER TABLE points_projects
    ADD CONSTRAINT points_projects_entityid_commants_id FOREIGN KEY (entityid)
        REFERENCES projects (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
        NOT VALID;
-- CHANGELOG
CREATE TABLE points_changelog (LIKE points_base INCLUDING ALL) INHERITS (points_base);
CREATE TABLE points_aggregated_changelog (LIKE points_aggregated_base INCLUDING ALL) INHERITS (points_aggregated_base);

CREATE OR REPLACE FUNCTION create_aggregate_changelog_points()
    RETURNS trigger AS
$BODY$
BEGIN
    INSERT INTO points_aggregated_changelog (entityid,updatedate,count,value,avg)
    VALUES (NEW.id, now(), 0, 0, 0 )
    ON CONFLICT (entityid)
        DO NOTHING;
    RETURN NULL;
END;
$BODY$
    LANGUAGE plpgsql VOLATILE;

CREATE OR REPLACE FUNCTION aggregate_changelog_points()
    RETURNS trigger AS
$BODY$
DECLARE
    _avg double precision := 0;
    _count INTEGER := 0;
BEGIN
    -- minVote = 10
    -- avgConst = 7
    SELECT avg, count into _avg, _count + 1
    FROM points_aggregated_changelog WHERE entityid = NEW.entityid LIMIT 1 FOR UPDATE;

    UPDATE points_aggregated_changelog
    SET count = _count,
        avg = _avg + (NEW.value - _avg) / _count,
        value = ((_avg + (NEW.value - _avg) / _count) * _count + 7 * 10) /  (_count + 10)
    WHERE entityid = NEW.entityid;
    RETURN NULL;
END;
$BODY$
    LANGUAGE plpgsql VOLATILE;

CREATE TRIGGER create_aggregate_points_changelog_trigger
    AFTER INSERT
    ON changelog
    FOR EACH ROW
EXECUTE PROCEDURE create_aggregate_changelog_points();

CREATE TRIGGER aggregate_points_changelog_trigger
    AFTER INSERT
    ON points_changelog
    FOR EACH ROW
EXECUTE PROCEDURE aggregate_changelog_points();

ALTER TABLE points_changelog
    ADD CONSTRAINT points_changelog_entityid_changelog_id FOREIGN KEY (entityid)
        REFERENCES changelog (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
        NOT VALID;
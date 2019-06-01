
-----------------------------| apps |-----------------------------
DROP TABLE IF EXISTS apps CASCADE;
CREATE TABLE apps
(
    aid             BIGSERIAL   PRIMARY KEY,
    link            TEXT        NOT NULL UNIQUE,
    name            TEXT        NOT NULL UNIQUE,
    image           TEXT        NOT NULL,
    url             TEXT        NOT NULL DEFAULT '',
    about           TEXT        NOT NULL,
    installations   BIGINT      NOT NULL DEFAULT 0,
    category        TEXT        NOT NULL
);

-----------------------------| installed apps |-----------------------------
DROP TABLE IF EXISTS userapps;
CREATE TABLE userapps (
    uid     BIGINT,
    aid     BIGINT,

    FOREIGN KEY (uid) REFERENCES users(uid),
    FOREIGN KEY (aid) REFERENCES apps(aid),
    PRIMARY KEY (uid, aid)
);

CREATE OR REPLACE FUNCTION update_installations() RETURNS TRIGGER AS '
    BEGIN
        UPDATE apps
            SET installations = installations + 1
            WHERE aid = NEW.aid;
            RETURN NULL;
    END;
    ' LANGUAGE plpgsql;

CREATE TRIGGER update_installations
    AFTER INSERT ON userapps FOR EACH ROW 
    EXECUTE PROCEDURE update_installations();

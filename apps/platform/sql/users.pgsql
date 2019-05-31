
-----------------------------| user |-----------------------------
DROP TABLE IF EXISTS users CASCADE;
CREATE TABLE users
(
    uid         BIGINT      PRIMARY KEY,        --\e\-- comes from auth microservice
    name        TEXT        NOT NULL UNIQUE,    --\e\-- comes from auth microservice
    avatar      TEXT        DEFAULT 'img/avatars/default'
);

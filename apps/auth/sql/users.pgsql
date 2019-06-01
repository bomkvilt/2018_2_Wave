-----------------------------| user |-----------------------------
DROP TABLE IF EXISTS users CASCADE;
CREATE TABLE users
(
    uid         BIGSERIAL   PRIMARY KEY,
    name        TEXT        NOT NULL UNIQUE,
    password    TEXT        NOT NULL
);

-----------------------------| session |-----------------------------
DROP TABLE IF EXISTS sessions;
CREATE TABLE sessions
(
    uid         BIGINT      NOT NULL,
    cookie      TEXT        PRIMARY KEY,

    FOREIGN KEY (uid) REFERENCES users(uid)
);

-----------------------------| admins |-----------------------------
DROP TABLE IF EXISTS admins;
CREATE TABLE admins
(
    uid         PRIMARY KEY,

    FOREIGN KEY (uid) REFERENCES users(uid)
);

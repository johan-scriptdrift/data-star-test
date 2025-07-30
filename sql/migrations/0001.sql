CREATE TABLE users
(
    id            INTEGER PRIMARY KEY,
    first_name    TEXT NOT NULL,
    last_name     TEXT NOT NULL,
    email         TEXT NOT NULL UNIQUE,
    password_hash BLOB NOT NULL,
    created_at    DATETIME NOT NULL
);

CREATE TABLE locations
(
    id          INTEGER PRIMARY KEY,
    lat         REAL NOT NULL,
    long        REAL NOT NULL,
    created_at  DATETIME NOT NULL
);
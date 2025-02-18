-- Enable citext extension if not already enabled
CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users (
                                     id          BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
                                     created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                     updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                     username    TEXT NOT NULL,
                                     email       CITEXT UNIQUE NOT NULL,
                                     address     TEXT,
                                     password    BYTEA NOT NULL,
                                     user_role   TEXT,
    status      TEXT,
    shop        TEXT,
    version     INT DEFAULT 1
    );
 CREATE TABLE IF NOT EXISTS users (
     id         bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY
     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
     username       TEXT NOT NULL,
     email      CITEXT UNIQUE NOT NULL,
     address    TEXT,
     password   Bytea NOT NULL,
     user_role   TEXT,
     status     TEXT,
     shop       TEXT,
     version    INT DEFAULT 1
     );
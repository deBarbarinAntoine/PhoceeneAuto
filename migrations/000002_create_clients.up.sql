CREATE TABLE IF NOT EXISTS clients (
    id           bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    created_at  TIMESTAMP(0) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP(0) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMP NULL,
    first_name  TEXT NOT NULL,
    last_name   TEXT NOT NULL,
    email       CITEXT UNIQUE NOT NULL,
    phone       TEXT,
    status      TEXT,
    address     TEXT,
    shop        TEXT,
    street TEXT,
    complement TEXT,
    city TEXT,
    zip_code TEXT,
    state TEXT,
    version     INT DEFAULT 1
);

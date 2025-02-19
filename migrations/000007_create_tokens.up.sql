CREATE TABLE tokens (
                       id SERIAL PRIMARY KEY,
                       hash TEXT NOT NULL UNIQUE,
                       scope TEXT NOT NULL,
                       expiry TIMESTAMP NOT NULL
);
CREATE INDEX idx_token_expiry ON tokens (expiry);

CREATE TABLE IF NOT EXISTS transactions (
    id           bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    client_id   INT,
    user_id     INT,
    status      TEXT,
    version     INT DEFAULT 1,
    lease_amount   FLOAT[],
    FOREIGN KEY (client_id) REFERENCES clients(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    );


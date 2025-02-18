 CREATE TABLE IF NOT EXISTS users (
     id         INT PRIMARY KEY AUTO_INCREMENT,
     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
     username       TEXT NOT NULL,
     email      CITEXT UNIQUE NOT NULL,
     address    TEXT,
     password   Bytea NOT NULL,
     userRole   TEXT,
     status     TEXT,
     shop       TEXT,
     version    INT DEFAULT 1
     );
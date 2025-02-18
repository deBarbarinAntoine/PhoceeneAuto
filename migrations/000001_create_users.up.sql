 CREATE TABLE IF NOT EXISTS users (
     id         INT PRIMARY KEY AUTO_INCREMENT,
     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
     username       TEXT NOT NULL,
     email      TEXT UNIQUE NOT NULL,
     address    CITEXT,
     password   Bytea NOT NULL,
     userRole   TEXT,
     status     TEXT,
     shop       TEXT,
     version    INT DEFAULT 1
     );
CREATE TABLE IF NOT EXISTS clients (
                                       id          INT PRIMARY KEY AUTO_INCREMENT,
                                       created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                       updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                                       deleted_at  TIMESTAMP NULL,
                                       first_name  TEXT NOT NULL,
                                       last_name   TEXT NOT NULL,
                                       email       CITEXT UNIQUE NOT NULL,
                                       phone       TEXT,
                                       status      TEXT,
                                       address     TEXT,
                                       shop        TEXT,
                                       version     INT DEFAULT 1
);

CREATE TABLE IF NOT EXISTS car_products (
    id          INT PRIMARY KEY AUTO_INCREMENT,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    status      TEXT,
    kilometers  FLOAT,
    owner_nb    INT,
    color       TEXT,
    price       FLOAT,
    shop        TEXT,
    version     INT DEFAULT 1,
    cat_id INT,
    FOREIGN KEY (cat_id) REFERENCES cars_catalog(id) ON DELETE SET NULL
    );

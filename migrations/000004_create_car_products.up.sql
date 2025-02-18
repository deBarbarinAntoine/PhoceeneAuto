CREATE TABLE IF NOT EXISTS car_products (
    id           bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    created_at  TIMESTAMP(0) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP(0) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    status      TEXT,
    kilometers  FLOAT,
    owner_nb    INT,
    color       TEXT,
    price       FLOAT,
    shop        TEXT,
    version     INT DEFAULT 1,
    cat_id INT,
    FOREIGN KEY (cat_id) REFERENCES cars_catalog(id)
    );

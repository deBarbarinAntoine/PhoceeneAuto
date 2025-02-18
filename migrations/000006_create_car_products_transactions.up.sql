CREATE TABLE IF NOT EXISTS car_products_transactions (
    transaction_id INT,
    car_product_id INT,
    PRIMARY KEY (transaction_id, car_product_id),
    FOREIGN KEY (transaction_id) REFERENCES transactions(id) ON DELETE CASCADE,
    FOREIGN KEY (car_product_id) REFERENCES car_products(id) ON DELETE CASCADE
);
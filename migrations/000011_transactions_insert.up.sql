-- Inserting three transactions into the transactions table
INSERT INTO transactions (
    client_id,
    user_id,
    status,
    lease_amount
) VALUES
(2, 1, 'PROCESSING', ARRAY[2500.00]),
(3, 2, 'ONGOING', ARRAY[1500.00, 1500.00]),
(4, 3, 'DONE', ARRAY[]::FLOAT[]);

INSERT INTO users (
    id,
    username,
    email,
    password_hash,
    user_role,
    status,
    shop,
    street,
    complement,
    city,
    zip_code,
    state
) VALUES
      (1, 'admin', 'admin@example.com', '$2a$12$mJcsAgeaac125wgl/I52oumgpEzVaoJJfWvFC5kB7IlfPNj61ET0i', 'ADMIN', 'ACTIVE', 'HEADQUARTERS', '789 Oak St', '', 'Springfield', '13579', 'IL'),
      (2, 'user1', 'user1@example.com', '$2a$12$mJcsAgeaac125wgl/I52oumgpEzVaoJJfWvFC5kB7IlfPNj61ET0i', 'USER', 'ACTIVE', 'HEADQUARTERS', '123 Main St', 'Apt 4B', 'Springfield', '12345', 'IL'),
      (3, 'user2', 'user2@example.com', '$2a$12$mJcsAgeaac125wgl/I52oumgpEzVaoJJfWvFC5kB7IlfPNj61ET0i', 'USER', 'ACTIVE', 'HEADQUARTERS', '456 Elm St', 'Suite 3C', 'Springfield', '67890', 'IL');

INSERT INTO roles (name, description)
VALUES
    ('ADMIN', 'Администратор системы'),
    ('USER', 'Стандартный пользователь')
ON CONFLICT (name) DO NOTHING;
INSERT INTO roles (name, description)
VALUES
    ('ADMIN', 'Администратор системы'),
    ('USER', 'Стандартный пользователь')
ON CONFLICT (name) DO NOTHING;

INSERT INTO categories (name)
VALUES
    ('web'),
    ('osint'),
    ('forensics'),
    ('crypto'),
    ('pwn')
ON CONFLICT (name) DO NOTHING;
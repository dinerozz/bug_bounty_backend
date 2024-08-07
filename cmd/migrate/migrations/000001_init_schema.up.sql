CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
                       id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                       username VARCHAR(255) UNIQUE NOT NULL,
                       email VARCHAR(255) UNIQUE NOT NULL,
                       password VARCHAR(255) NOT NULL,
                       points INTEGER
);

CREATE TABLE IF NOT EXISTS roles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) UNIQUE NOT NULL,
    description VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS user_roles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    role_id UUID REFERENCES roles(id),
    user_id UUID REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS teams (
                       id SERIAL PRIMARY KEY,
                       name VARCHAR(255) UNIQUE,
                       description VARCHAR(255),
                       owner_id UUID UNIQUE REFERENCES users(id),
                       invite_token VARCHAR(255) UNIQUE,
                       points INTEGER
);

CREATE TABLE IF NOT EXISTS team_members (
                              team_id INTEGER REFERENCES teams(id),
                              user_id UUID REFERENCES users(id),
                              CONSTRAINT pk_team_members PRIMARY KEY (team_id, user_id)
);

CREATE TABLE IF NOT EXISTS invites (
                         id SERIAL PRIMARY KEY,
                         team_id INTEGER REFERENCES teams(id),
                         token VARCHAR(255) UNIQUE NOT NULL,
                         expires_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) UNIQUE
);

CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    author_id UUID references users(id),
    category_id UUID references categories(id),
    points INTEGER,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    title VARCHAR(255),
    task_description TEXT
);

CREATE TABLE IF NOT EXISTS reports (
    id SERIAL PRIMARY KEY,
    author_id UUID references users(id),
    category_id UUID references categories(id),
    title VARCHAR(255),
    status VARCHAR(255) NOT NULL DEFAULT 'PENDING',
    description TEXT,
    team_id INTEGER references teams(id)
);

CREATE TABLE IF NOT EXISTS report_reviews (
    id SERIAL PRIMARY KEY,
    report_id int references reports(id),
    reviewer_id UUID references users(id),
    review_text TEXT,
    status VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS report_conversations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    report_id int references reports(id),
    user_id UUID references users(id),
    message TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);
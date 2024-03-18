CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       username VARCHAR(255) UNIQUE NOT NULL,
                       email VARCHAR(255) UNIQUE NOT NULL,
                       password VARCHAR(255) NOT NULL
);

CREATE TABLE teams (
                       id SERIAL PRIMARY KEY,
                       name VARCHAR(255) UNIQUE NOT NULL,
                       owner_id INTEGER REFERENCES users(id)
);

CREATE TABLE team_members (
                              team_id INTEGER REFERENCES teams(id),
                              user_id INTEGER REFERENCES users(id),
                              CONSTRAINT pk_team_members PRIMARY KEY (team_id, user_id)
);

CREATE TABLE invites (
                         id SERIAL PRIMARY KEY,
                         team_id INTEGER REFERENCES teams(id),
                         token VARCHAR(255) UNIQUE NOT NULL,
                         expires_at TIMESTAMP NOT NULL
);
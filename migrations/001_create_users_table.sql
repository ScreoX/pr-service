-- +goose Up
CREATE TABLE users
(
    id        TEXT PRIMARY KEY,
    username  VARCHAR(255) NOT NULL,
    team_name VARCHAR(255) NOT NULL,
    is_active BOOLEAN      NOT NULL DEFAULT TRUE
);

-- +goose Down
DROP TABLE IF EXISTS users;
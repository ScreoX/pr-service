-- +goose Up
CREATE TABLE teams
(
    id         TEXT PRIMARY KEY,
    team_name  VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS teams;
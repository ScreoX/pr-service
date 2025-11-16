-- +goose Up
CREATE TABLE pull_requests
(
    id                TEXT PRIMARY KEY,
    pull_request_name VARCHAR(255) NOT NULL,
    author_id         TEXT         NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    status            VARCHAR(50)  NOT NULL,
    created_at        TIMESTAMPTZ DEFAULT NOW(),
    merged_at         TIMESTAMPTZ
);

-- +goose Down
DROP TABLE IF EXISTS pull_requests;
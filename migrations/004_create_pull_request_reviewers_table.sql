-- +goose Up
CREATE TABLE pull_request_reviewers
(
    pull_request_id TEXT NOT NULL REFERENCES pull_requests (id) ON DELETE CASCADE,
    user_id         TEXT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    PRIMARY KEY (pull_request_id, user_id)
);

-- +goose Down
DROP TABLE IF EXISTS pull_request_reviewers;
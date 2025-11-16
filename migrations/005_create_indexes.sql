-- +goose Up
CREATE INDEX idx_pull_requests_status_author ON pull_requests (status, author_id);
CREATE INDEX idx_pull_request_reviewers_user_id ON pull_request_reviewers (user_id);

-- +goose Down
DROP INDEX IF EXISTS idx_pull_requests_status_author;
DROP INDEX IF EXISTS idx_pull_request_reviewers_user_id;
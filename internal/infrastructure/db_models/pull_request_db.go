package db_models

type PullRequest struct {
	ID        string  `db:"id"`
	Name      string  `db:"pull_request_name"`
	AuthorID  string  `db:"author_id"`
	Status    string  `db:"status"`
	CreatedAt string  `db:"created_at"`
	MergedAt  *string `db:"merged_at"`
}

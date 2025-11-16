package dto

type CreatePullRequest struct {
	PullRequestID   string `json:"pull_request_id" binding:"required"`
	PullRequestName string `json:"pull_request_name" binding:"required"`
	AuthorID        string `json:"author_id" binding:"required"`
}

type MergePullRequest struct {
	PullRequestID string `json:"pull_request_id" binding:"required"`
}

type ReassignReviewerRequest struct {
	PullRequestID string `json:"pull_request_id" binding:"required"`
	OldReviewerID string `json:"old_reviewer_id" binding:"required"`
}

type PullRequestResponse struct {
	PullRequestID     string   `json:"pull_request_id"`
	PullRequestName   string   `json:"pull_request_name"`
	AuthorID          string   `json:"author_id"`
	Status            string   `json:"status"`
	AssignedReviewers []string `json:"assigned_reviewers"`
	CreatedAt         string   `json:"created_at"`
	MergedAt          *string  `json:"merged_at,omitempty"`
}

type PullRequestReassignResponse struct {
	PullRequest PullRequestResponse `json:"pr"`
	ReplacedBy  string              `json:"replaced_by"`
}

type PullRequestShort struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
	Status          string `json:"status"`
}

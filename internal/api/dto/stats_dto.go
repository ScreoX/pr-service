package dto

type UserStats struct {
	UserID              string `json:"user_id"`
	Username            string `json:"username"`
	PullRequestsCreated int    `json:"prs_created"`
	TeamName            string `json:"team_name"`
}

type TeamStats struct {
	TeamName          string `json:"team_name"`
	MemberCount       int    `json:"member_count"`
	ActiveMembers     int    `json:"active_members"`
	PullRequestsCount int    `json:"prs_count"`
}

type ReviewAssignment struct {
	PullRequestID   string `json:"pr_id"`
	PullRequestName string `json:"pr_name"`
	AuthorID        string `json:"author_id"`
	Status          string `json:"status"`
}

type StatsResponse struct {
	TotalPullRequests  int                `json:"total_prs"`
	OpenPullRequests   int                `json:"open_prs"`
	MergedPullRequests int                `json:"merged_prs"`
	UsersStats         []UserStats        `json:"users_stats"`
	TeamsStats         []TeamStats        `json:"teams_stats"`
	ReviewAssignments  []ReviewAssignment `json:"review_assignments"`
}

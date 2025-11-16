package dto

type TeamMember struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type CreateTeamRequest struct {
	TeamName string       `json:"team_name" binding:"required"`
	Members  []TeamMember `json:"members" binding:"required,min=1"`
}

type TeamResponse struct {
	TeamName string       `json:"team_name"`
	Members  []TeamMember `json:"members"`
}

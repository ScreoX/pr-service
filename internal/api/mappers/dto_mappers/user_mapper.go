package dto_mappers

import (
	"pr-service/internal/api/dto"
	"pr-service/internal/domain/entities"
	"pr-service/internal/domain/value_objects"
)

func ToUserStatusResponseDTO(user entities.User) dto.UserStatusResponse {
	return dto.UserStatusResponse{
		UserID:   string(user.ID),
		Username: user.Username,
		TeamName: string(user.Team),
		IsActive: user.IsActive,
	}
}

func ToUserReviewsResponseDTO(userID value_objects.UserID, pullRequests []entities.PullRequest) dto.UserReviewsResponse {
	pullRequestShortDTOs := make([]dto.PullRequestShort, len(pullRequests))

	for i, pullRequest := range pullRequests {
		pullRequestShortDTOs[i] = ToPullRequestShortDTO(pullRequest)
	}

	return dto.UserReviewsResponse{
		UserID:       string(userID),
		PullRequests: pullRequestShortDTOs,
	}
}

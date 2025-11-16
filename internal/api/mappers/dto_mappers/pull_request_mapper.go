package dto_mappers

import (
	"pr-service/internal/api/dto"
	"pr-service/internal/domain/entities"
	"pr-service/internal/domain/value_objects"
)

const dateFormat = "2006-01-02T15:04:05Z"

func FromCreatePullRequestDTO(dto dto.CreatePullRequest) (value_objects.PullRequestID, string, value_objects.UserID) {
	pullRequestID := value_objects.PullRequestID(dto.PullRequestID)
	pullRequestName := dto.PullRequestName
	authorID := value_objects.UserID(dto.AuthorID)

	return pullRequestID, pullRequestName, authorID
}

func FromReassignReviewerRequestDTO(request dto.ReassignReviewerRequest) (value_objects.PullRequestID, value_objects.UserID) {
	pullRequestID := value_objects.PullRequestID(request.PullRequestID)
	oldReviewerID := value_objects.UserID(request.OldReviewerID)

	return pullRequestID, oldReviewerID
}

func ToPullRequestReassignResponseDTO(pullRequest entities.PullRequest, newReviewerID value_objects.UserID) dto.PullRequestReassignResponse {
	return dto.PullRequestReassignResponse{
		PullRequest: ToPullRequestResponseDTO(pullRequest),
		ReplacedBy:  string(newReviewerID),
	}
}

func ToPullRequestResponseDTO(pullRequest entities.PullRequest) dto.PullRequestResponse {
	var mergedAt *string
	if pullRequest.MergedAt != nil {
		mergedAtStr := pullRequest.MergedAt.Format(dateFormat)
		mergedAt = &mergedAtStr
	}

	return dto.PullRequestResponse{
		PullRequestID:     string(pullRequest.ID),
		PullRequestName:   pullRequest.Name,
		AuthorID:          string(pullRequest.AuthorID),
		Status:            string(pullRequest.Status),
		AssignedReviewers: toStringSlice(pullRequest.Reviewers()),
		CreatedAt:         pullRequest.CreatedAt.Format(dateFormat),
		MergedAt:          mergedAt,
	}
}

func ToPullRequestShortDTO(pullRequest entities.PullRequest) dto.PullRequestShort {
	return dto.PullRequestShort{
		PullRequestID:   string(pullRequest.ID),
		PullRequestName: pullRequest.Name,
		AuthorID:        string(pullRequest.AuthorID),
		Status:          string(pullRequest.Status),
	}
}

func toStringSlice(userIDs []value_objects.UserID) []string {
	var result []string

	for _, id := range userIDs {
		result = append(result, string(id))
	}

	return result
}

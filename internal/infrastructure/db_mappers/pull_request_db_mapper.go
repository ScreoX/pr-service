package db_mappers

import (
	"time"

	"pr-service/internal/domain/entities"
	"pr-service/internal/domain/value_objects"
	"pr-service/internal/infrastructure/db_models"
)

func ToPullRequestDBModel(pullRequest entities.PullRequest) db_models.PullRequest {
	var mergedAt *string
	if pullRequest.MergedAt != nil {
		s := pullRequest.MergedAt.Format(time.RFC3339)
		mergedAt = &s
	}

	return db_models.PullRequest{
		ID:        string(pullRequest.ID),
		Name:      pullRequest.Name,
		AuthorID:  string(pullRequest.AuthorID),
		Status:    string(pullRequest.Status),
		CreatedAt: pullRequest.CreatedAt.Format(time.RFC3339),
		MergedAt:  mergedAt,
	}
}

func FromPullRequestDBModel(dbPullRequest db_models.PullRequest) entities.PullRequest {
	createdAt, err := time.Parse(time.RFC3339, dbPullRequest.CreatedAt)
	if err != nil {
		return entities.PullRequest{}
	}

	var mergedAt *time.Time
	if dbPullRequest.MergedAt != nil {
		t, err := time.Parse(time.RFC3339, *dbPullRequest.MergedAt)
		if err != nil {
			return entities.PullRequest{}
		}

		mergedAt = &t
	}

	return entities.PullRequest{
		ID:        value_objects.PullRequestID(dbPullRequest.ID),
		Name:      dbPullRequest.Name,
		AuthorID:  value_objects.UserID(dbPullRequest.AuthorID),
		Status:    entities.PullRequestStatus(dbPullRequest.Status),
		CreatedAt: createdAt,
		MergedAt:  mergedAt,
	}
}

package app

import (
	"context"
	"time"

	"pr-service/internal/domain/entities"
	"pr-service/internal/domain/value_objects"
)

type TxManager interface {
	Do(ctx context.Context, operation func(ctx context.Context) error) error
}

type TimeProvider interface {
	Now() time.Time
}

type RandomProvider interface {
	Shuffle(n int, swapFunc func(i, j int))
	Intn(n int) int
}

type UserRepository interface {
	GetByID(ctx context.Context, id value_objects.UserID) (entities.User, error)
	GetUsersByTeam(ctx context.Context, teamName value_objects.TeamName) ([]entities.User, error)
	GetAll(ctx context.Context) ([]entities.User, error)
	UpsertMembers(ctx context.Context, teamName value_objects.TeamName, members []entities.User) error
	SetIsActive(ctx context.Context, id value_objects.UserID, isActive bool) (entities.User, error)
}

type TeamRepository interface {
	Create(ctx context.Context, team entities.Team) error
	GetByName(ctx context.Context, name value_objects.TeamName) (entities.Team, error)
	GetAll(ctx context.Context) ([]entities.Team, error)
}

type PullRequestRepository interface {
	Create(ctx context.Context, pullRequest *entities.PullRequest) error
	Save(ctx context.Context, pullRequest *entities.PullRequest) error
	GetByID(ctx context.Context, id value_objects.PullRequestID) (*entities.PullRequest, error)
	GetByReviewer(ctx context.Context, reviewerID value_objects.UserID) ([]entities.PullRequest, error)
	GetAll(ctx context.Context) ([]entities.PullRequest, error)
	ReassignReviewer(ctx context.Context, pullRequestID value_objects.PullRequestID, oldReviewerID value_objects.UserID, newReviewerID value_objects.UserID) error
}

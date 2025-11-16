package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"pr-service/internal/domain/entities"
	"pr-service/internal/domain/value_objects"
)

type UserRepository struct {
	mock.Mock
}

func (m *UserRepository) GetByID(ctx context.Context, id value_objects.UserID) (entities.User, error) {
	args := m.Called(ctx, id)

	return args.Get(0).(entities.User), args.Error(1)
}

func (m *UserRepository) GetUsersByTeam(ctx context.Context, teamName value_objects.TeamName) ([]entities.User, error) {
	args := m.Called(ctx, teamName)

	return args.Get(0).([]entities.User), args.Error(1)
}

func (m *UserRepository) GetAll(ctx context.Context) ([]entities.User, error) {
	args := m.Called(ctx)

	return args.Get(0).([]entities.User), args.Error(1)
}

func (m *UserRepository) UpsertMembers(ctx context.Context, teamName value_objects.TeamName, members []entities.User) error {
	args := m.Called(ctx, teamName, members)

	return args.Error(0)
}

func (m *UserRepository) SetIsActive(ctx context.Context, userID value_objects.UserID, isActive bool) (entities.User, error) {
	args := m.Called(ctx, userID, isActive)

	return args.Get(0).(entities.User), args.Error(1)
}

type TeamRepository struct {
	mock.Mock
}

func (m *TeamRepository) GetByName(ctx context.Context, name value_objects.TeamName) (entities.Team, error) {
	args := m.Called(ctx, name)

	return args.Get(0).(entities.Team), args.Error(1)
}

func (m *TeamRepository) Create(ctx context.Context, team entities.Team) error {
	args := m.Called(ctx, team)

	return args.Error(0)
}

func (m *TeamRepository) GetAll(ctx context.Context) ([]entities.Team, error) {
	args := m.Called(ctx)

	return args.Get(0).([]entities.Team), args.Error(1)
}

type PullRequestRepository struct {
	mock.Mock
}

func (m *PullRequestRepository) GetByID(ctx context.Context, id value_objects.PullRequestID) (*entities.PullRequest, error) {
	args := m.Called(ctx, id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entities.PullRequest), args.Error(1)
}

func (m *PullRequestRepository) Create(ctx context.Context, pullRequest *entities.PullRequest) error {
	args := m.Called(ctx, pullRequest)

	return args.Error(0)
}

func (m *PullRequestRepository) Save(ctx context.Context, pullRequest *entities.PullRequest) error {
	args := m.Called(ctx, pullRequest)

	return args.Error(0)
}

func (m *PullRequestRepository) GetByReviewer(ctx context.Context, reviewerID value_objects.UserID) ([]entities.PullRequest, error) {
	args := m.Called(ctx, reviewerID)

	return args.Get(0).([]entities.PullRequest), args.Error(1)
}

func (m *PullRequestRepository) GetAll(ctx context.Context) ([]entities.PullRequest, error) {
	args := m.Called(ctx)

	return args.Get(0).([]entities.PullRequest), args.Error(1)
}

func (m *PullRequestRepository) ReassignReviewer(ctx context.Context, pullRequestID value_objects.PullRequestID, oldReviewerID value_objects.UserID, newReviewerID value_objects.UserID) error {
	args := m.Called(ctx, pullRequestID, oldReviewerID, newReviewerID)

	return args.Error(0)
}

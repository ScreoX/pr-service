package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"pr-service/internal/app"
	"pr-service/internal/app/services/mocks"
	"pr-service/internal/domain"
	"pr-service/internal/domain/entities"
	"pr-service/internal/domain/value_objects"
)

func TestPullRequestService_Create(t *testing.T) {
	ctx := context.Background()
	fixedTime := time.Now()

	t.Run("successfully create pull request with reviewers", func(t *testing.T) {
		userRepository := &mocks.UserRepository{}
		teamRepository := &mocks.TeamRepository{}
		pullRequestRepository := &mocks.PullRequestRepository{}
		txManager := &mocks.TxManager{}
		timeProvider := &mocks.TimeProvider{}
		random := &mocks.RandomProvider{}

		pullRequestID := value_objects.PullRequestID("pull-request-1")
		pullRequestName := "Test Pull Request"
		authorID := value_objects.UserID("author1")
		author := entities.User{
			ID:       authorID,
			Username: "author",
			Team:     "backend",
			IsActive: true,
		}
		team := entities.Team{Name: "backend"}
		teamMembers := []entities.User{
			{ID: "user1", Username: "user1", Team: "backend", IsActive: true},
			{ID: "user2", Username: "user2", Team: "backend", IsActive: true},
			author,
		}

		pullRequestRepository.On("GetByID", ctx, pullRequestID).Return(nil, domain.ErrPRNotFound)
		userRepository.On("GetByID", ctx, authorID).Return(author, nil)
		teamRepository.On("GetByName", ctx, author.Team).Return(team, nil)
		userRepository.On("GetUsersByTeam", ctx, team.Name).Return(teamMembers, nil)
		timeProvider.On("Now").Return(fixedTime)
		random.On("Shuffle", mock.AnythingOfType("int"), mock.AnythingOfType("func(int, int)")).
			Run(func(args mock.Arguments) {
				swap := args.Get(1).(func(i, j int))
				swap(0, 1)
			})
		pullRequestRepository.On("Create", ctx, mock.AnythingOfType("*entities.PullRequest")).Return(nil)
		txManager.On("Do", ctx, mock.AnythingOfType("func(context.Context) error")).Return(nil)

		service := NewPullRequestService(userRepository, teamRepository, pullRequestRepository, txManager, timeProvider, random)
		result, err := service.Create(ctx, pullRequestID, pullRequestName, authorID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, pullRequestID, result.ID)
		assert.Equal(t, pullRequestName, result.Name)
		assert.Equal(t, authorID, result.AuthorID)
		assert.Len(t, result.Reviewers(), 2)
	})

	t.Run("fail when txManager is nil", func(t *testing.T) {
		userRepository := &mocks.UserRepository{}
		teamRepository := &mocks.TeamRepository{}
		pullRequestRepository := &mocks.PullRequestRepository{}
		timeProvider := &mocks.TimeProvider{}
		random := &mocks.RandomProvider{}

		service := NewPullRequestService(userRepository, teamRepository, pullRequestRepository, nil, timeProvider, random)
		result, err := service.Create(ctx, "pull-request-1", "Test Pull Request", "author1")

		assert.Error(t, err)
		assert.True(t, errors.Is(err, app.ErrTransactionRequired))
		assert.Nil(t, result)
	})

	t.Run("fail when pull request already exists", func(t *testing.T) {
		userRepository := &mocks.UserRepository{}
		teamRepository := &mocks.TeamRepository{}
		pullRequestRepository := &mocks.PullRequestRepository{}
		txManager := &mocks.TxManager{}
		timeProvider := &mocks.TimeProvider{}
		random := &mocks.RandomProvider{}

		pullRequestID := value_objects.PullRequestID("pull-request-1")
		existingPullRequest := &entities.PullRequest{ID: pullRequestID}

		pullRequestRepository.On("GetByID", ctx, pullRequestID).Return(existingPullRequest, nil)
		txManager.On("Do", ctx, mock.AnythingOfType("func(context.Context) error")).Return(domain.ErrPRExists)

		service := NewPullRequestService(userRepository, teamRepository, pullRequestRepository, txManager, timeProvider, random)
		result, err := service.Create(ctx, pullRequestID, "Test Pull Request", "author1")

		assert.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrPRExists))
		assert.Nil(t, result)
	})

	t.Run("fail when author not found", func(t *testing.T) {
		userRepository := &mocks.UserRepository{}
		teamRepository := &mocks.TeamRepository{}
		pullRequestRepository := &mocks.PullRequestRepository{}
		txManager := &mocks.TxManager{}
		timeProvider := &mocks.TimeProvider{}
		random := &mocks.RandomProvider{}

		pullRequestID := value_objects.PullRequestID("pull-request-1")
		authorID := value_objects.UserID("author1")

		pullRequestRepository.On("GetByID", ctx, pullRequestID).Return(nil, domain.ErrPRNotFound)
		userRepository.On("GetByID", ctx, authorID).Return(entities.User{}, domain.ErrUserNotFound)
		txManager.On("Do", ctx, mock.AnythingOfType("func(context.Context) error")).Return(domain.ErrUserNotFound)

		service := NewPullRequestService(userRepository, teamRepository, pullRequestRepository, txManager, timeProvider, random)
		result, err := service.Create(ctx, pullRequestID, "Test Pull Request", authorID)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrUserNotFound))
		assert.Nil(t, result)
	})

	t.Run("fail when no active candidates", func(t *testing.T) {
		userRepository := &mocks.UserRepository{}
		teamRepository := &mocks.TeamRepository{}
		pullRequestRepository := &mocks.PullRequestRepository{}
		txManager := &mocks.TxManager{}
		timeProvider := &mocks.TimeProvider{}
		random := &mocks.RandomProvider{}

		pullRequestID := value_objects.PullRequestID("pull-request-1")
		authorID := value_objects.UserID("author1")
		author := entities.User{
			ID:       authorID,
			Username: "author",
			Team:     "backend",
			IsActive: true,
		}
		team := entities.Team{Name: "backend"}
		teamMembers := []entities.User{author}

		pullRequestRepository.On("GetByID", ctx, pullRequestID).Return(nil, domain.ErrPRNotFound)
		userRepository.On("GetByID", ctx, authorID).Return(author, nil)
		teamRepository.On("GetByName", ctx, author.Team).Return(team, nil)
		userRepository.On("GetUsersByTeam", ctx, team.Name).Return(teamMembers, nil)
		timeProvider.On("Now").Return(fixedTime)
		txManager.On("Do", ctx, mock.AnythingOfType("func(context.Context) error")).Return(domain.ErrNoCandidate)

		service := NewPullRequestService(userRepository, teamRepository, pullRequestRepository, txManager, timeProvider, random)
		result, err := service.Create(ctx, pullRequestID, "Test Pull Request", authorID)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrNoCandidate))
		assert.Nil(t, result)
	})
}

func TestPullRequestService_Merge(t *testing.T) {
	ctx := context.Background()
	fixedTime := time.Now()

	t.Run("successfully merge pull request", func(t *testing.T) {
		userRepository := &mocks.UserRepository{}
		teamRepository := &mocks.TeamRepository{}
		pullRequestRepository := &mocks.PullRequestRepository{}
		txManager := &mocks.TxManager{}
		timeProvider := &mocks.TimeProvider{}
		random := &mocks.RandomProvider{}

		pullRequestID := value_objects.PullRequestID("pull-request-1")
		pullRequest := &entities.PullRequest{
			ID:     pullRequestID,
			Name:   "Test Pull Request",
			Status: entities.StatusOpen,
		}

		pullRequestRepository.On("GetByID", ctx, pullRequestID).Return(pullRequest, nil)
		timeProvider.On("Now").Return(fixedTime)
		pullRequestRepository.On("Save", ctx, pullRequest).Return(nil)

		service := NewPullRequestService(userRepository, teamRepository, pullRequestRepository, txManager, timeProvider, random)
		result, err := service.Merge(ctx, pullRequestID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, entities.StatusMerged, result.Status)
	})

	t.Run("fail when pull request not found", func(t *testing.T) {
		userRepository := &mocks.UserRepository{}
		teamRepository := &mocks.TeamRepository{}
		pullRequestRepository := &mocks.PullRequestRepository{}
		txManager := &mocks.TxManager{}
		timeProvider := &mocks.TimeProvider{}
		random := &mocks.RandomProvider{}

		pullRequestID := value_objects.PullRequestID("pull-request-1")

		pullRequestRepository.On("GetByID", ctx, pullRequestID).Return(nil, domain.ErrPRNotFound)

		service := NewPullRequestService(userRepository, teamRepository, pullRequestRepository, txManager, timeProvider, random)
		result, err := service.Merge(ctx, pullRequestID)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrPRNotFound))
		assert.Nil(t, result)
	})

	t.Run("fail when save fails", func(t *testing.T) {
		userRepository := &mocks.UserRepository{}
		teamRepository := &mocks.TeamRepository{}
		pullRequestRepository := &mocks.PullRequestRepository{}
		txManager := &mocks.TxManager{}
		timeProvider := &mocks.TimeProvider{}
		random := &mocks.RandomProvider{}

		pullRequestID := value_objects.PullRequestID("pull-request-1")
		pullRequest := &entities.PullRequest{
			ID:     pullRequestID,
			Name:   "Test Pull Request",
			Status: entities.StatusOpen,
		}

		pullRequestRepository.On("GetByID", ctx, pullRequestID).Return(pullRequest, nil)
		timeProvider.On("Now").Return(fixedTime)
		pullRequestRepository.On("Save", ctx, pullRequest).Return(errors.New("save error"))

		service := NewPullRequestService(userRepository, teamRepository, pullRequestRepository, txManager, timeProvider, random)
		result, err := service.Merge(ctx, pullRequestID)

		assert.Error(t, err)
		assert.Equal(t, "save error", err.Error())
		assert.Nil(t, result)
	})
}

func TestPullRequestService_ReassignReviewer(t *testing.T) {
	ctx := context.Background()

	t.Run("successfully reassign reviewer", func(t *testing.T) {
		userRepository := &mocks.UserRepository{}
		teamRepository := &mocks.TeamRepository{}
		pullRequestRepository := &mocks.PullRequestRepository{}
		txManager := &mocks.TxManager{}
		timeProvider := &mocks.TimeProvider{}
		random := &mocks.RandomProvider{}

		pullRequestID := value_objects.PullRequestID("pull-request-1")
		oldReviewerID := value_objects.UserID("reviewer1")
		newReviewerID := value_objects.UserID("reviewer2")
		authorID := value_objects.UserID("author1")

		initialPullRequest := &entities.PullRequest{
			ID:       pullRequestID,
			Name:     "Test Pull Request",
			AuthorID: authorID,
			Status:   entities.StatusOpen,
		}
		initialPullRequest.AddReviewers([]value_objects.UserID{oldReviewerID})

		updatedPullRequest := &entities.PullRequest{
			ID:       pullRequestID,
			Name:     "Test Pull Request",
			AuthorID: authorID,
			Status:   entities.StatusOpen,
		}
		updatedPullRequest.AddReviewers([]value_objects.UserID{newReviewerID})

		author := entities.User{
			ID:       authorID,
			Username: "author",
			Team:     "backend",
			IsActive: true,
		}

		oldReviewer := entities.User{
			ID:       oldReviewerID,
			Username: "reviewer1",
			Team:     "backend",
			IsActive: true,
		}

		team := entities.Team{Name: "backend"}
		teamMembers := []entities.User{
			{ID: newReviewerID, Username: "reviewer2", Team: "backend", IsActive: true},
			{ID: "user3", Username: "user3", Team: "backend", IsActive: true},
		}

		pullRequestRepository.On("GetByID", ctx, pullRequestID).Return(initialPullRequest, nil).Once()

		userRepository.On("GetByID", ctx, oldReviewerID).Return(oldReviewer, nil)
		userRepository.On("GetByID", ctx, authorID).Return(author, nil)
		teamRepository.On("GetByName", ctx, author.Team).Return(team, nil)
		userRepository.On("GetUsersByTeam", ctx, team.Name).Return(teamMembers, nil)
		random.On("Intn", 2).Return(0)

		pullRequestRepository.On("ReassignReviewer", ctx, pullRequestID, oldReviewerID, newReviewerID).Return(nil)

		pullRequestRepository.On("GetByID", ctx, pullRequestID).Return(updatedPullRequest, nil).Once()

		txManager.On("Do", ctx, mock.AnythingOfType("func(context.Context) error")).Return(nil)

		service := NewPullRequestService(userRepository, teamRepository, pullRequestRepository, txManager, timeProvider, random)
		resultPullRequest, resultReviewer, err := service.ReassignReviewer(ctx, pullRequestID, oldReviewerID)

		assert.NoError(t, err)
		assert.NotNil(t, resultPullRequest)
		assert.Equal(t, newReviewerID, resultReviewer)
		assert.Contains(t, resultPullRequest.Reviewers(), newReviewerID)
		assert.NotContains(t, resultPullRequest.Reviewers(), oldReviewerID)
	})

	t.Run("fail when txManager is nil", func(t *testing.T) {
		userRepository := &mocks.UserRepository{}
		teamRepository := &mocks.TeamRepository{}
		pullRequestRepository := &mocks.PullRequestRepository{}
		timeProvider := &mocks.TimeProvider{}
		random := &mocks.RandomProvider{}

		service := NewPullRequestService(userRepository, teamRepository, pullRequestRepository, nil, timeProvider, random)
		resultPullRequest, resultReviewer, err := service.ReassignReviewer(ctx, "pull-request-1", "reviewer1")

		assert.Error(t, err)
		assert.True(t, errors.Is(err, app.ErrTransactionRequired))
		assert.Nil(t, resultPullRequest)
		assert.Equal(t, value_objects.UserID(""), resultReviewer)
	})

	t.Run("fail when pull request not found", func(t *testing.T) {
		userRepository := &mocks.UserRepository{}
		teamRepository := &mocks.TeamRepository{}
		pullRequestRepository := &mocks.PullRequestRepository{}
		txManager := &mocks.TxManager{}
		timeProvider := &mocks.TimeProvider{}
		random := &mocks.RandomProvider{}

		pullRequestID := value_objects.PullRequestID("pull-request-1")

		pullRequestRepository.On("GetByID", ctx, pullRequestID).Return(nil, domain.ErrPRNotFound)
		txManager.On("Do", ctx, mock.AnythingOfType("func(context.Context) error")).Return(domain.ErrPRNotFound)

		service := NewPullRequestService(userRepository, teamRepository, pullRequestRepository, txManager, timeProvider, random)
		resultPullRequest, resultReviewer, err := service.ReassignReviewer(ctx, pullRequestID, "reviewer1")

		assert.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrPRNotFound))
		assert.Nil(t, resultPullRequest)
		assert.Equal(t, value_objects.UserID(""), resultReviewer)
	})

	t.Run("fail when pull request is merged", func(t *testing.T) {
		userRepository := &mocks.UserRepository{}
		teamRepository := &mocks.TeamRepository{}
		pullRequestRepository := &mocks.PullRequestRepository{}
		txManager := &mocks.TxManager{}
		timeProvider := &mocks.TimeProvider{}
		random := &mocks.RandomProvider{}

		pullRequestID := value_objects.PullRequestID("pull-request-1")
		pullRequest := &entities.PullRequest{
			ID:     pullRequestID,
			Name:   "Test Pull Request",
			Status: entities.StatusMerged,
		}

		pullRequestRepository.On("GetByID", ctx, pullRequestID).Return(pullRequest, nil)
		txManager.On("Do", ctx, mock.AnythingOfType("func(context.Context) error")).Return(domain.ErrPRMerged)

		service := NewPullRequestService(userRepository, teamRepository, pullRequestRepository, txManager, timeProvider, random)
		resultPullRequest, resultReviewer, err := service.ReassignReviewer(ctx, pullRequestID, "reviewer1")

		assert.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrPRMerged))
		assert.Nil(t, resultPullRequest)
		assert.Equal(t, value_objects.UserID(""), resultReviewer)
	})

	t.Run("fail when reviewer not assigned", func(t *testing.T) {
		userRepository := &mocks.UserRepository{}
		teamRepository := &mocks.TeamRepository{}
		pullRequestRepository := &mocks.PullRequestRepository{}
		txManager := &mocks.TxManager{}
		timeProvider := &mocks.TimeProvider{}
		random := &mocks.RandomProvider{}

		pullRequestID := value_objects.PullRequestID("pull-request-1")
		pullRequest := &entities.PullRequest{
			ID:     pullRequestID,
			Name:   "Test Pull Request",
			Status: entities.StatusOpen,
		}
		pullRequest.AddReviewers([]value_objects.UserID{"other_reviewer"})

		pullRequestRepository.On("GetByID", ctx, pullRequestID).Return(pullRequest, nil)
		txManager.On("Do", ctx, mock.AnythingOfType("func(context.Context) error")).Return(domain.ErrNotAssigned)

		service := NewPullRequestService(userRepository, teamRepository, pullRequestRepository, txManager, timeProvider, random)
		resultPullRequest, resultReviewer, err := service.ReassignReviewer(ctx, pullRequestID, "reviewer1")

		assert.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrNotAssigned))
		assert.Nil(t, resultPullRequest)
		assert.Equal(t, value_objects.UserID(""), resultReviewer)
	})

	t.Run("fail when no candidates available", func(t *testing.T) {
		userRepository := &mocks.UserRepository{}
		teamRepository := &mocks.TeamRepository{}
		pullRequestRepository := &mocks.PullRequestRepository{}
		txManager := &mocks.TxManager{}
		timeProvider := &mocks.TimeProvider{}
		random := &mocks.RandomProvider{}

		pullRequestID := value_objects.PullRequestID("pull-request-1")
		oldReviewerID := value_objects.UserID("reviewer1")
		authorID := value_objects.UserID("author1")

		pullRequest := &entities.PullRequest{
			ID:       pullRequestID,
			Name:     "Test Pull Request",
			AuthorID: authorID,
			Status:   entities.StatusOpen,
		}
		pullRequest.AddReviewers([]value_objects.UserID{oldReviewerID})

		author := entities.User{
			ID:       authorID,
			Username: "author",
			Team:     "backend",
			IsActive: true,
		}
		team := entities.Team{Name: "backend"}
		teamMembers := []entities.User{author}

		pullRequestRepository.On("GetByID", ctx, pullRequestID).Return(pullRequest, nil)
		userRepository.On("GetByID", ctx, authorID).Return(author, nil)
		teamRepository.On("GetByName", ctx, author.Team).Return(team, nil)
		userRepository.On("GetUsersByTeam", ctx, team.Name).Return(teamMembers, nil)
		txManager.On("Do", ctx, mock.AnythingOfType("func(context.Context) error")).Return(domain.ErrNoCandidate)

		service := NewPullRequestService(userRepository, teamRepository, pullRequestRepository, txManager, timeProvider, random)
		resultPullRequest, resultReviewer, err := service.ReassignReviewer(ctx, pullRequestID, oldReviewerID)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrNoCandidate))
		assert.Nil(t, resultPullRequest)
		assert.Equal(t, value_objects.UserID(""), resultReviewer)
	})
}

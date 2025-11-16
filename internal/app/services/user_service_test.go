package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"pr-service/internal/app/services/mocks"
	"pr-service/internal/domain"
	"pr-service/internal/domain/entities"
	"pr-service/internal/domain/value_objects"
)

func TestUserService_SetActiveStatus(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		userID        value_objects.UserID
		isActive      bool
		setupMocks    func(userRepository *mocks.UserRepository, pullRequestRepository *mocks.PullRequestRepository)
		expectedUser  entities.User
		expectedError error
	}{
		{
			name:     "successfully activate user",
			userID:   value_objects.UserID("user1"),
			isActive: true,
			setupMocks: func(userRepository *mocks.UserRepository, pullRequestRepository *mocks.PullRequestRepository) {
				expectedUser := entities.User{
					ID:       value_objects.UserID("user1"),
					Username: "Alice",
					Team:     value_objects.TeamName("backend"),
					IsActive: true,
				}
				userRepository.On("SetIsActive", ctx, value_objects.UserID("user1"), true).Return(expectedUser, nil)
			},
			expectedUser: entities.User{
				ID:       value_objects.UserID("user1"),
				Username: "Alice",
				Team:     value_objects.TeamName("backend"),
				IsActive: true,
			},
			expectedError: nil,
		},
		{
			name:     "successfully deactivate user",
			userID:   value_objects.UserID("user2"),
			isActive: false,
			setupMocks: func(userRepository *mocks.UserRepository, pullRequestRepository *mocks.PullRequestRepository) {
				expectedUser := entities.User{
					ID:       value_objects.UserID("user2"),
					Username: "Bob",
					Team:     value_objects.TeamName("backend"),
					IsActive: false,
				}
				userRepository.On("SetIsActive", ctx, value_objects.UserID("user2"), false).Return(expectedUser, nil)
			},
			expectedUser: entities.User{
				ID:       value_objects.UserID("user2"),
				Username: "Bob",
				Team:     value_objects.TeamName("backend"),
				IsActive: false,
			},
			expectedError: nil,
		},
		{
			name:     "fail when user not found",
			userID:   value_objects.UserID("nonexistent"),
			isActive: true,
			setupMocks: func(userRepository *mocks.UserRepository, pullRequestRepository *mocks.PullRequestRepository) {
				userRepository.On("SetIsActive", ctx, value_objects.UserID("nonexistent"), true).Return(entities.User{}, domain.ErrUserNotFound)
			},
			expectedUser:  entities.User{},
			expectedError: domain.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepository := &mocks.UserRepository{}
			pullRequestRepository := &mocks.PullRequestRepository{}
			tt.setupMocks(userRepository, pullRequestRepository)

			service := NewUserService(userRepository, pullRequestRepository)

			resultUser, err := service.SetActiveStatus(ctx, tt.userID, tt.isActive)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, domain.ErrUserNotFound))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser, resultUser)
			}

			userRepository.AssertExpectations(t)
		})
	}
}

func TestUserService_GetUserReviews(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	tests := []struct {
		name                 string
		userID               value_objects.UserID
		setupMocks           func(userRepository *mocks.UserRepository, pullRequestRepository *mocks.PullRequestRepository)
		expectedPullRequests []entities.PullRequest
		expectedError        error
	}{
		{
			name:   "successfully get user reviews with multiple pull requests",
			userID: value_objects.UserID("user1"),
			setupMocks: func(userRepository *mocks.UserRepository, pullRequestRepository *mocks.PullRequestRepository) {
				user := entities.User{
					ID:       value_objects.UserID("user1"),
					Username: "Alice",
					Team:     value_objects.TeamName("backend"),
					IsActive: true,
				}
				userRepository.On("GetByID", ctx, value_objects.UserID("user1")).Return(user, nil)

				pullRequest1 := entities.NewPullRequest(
					"pullRequest1",
					"Feature A",
					"user2",
					now.Add(-24*time.Hour),
				)
				pullRequest1.AddReviewers([]value_objects.UserID{
					"user1",
					"user3",
				})

				pullRequest2 := entities.NewPullRequest(
					"pullRequest2",
					"Feature B",
					"user3",
					now.Add(-12*time.Hour),
				)
				pullRequest2.AddReviewers([]value_objects.UserID{"user1"})
				pullRequest2.Merge(now)

				pullRequests := []entities.PullRequest{*pullRequest1, *pullRequest2}
				pullRequestRepository.On("GetByReviewer", ctx, value_objects.UserID("user1")).Return(pullRequests, nil)
			},
			expectedPullRequests: func() []entities.PullRequest {
				pullRequest1 := entities.NewPullRequest(
					"pullRequest1",
					"Feature A",
					"user2",
					now.Add(-24*time.Hour),
				)
				pullRequest1.AddReviewers([]value_objects.UserID{
					"user1",
					"user3",
				})

				pullRequest2 := entities.NewPullRequest(
					"pullRequest2",
					"Feature B",
					"user3",
					now.Add(-12*time.Hour),
				)
				pullRequest2.AddReviewers([]value_objects.UserID{"user1"})
				pullRequest2.Merge(now)

				return []entities.PullRequest{*pullRequest1, *pullRequest2}
			}(),
			expectedError: nil,
		},
		{
			name:   "return empty list when user has no reviews",
			userID: value_objects.UserID("user2"),
			setupMocks: func(userRepository *mocks.UserRepository, pullRequestRepository *mocks.PullRequestRepository) {
				user := entities.User{
					ID:       value_objects.UserID("user2"),
					Username: "Bob",
					Team:     value_objects.TeamName("backend"),
					IsActive: true,
				}
				userRepository.On("GetByID", ctx, value_objects.UserID("user2")).Return(user, nil)
				pullRequestRepository.On("GetByReviewer", ctx, value_objects.UserID("user2")).Return([]entities.PullRequest{}, nil)
			},
			expectedPullRequests: []entities.PullRequest{},
			expectedError:        nil,
		},
		{
			name:   "fail when user not found",
			userID: value_objects.UserID("nonexistent"),
			setupMocks: func(userRepository *mocks.UserRepository, pullRequestRepository *mocks.PullRequestRepository) {
				userRepository.On("GetByID", ctx, value_objects.UserID("nonexistent")).Return(entities.User{}, domain.ErrUserNotFound)
			},
			expectedPullRequests: nil,
			expectedError:        domain.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepository := &mocks.UserRepository{}
			pullRequestRepository := &mocks.PullRequestRepository{}
			tt.setupMocks(userRepository, pullRequestRepository)

			service := NewUserService(userRepository, pullRequestRepository)

			resultPullRequests, err := service.GetUserReviews(ctx, tt.userID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, domain.ErrUserNotFound))
				assert.Nil(t, resultPullRequests)
			} else {
				assert.NoError(t, err)
				assert.Len(t, resultPullRequests, len(tt.expectedPullRequests))

				for i, expectedPullRequest := range tt.expectedPullRequests {
					assert.Equal(t, expectedPullRequest.ID, resultPullRequests[i].ID)
					assert.Equal(t, expectedPullRequest.Name, resultPullRequests[i].Name)
					assert.Equal(t, expectedPullRequest.AuthorID, resultPullRequests[i].AuthorID)
					assert.Equal(t, expectedPullRequest.Status, resultPullRequests[i].Status)
					assert.ElementsMatch(t, expectedPullRequest.Reviewers(), resultPullRequests[i].Reviewers())
				}
			}

			userRepository.AssertExpectations(t)
			pullRequestRepository.AssertExpectations(t)
		})
	}
}

func TestUserService_GetUserReviews_DoesNotCallPullRequestRepositoryWhenUserNotFound(t *testing.T) {
	ctx := context.Background()
	userRepository := &mocks.UserRepository{}
	pullRequestRepository := &mocks.PullRequestRepository{}

	userRepository.On("GetByID", ctx, value_objects.UserID("nonexistent")).Return(entities.User{}, domain.ErrUserNotFound)

	service := NewUserService(userRepository, pullRequestRepository)

	resultPullRequests, err := service.GetUserReviews(ctx, "nonexistent")

	assert.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrUserNotFound))
	assert.Nil(t, resultPullRequests)

	pullRequestRepository.AssertNotCalled(t, "GetByReviewer", mock.Anything, mock.Anything)
	userRepository.AssertExpectations(t)
}

package services

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"pr-service/internal/api/dto"
	"pr-service/internal/app/services/mocks"
	"pr-service/internal/domain/entities"
)

func TestStatsService_GetStats(t *testing.T) {
	ctx := context.Background()

	t.Run("successfully get statistics with complete data", func(t *testing.T) {
		userRepository := &mocks.UserRepository{}
		teamRepository := &mocks.TeamRepository{}
		pullRequestRepository := &mocks.PullRequestRepository{}

		users := []entities.User{
			{
				ID:       "user1",
				Username: "Alice",
				Team:     "backend",
				IsActive: true,
			},
			{
				ID:       "user2",
				Username: "Bob",
				Team:     "backend",
				IsActive: true,
			},
			{
				ID:       "user3",
				Username: "Charlie",
				Team:     "frontend",
				IsActive: true,
			},
			{
				ID:       "user4",
				Username: "David",
				Team:     "frontend",
				IsActive: false,
			},
		}

		teams := []entities.Team{
			{Name: "backend"},
			{Name: "frontend"},
		}

		pullRequests := []entities.PullRequest{
			{
				ID:       "pull-request-1",
				Name:     "Add feature A",
				AuthorID: "user1",
				Status:   entities.StatusOpen,
			},
			{
				ID:       "pull-request-2",
				Name:     "Fix bug B",
				AuthorID: "user1",
				Status:   entities.StatusOpen,
			},
			{
				ID:       "pull-request-3",
				Name:     "Refactor C",
				AuthorID: "user3",
				Status:   entities.StatusMerged,
			},
		}

		userRepository.On("GetAll", ctx).Return(users, nil)
		teamRepository.On("GetAll", ctx).Return(teams, nil)
		pullRequestRepository.On("GetAll", ctx).Return(pullRequests, nil)

		service := NewStatsService(userRepository, teamRepository, pullRequestRepository)
		stats, err := service.GetStats(ctx)

		assert.NoError(t, err)
		assert.NotNil(t, stats)

		assert.Equal(t, 3, stats.TotalPullRequests)
		assert.Equal(t, 2, stats.OpenPullRequests)
		assert.Equal(t, 1, stats.MergedPullRequests)

		assert.Len(t, stats.UsersStats, 4)

		aliceStats := findUserStats(stats.UsersStats, "user1")
		assert.NotNil(t, aliceStats)
		assert.Equal(t, "Alice", aliceStats.Username)
		assert.Equal(t, "backend", aliceStats.TeamName)
		assert.Equal(t, 2, aliceStats.PullRequestsCreated)

		charlieStats := findUserStats(stats.UsersStats, "user3")
		assert.NotNil(t, charlieStats)
		assert.Equal(t, "Charlie", charlieStats.Username)
		assert.Equal(t, "frontend", charlieStats.TeamName)
		assert.Equal(t, 1, charlieStats.PullRequestsCreated)

		assert.Len(t, stats.TeamsStats, 2)

		backendStats := findTeamStats(stats.TeamsStats, "backend")
		assert.NotNil(t, backendStats)
		assert.Equal(t, 2, backendStats.MemberCount)
		assert.Equal(t, 2, backendStats.ActiveMembers)
		assert.Equal(t, 2, backendStats.PullRequestsCount)

		frontendStats := findTeamStats(stats.TeamsStats, "frontend")
		assert.NotNil(t, frontendStats)
		assert.Equal(t, 2, frontendStats.MemberCount)
		assert.Equal(t, 1, frontendStats.ActiveMembers)
		assert.Equal(t, 1, frontendStats.PullRequestsCount)

		assert.Len(t, stats.ReviewAssignments, 3)

		userRepository.AssertCalled(t, "GetAll", ctx)
		teamRepository.AssertCalled(t, "GetAll", ctx)
		pullRequestRepository.AssertCalled(t, "GetAll", ctx)
	})

	t.Run("return empty statistics when no data", func(t *testing.T) {
		userRepository := &mocks.UserRepository{}
		teamRepository := &mocks.TeamRepository{}
		pullRequestRepository := &mocks.PullRequestRepository{}

		userRepository.On("GetAll", ctx).Return([]entities.User{}, nil)
		teamRepository.On("GetAll", ctx).Return([]entities.Team{}, nil)
		pullRequestRepository.On("GetAll", ctx).Return([]entities.PullRequest{}, nil)

		service := NewStatsService(userRepository, teamRepository, pullRequestRepository)
		stats, err := service.GetStats(ctx)

		assert.NoError(t, err)
		assert.NotNil(t, stats)
		assert.Equal(t, 0, stats.TotalPullRequests)
		assert.Equal(t, 0, stats.OpenPullRequests)
		assert.Equal(t, 0, stats.MergedPullRequests)
		assert.Empty(t, stats.UsersStats)
		assert.Empty(t, stats.TeamsStats)
		assert.Empty(t, stats.ReviewAssignments)
	})

	t.Run("fail when user repository returns error", func(t *testing.T) {
		userRepository := &mocks.UserRepository{}
		teamRepository := &mocks.TeamRepository{}
		pullRequestRepository := &mocks.PullRequestRepository{}

		userRepository.On("GetAll", ctx).Return([]entities.User{}, errors.New("database error"))

		service := NewStatsService(userRepository, teamRepository, pullRequestRepository)
		stats, err := service.GetStats(ctx)

		assert.Error(t, err)
		assert.Nil(t, stats)
		assert.Equal(t, "database error", err.Error())

		userRepository.AssertCalled(t, "GetAll", ctx)
		teamRepository.AssertNotCalled(t, "GetAll", ctx)
		pullRequestRepository.AssertNotCalled(t, "GetAll", ctx)
	})

	t.Run("fail when team repository returns error", func(t *testing.T) {
		userRepository := &mocks.UserRepository{}
		teamRepository := &mocks.TeamRepository{}
		pullRequestRepository := &mocks.PullRequestRepository{}

		users := []entities.User{
			{ID: "user1", Username: "Test", Team: "backend", IsActive: true},
		}

		userRepository.On("GetAll", ctx).Return(users, nil)
		teamRepository.On("GetAll", ctx).Return([]entities.Team{}, errors.New("team database error"))

		service := NewStatsService(userRepository, teamRepository, pullRequestRepository)
		stats, err := service.GetStats(ctx)

		assert.Error(t, err)
		assert.Nil(t, stats)
		assert.Equal(t, "team database error", err.Error())

		userRepository.AssertCalled(t, "GetAll", ctx)
		teamRepository.AssertCalled(t, "GetAll", ctx)
		pullRequestRepository.AssertNotCalled(t, "GetAll", ctx)
	})

	t.Run("fail when pull request repository returns error", func(t *testing.T) {
		userRepository := &mocks.UserRepository{}
		teamRepository := &mocks.TeamRepository{}
		pullRequestRepository := &mocks.PullRequestRepository{}

		users := []entities.User{
			{ID: "user1", Username: "Test", Team: "backend", IsActive: true},
		}
		teams := []entities.Team{
			{Name: "backend"},
		}

		userRepository.On("GetAll", ctx).Return(users, nil)
		teamRepository.On("GetAll", ctx).Return(teams, nil)
		pullRequestRepository.On("GetAll", ctx).Return([]entities.PullRequest{}, errors.New("pr database error"))

		service := NewStatsService(userRepository, teamRepository, pullRequestRepository)
		stats, err := service.GetStats(ctx)

		assert.Error(t, err)
		assert.Nil(t, stats)
		assert.Equal(t, "pr database error", err.Error())

		userRepository.AssertCalled(t, "GetAll", ctx)
		teamRepository.AssertCalled(t, "GetAll", ctx)
		pullRequestRepository.AssertCalled(t, "GetAll", ctx)
	})
}

func findUserStats(userStats []dto.UserStats, userID string) *dto.UserStats {
	for _, stats := range userStats {
		if stats.UserID == userID {
			return &stats
		}
	}
	return nil
}

func findTeamStats(teamStats []dto.TeamStats, teamName string) *dto.TeamStats {
	for _, stats := range teamStats {
		if stats.TeamName == teamName {
			return &stats
		}
	}
	return nil
}

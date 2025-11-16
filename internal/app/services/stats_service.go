package services

import (
	"context"
	"pr-service/internal/api/dto"
	"pr-service/internal/app"
	"pr-service/internal/domain/entities"
)

type StatsService interface {
	GetStats(ctx context.Context) (*dto.StatsResponse, error)
}

type statsService struct {
	userRepository        app.UserRepository
	teamRepository        app.TeamRepository
	pullRequestRepository app.PullRequestRepository
}

func NewStatsService(userRepository app.UserRepository, teamRepository app.TeamRepository, pullRequestRepository app.PullRequestRepository) StatsService {
	return &statsService{
		userRepository:        userRepository,
		teamRepository:        teamRepository,
		pullRequestRepository: pullRequestRepository,
	}
}

func (s *statsService) GetStats(ctx context.Context) (*dto.StatsResponse, error) {
	allUsers, err := s.userRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	allTeams, err := s.teamRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	allPullRequests, err := s.pullRequestRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	stats := &dto.StatsResponse{
		TotalPullRequests:  len(allPullRequests),
		OpenPullRequests:   countOpenPullRequests(allPullRequests),
		MergedPullRequests: countMergedPullRequests(allPullRequests),
		UsersStats:         calculateUserStats(allUsers, allPullRequests),
		TeamsStats:         calculateTeamStats(allTeams, allUsers, allPullRequests),
		ReviewAssignments:  calculateReviewAssignments(allPullRequests),
	}

	return stats, nil
}

func countOpenPullRequests(pullRequests []entities.PullRequest) int {
	count := 0

	for _, pullRequest := range pullRequests {
		if pullRequest.Status == entities.StatusOpen {
			count++
		}
	}

	return count
}

func countMergedPullRequests(pullRequests []entities.PullRequest) int {
	count := 0

	for _, pullRequest := range pullRequests {
		if pullRequest.Status == entities.StatusMerged {
			count++
		}
	}

	return count
}

func calculateUserStats(users []entities.User, pullRequests []entities.PullRequest) []dto.UserStats {
	var userStats []dto.UserStats

	for _, user := range users {
		stats := dto.UserStats{
			UserID:   string(user.ID),
			Username: user.Username,
			TeamName: string(user.Team),
		}

		for _, pullRequest := range pullRequests {
			if pullRequest.AuthorID == user.ID {
				stats.PullRequestsCreated++
			}
		}

		userStats = append(userStats, stats)
	}

	return userStats
}

func calculateTeamStats(teams []entities.Team, users []entities.User, pullRequests []entities.PullRequest) []dto.TeamStats {
	var teamStats []dto.TeamStats

	for _, team := range teams {
		stats := dto.TeamStats{
			TeamName: string(team.Name),
		}

		for _, user := range users {
			if user.Team == team.Name {
				stats.MemberCount++
				if user.IsActive {
					stats.ActiveMembers++
				}
			}
		}

		for _, pullRequest := range pullRequests {
			for _, user := range users {
				if user.ID == pullRequest.AuthorID && user.Team == team.Name {
					stats.PullRequestsCount++
					break
				}
			}
		}

		teamStats = append(teamStats, stats)
	}

	return teamStats
}

func calculateReviewAssignments(pullRequests []entities.PullRequest) []dto.ReviewAssignment {
	var assignments []dto.ReviewAssignment

	for _, pullRequest := range pullRequests {
		assignment := dto.ReviewAssignment{
			PullRequestID:   string(pullRequest.ID),
			PullRequestName: pullRequest.Name,
			AuthorID:        string(pullRequest.AuthorID),
			Status:          string(pullRequest.Status),
		}

		assignments = append(assignments, assignment)
	}

	return assignments
}

package services

import (
	"context"

	"pr-service/internal/app"
	"pr-service/internal/domain/entities"
	"pr-service/internal/domain/value_objects"
)

type UserService interface {
	SetActiveStatus(ctx context.Context, userID value_objects.UserID, isActive bool) (entities.User, error)
	GetUserReviews(ctx context.Context, userID value_objects.UserID) ([]entities.PullRequest, error)
}

type userService struct {
	userRepository  app.UserRepository
	pullRequestRepo app.PullRequestRepository
}

func NewUserService(userRepository app.UserRepository, pullRequestRepo app.PullRequestRepository) UserService {
	return &userService{
		userRepository:  userRepository,
		pullRequestRepo: pullRequestRepo,
	}
}

func (s *userService) SetActiveStatus(ctx context.Context, userID value_objects.UserID, isActive bool) (entities.User, error) {
	return s.userRepository.SetIsActive(ctx, userID, isActive)
}

func (s *userService) GetUserReviews(ctx context.Context, userID value_objects.UserID) ([]entities.PullRequest, error) {
	_, err := s.userRepository.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	pullRequests, err := s.pullRequestRepo.GetByReviewer(ctx, userID)
	if err != nil {
		return nil, err
	}

	return pullRequests, nil
}

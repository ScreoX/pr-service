package services

import (
	"context"

	"pr-service/internal/app"
	"pr-service/internal/domain"
	"pr-service/internal/domain/entities"
	"pr-service/internal/domain/value_objects"
)

type PullRequestService interface {
	Create(ctx context.Context, pullRequestID value_objects.PullRequestID, pullRequestName string, authorID value_objects.UserID) (*entities.PullRequest, error)
	Merge(ctx context.Context, pullRequestID value_objects.PullRequestID) (*entities.PullRequest, error)
	ReassignReviewer(ctx context.Context, pullRequestID value_objects.PullRequestID, oldReviewerID value_objects.UserID) (*entities.PullRequest, value_objects.UserID, error)
}

type pullRequestService struct {
	userRepository        app.UserRepository
	teamRepository        app.TeamRepository
	pullRequestRepository app.PullRequestRepository
	txManager             app.TxManager
	timeProvider          app.TimeProvider
	random                app.RandomProvider
}

func NewPullRequestService(userRepository app.UserRepository, teamRepository app.TeamRepository, pullRequestRepository app.PullRequestRepository, txManager app.TxManager, timeProvider app.TimeProvider, randomProvider app.RandomProvider) PullRequestService {
	return &pullRequestService{
		userRepository:        userRepository,
		teamRepository:        teamRepository,
		pullRequestRepository: pullRequestRepository,
		txManager:             txManager,
		timeProvider:          timeProvider,
		random:                randomProvider,
	}
}

func (s *pullRequestService) Create(ctx context.Context, pullRequestID value_objects.PullRequestID, pullRequestName string, authorID value_objects.UserID) (*entities.PullRequest, error) {
	if s.txManager == nil {
		return nil, app.ErrTransactionRequired
	}

	var resultPullRequest *entities.PullRequest

	operation := func(ctx context.Context) error {
		_, err := s.pullRequestRepository.GetByID(ctx, pullRequestID)
		if err == nil {
			return domain.ErrPRExists
		}

		author, err := s.userRepository.GetByID(ctx, authorID)
		if err != nil {
			return domain.ErrUserNotFound
		}

		team, err := s.teamRepository.GetByName(ctx, author.Team)
		if err != nil {
			return err
		}

		if !author.IsActive {
			return domain.ErrAuthorNotActive
		}

		teamMembers, err := s.userRepository.GetUsersByTeam(ctx, team.Name)
		if err != nil {
			return err
		}

		resultPullRequest = entities.NewPullRequest(pullRequestID, pullRequestName, authorID, s.timeProvider.Now())

		activeCandidates := s.filterActiveUsersExcludeAuthor(authorID, teamMembers)
		if len(activeCandidates) > 0 {
			activeCandidatesIDs := toUserIDs(activeCandidates)
			s.random.Shuffle(len(activeCandidatesIDs), func(i, j int) {
				activeCandidatesIDs[i], activeCandidatesIDs[j] = activeCandidatesIDs[j], activeCandidatesIDs[i]
			})

			resultPullRequest.AddReviewers(activeCandidatesIDs)
		}

		if err := s.pullRequestRepository.Create(ctx, resultPullRequest); err != nil {
			return err
		}

		return nil
	}

	if err := s.txManager.Do(ctx, operation); err != nil {
		return nil, err
	}

	return resultPullRequest, nil
}

func (s *pullRequestService) Merge(ctx context.Context, pullRequestID value_objects.PullRequestID) (*entities.PullRequest, error) {
	pullRequest, err := s.pullRequestRepository.GetByID(ctx, pullRequestID)
	if err != nil {
		return nil, err
	}

	pullRequest.Merge(s.timeProvider.Now())

	if err := s.pullRequestRepository.Save(ctx, pullRequest); err != nil {
		return nil, err
	}

	return pullRequest, nil
}

func (s *pullRequestService) ReassignReviewer(ctx context.Context, pullRequestID value_objects.PullRequestID, oldReviewerID value_objects.UserID) (*entities.PullRequest, value_objects.UserID, error) {
	if s.txManager == nil {
		return nil, "", app.ErrTransactionRequired
	}

	var resultPullRequest *entities.PullRequest
	var newReviewerID value_objects.UserID

	operation := func(ctx context.Context) error {
		pullRequest, err := s.pullRequestRepository.GetByID(ctx, pullRequestID)
		if err != nil {
			return err
		}

		if pullRequest.IsMerged() {
			return domain.ErrPRMerged
		}

		_, err = s.userRepository.GetByID(ctx, oldReviewerID)
		if err != nil {
			return domain.ErrUserNotFound
		}

		if !pullRequest.IsReviewer(oldReviewerID) {
			return domain.ErrNotAssigned
		}

		author, err := s.userRepository.GetByID(ctx, pullRequest.AuthorID)
		if err != nil {
			return err
		}

		team, err := s.teamRepository.GetByName(ctx, author.Team)
		if err != nil {
			return err
		}

		teamMembers, err := s.userRepository.GetUsersByTeam(ctx, team.Name)
		if err != nil {
			return err
		}

		activeCandidates := s.filterActiveUsersExcludeAuthorAndReviewer(pullRequest, pullRequest.AuthorID, oldReviewerID, teamMembers)
		if len(activeCandidates) == 0 {
			return domain.ErrNoCandidate
		}

		activeCandidatesIDs := toUserIDs(activeCandidates)
		newReviewerID = activeCandidatesIDs[s.random.Intn(len(activeCandidatesIDs))]

		err = s.pullRequestRepository.ReassignReviewer(ctx, pullRequestID, oldReviewerID, newReviewerID)
		if err != nil {
			return err
		}

		resultPullRequest, err = s.pullRequestRepository.GetByID(ctx, pullRequestID)
		if err != nil {
			return err
		}

		return nil
	}

	if err := s.txManager.Do(ctx, operation); err != nil {
		return nil, "", err
	}

	return resultPullRequest, newReviewerID, nil
}

func (s *pullRequestService) filterActiveUsersExcludeAuthor(authorID value_objects.UserID, candidates []entities.User) []entities.User {
	var activeCandidates []entities.User

	for _, candidate := range candidates {
		if candidate.IsActive && candidate.ID != authorID {
			activeCandidates = append(activeCandidates, candidate)
		}
	}

	return activeCandidates
}

func (s *pullRequestService) filterActiveUsersExcludeAuthorAndReviewer(pullRequest *entities.PullRequest, authorID, reviewerID value_objects.UserID, candidates []entities.User) []entities.User {
	var activeCandidates []entities.User

	for _, candidate := range candidates {
		isAlreadyReviewer := false
		for _, existingReviewer := range pullRequest.Reviewers() {
			if existingReviewer == candidate.ID && existingReviewer != reviewerID {
				isAlreadyReviewer = true
				break
			}
		}

		if candidate.IsActive && candidate.ID != authorID && candidate.ID != reviewerID && !isAlreadyReviewer {
			activeCandidates = append(activeCandidates, candidate)
		}
	}

	return activeCandidates
}

func toUserIDs(users []entities.User) []value_objects.UserID {
	userIDs := make([]value_objects.UserID, len(users))

	for i, user := range users {
		userIDs[i] = user.ID
	}

	return userIDs
}

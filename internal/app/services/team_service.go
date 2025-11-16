package services

import (
	"context"
	"errors"

	"pr-service/internal/app"
	"pr-service/internal/domain"
	"pr-service/internal/domain/entities"
	"pr-service/internal/domain/value_objects"
)

type TeamService interface {
	Create(ctx context.Context, teamName value_objects.TeamName, members []entities.User) (entities.Team, []entities.User, error)
	GetByName(ctx context.Context, teamName value_objects.TeamName) (entities.Team, []entities.User, error)
}

type teamService struct {
	userRepository app.UserRepository
	teamRepository app.TeamRepository
	txManager      app.TxManager
}

func NewTeamService(userRepository app.UserRepository, teamRepository app.TeamRepository, txManager app.TxManager) TeamService {
	return &teamService{
		userRepository: userRepository,
		teamRepository: teamRepository,
		txManager:      txManager,
	}
}

func (s *teamService) Create(ctx context.Context, teamName value_objects.TeamName, members []entities.User) (entities.Team, []entities.User, error) {
	if s.txManager == nil {
		return entities.Team{}, nil, app.ErrTransactionRequired
	}

	var resultTeam entities.Team
	var resultTeamMembers []entities.User

	operation := func(ctx context.Context) error {
		_, err := s.teamRepository.GetByName(ctx, teamName)
		if err == nil {
			return domain.ErrTeamExists
		} else if !errors.Is(err, domain.ErrTeamNotFound) {
			return err
		}

		if createErr := s.teamRepository.Create(ctx, entities.Team{Name: teamName}); createErr != nil {
			return createErr
		}

		if upsertErr := s.userRepository.UpsertMembers(ctx, teamName, members); upsertErr != nil {
			return upsertErr
		}

		resultTeamMembers, err = s.userRepository.GetUsersByTeam(ctx, teamName)
		if err != nil {
			return err
		}

		resultTeam = entities.Team{Name: teamName}

		return nil
	}

	return resultTeam, resultTeamMembers, s.txManager.Do(ctx, operation)
}

func (s *teamService) GetByName(ctx context.Context, teamName value_objects.TeamName) (entities.Team, []entities.User, error) {
	team, err := s.teamRepository.GetByName(ctx, teamName)
	if err != nil {
		return entities.Team{}, nil, err
	}

	users, err := s.userRepository.GetUsersByTeam(ctx, teamName)
	if err != nil {
		return entities.Team{}, nil, err
	}

	return team, users, nil
}

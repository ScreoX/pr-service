package services

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"pr-service/internal/app"
	"pr-service/internal/app/services/mocks"
	"pr-service/internal/domain"
	"pr-service/internal/domain/entities"
	"pr-service/internal/domain/value_objects"
)

func TestTeamService_Create(t *testing.T) {
	ctx := context.Background()

	t.Run("successfully create team with members", func(t *testing.T) {
		userRepository := &mocks.UserRepository{}
		teamRepository := &mocks.TeamRepository{}
		txManager := &mocks.TxManager{}

		teamName := value_objects.TeamName("backend")
		members := []entities.User{
			{
				ID:       value_objects.UserID("user1"),
				Username: "Alice",
				Team:     value_objects.TeamName("backend"),
				IsActive: true,
			},
		}

		teamRepository.On("GetByName", ctx, teamName).Once().
			Return(entities.Team{}, domain.ErrTeamNotFound)

		teamRepository.On("Create", ctx, entities.Team{Name: teamName}).Once().
			Return(nil)

		userRepository.On("UpsertMembers", ctx, teamName, members).Once().
			Return(nil)

		userRepository.On("GetUsersByTeam", ctx, teamName).Once().
			Return(members, nil)

		txManager.On("Do", ctx, mock.AnythingOfType("func(context.Context) error")).Once().
			Return(nil)

		service := NewTeamService(userRepository, teamRepository, txManager)
		resultTeam, resultUsers, err := service.Create(ctx, teamName, members)

		assert.NoError(t, err)
		assert.Equal(t, entities.Team{Name: teamName}, resultTeam)
		assert.Equal(t, members, resultUsers)

		userRepository.AssertExpectations(t)
		teamRepository.AssertExpectations(t)
		txManager.AssertExpectations(t)
	})

	t.Run("return error when txManager is nil", func(t *testing.T) {
		userRepository := &mocks.UserRepository{}
		teamRepository := &mocks.TeamRepository{}

		teamName := value_objects.TeamName("backend")
		members := []entities.User{
			{
				ID:       value_objects.UserID("user1"),
				Username: "Alice",
				Team:     value_objects.TeamName("backend"),
				IsActive: true,
			},
		}

		service := NewTeamService(userRepository, teamRepository, nil)
		resultTeam, resultUsers, err := service.Create(ctx, teamName, members)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, app.ErrTransactionRequired))
		assert.Equal(t, entities.Team{}, resultTeam)
		assert.Nil(t, resultUsers)

		userRepository.AssertNotCalled(t, "UpsertMembers")
		userRepository.AssertNotCalled(t, "GetUsersByTeam")
		teamRepository.AssertNotCalled(t, "GetByName")
		teamRepository.AssertNotCalled(t, "Create")
	})

	t.Run("return transaction error when txManager fails", func(t *testing.T) {
		userRepository := &mocks.UserRepository{}
		teamRepository := &mocks.TeamRepository{}
		txManager := &mocks.TxManager{}

		teamName := value_objects.TeamName("backend")
		members := []entities.User{
			{
				ID:       value_objects.UserID("user1"),
				Username: "Alice",
				Team:     value_objects.TeamName("backend"),
				IsActive: true,
			},
		}

		txManager.On("Do", ctx, mock.AnythingOfType("func(context.Context) error")).Once().
			Return(errors.New("transaction failed"))

		service := NewTeamService(userRepository, teamRepository, txManager)
		resultTeam, resultUsers, err := service.Create(ctx, teamName, members)

		assert.Error(t, err)
		assert.Equal(t, "transaction failed", err.Error())
		assert.Equal(t, entities.Team{}, resultTeam)
		assert.Nil(t, resultUsers)

		userRepository.AssertNotCalled(t, "UpsertMembers")
		userRepository.AssertNotCalled(t, "GetUsersByTeam")
		teamRepository.AssertNotCalled(t, "GetByName")
		teamRepository.AssertNotCalled(t, "Create")
		txManager.AssertExpectations(t)
	})

	t.Run("return error when operation returns error", func(t *testing.T) {
		userRepository := &mocks.UserRepository{}
		teamRepository := &mocks.TeamRepository{}
		txManager := &mocks.TxManager{}

		teamName := value_objects.TeamName("backend")
		members := []entities.User{
			{
				ID:       value_objects.UserID("user1"),
				Username: "Alice",
				Team:     value_objects.TeamName("backend"),
				IsActive: true,
			},
		}

		txManager.On("Do", ctx, mock.AnythingOfType("func(context.Context) error")).Once().
			Return(domain.ErrTeamExists)

		service := NewTeamService(userRepository, teamRepository, txManager)
		resultTeam, resultUsers, err := service.Create(ctx, teamName, members)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrTeamExists))
		assert.Equal(t, entities.Team{}, resultTeam)
		assert.Nil(t, resultUsers)

		txManager.AssertExpectations(t)
	})
}

func TestTeamService_GetByName(t *testing.T) {
	ctx := context.Background()

	t.Run("successfully get team with members", func(t *testing.T) {
		userRepository := &mocks.UserRepository{}
		teamRepository := &mocks.TeamRepository{}

		teamName := value_objects.TeamName("backend")
		expectedTeam := entities.Team{Name: teamName}
		expectedUsers := []entities.User{
			{
				ID:       value_objects.UserID("user1"),
				Username: "Alice",
				Team:     teamName,
				IsActive: true,
			},
		}

		teamRepository.On("GetByName", ctx, teamName).
			Return(expectedTeam, nil)

		userRepository.On("GetUsersByTeam", ctx, teamName).
			Return(expectedUsers, nil)

		service := NewTeamService(userRepository, teamRepository, nil)
		resultTeam, resultUsers, err := service.GetByName(ctx, teamName)

		assert.NoError(t, err)
		assert.Equal(t, expectedTeam, resultTeam)
		assert.Equal(t, expectedUsers, resultUsers)

		teamRepository.AssertExpectations(t)
		userRepository.AssertExpectations(t)
	})

	t.Run("return error when team not found", func(t *testing.T) {
		userRepository := &mocks.UserRepository{}
		teamRepository := &mocks.TeamRepository{}

		teamName := value_objects.TeamName("nonexistent")

		teamRepository.On("GetByName", ctx, teamName).
			Return(entities.Team{}, domain.ErrTeamNotFound)

		service := NewTeamService(userRepository, teamRepository, nil)
		resultTeam, resultUsers, err := service.GetByName(ctx, teamName)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrTeamNotFound))
		assert.Equal(t, entities.Team{}, resultTeam)
		assert.Nil(t, resultUsers)

		teamRepository.AssertExpectations(t)
		userRepository.AssertNotCalled(t, "GetUsersByTeam")
	})

	t.Run("return error when get team fails", func(t *testing.T) {
		userRepository := &mocks.UserRepository{}
		teamRepository := &mocks.TeamRepository{}

		teamName := value_objects.TeamName("backend")

		teamRepository.On("GetByName", ctx, teamName).
			Return(entities.Team{}, errors.New("database error"))

		service := NewTeamService(userRepository, teamRepository, nil)
		resultTeam, resultUsers, err := service.GetByName(ctx, teamName)

		assert.Error(t, err)
		assert.Equal(t, "database error", err.Error())
		assert.Equal(t, entities.Team{}, resultTeam)
		assert.Nil(t, resultUsers)

		teamRepository.AssertExpectations(t)
		userRepository.AssertNotCalled(t, "GetUsersByTeam")
	})

	t.Run("return empty users when team has no members", func(t *testing.T) {
		userRepository := &mocks.UserRepository{}
		teamRepository := &mocks.TeamRepository{}

		teamName := value_objects.TeamName("empty-team")
		expectedTeam := entities.Team{Name: teamName}

		teamRepository.On("GetByName", ctx, teamName).
			Return(expectedTeam, nil)

		userRepository.On("GetUsersByTeam", ctx, teamName).
			Return([]entities.User{}, nil)

		service := NewTeamService(userRepository, teamRepository, nil)
		resultTeam, resultUsers, err := service.GetByName(ctx, teamName)

		assert.NoError(t, err)
		assert.Equal(t, expectedTeam, resultTeam)
		assert.Equal(t, []entities.User{}, resultUsers)

		teamRepository.AssertExpectations(t)
		userRepository.AssertExpectations(t)
	})
}

func TestTeamService_Create_Validation(t *testing.T) {
	ctx := context.Background()

	t.Run("return empty team when creation fails", func(t *testing.T) {
		userRepository := &mocks.UserRepository{}
		teamRepository := &mocks.TeamRepository{}
		txManager := &mocks.TxManager{}

		teamName := value_objects.TeamName("backend")
		members := []entities.User{
			{
				ID:       value_objects.UserID("user1"),
				Username: "Alice",
				Team:     value_objects.TeamName("backend"),
				IsActive: true,
			},
		}

		txManager.On("Do", ctx, mock.AnythingOfType("func(context.Context) error")).Once().
			Return(domain.ErrTeamExists)

		service := NewTeamService(userRepository, teamRepository, txManager)
		resultTeam, resultUsers, err := service.Create(ctx, teamName, members)

		assert.Error(t, err)
		assert.Equal(t, entities.Team{}, resultTeam)
		assert.Nil(t, resultUsers)
		assert.True(t, errors.Is(err, domain.ErrTeamExists))
	})

	t.Run("return empty values on any transaction error", func(t *testing.T) {
		userRepository := &mocks.UserRepository{}
		teamRepository := &mocks.TeamRepository{}
		txManager := &mocks.TxManager{}

		teamName := value_objects.TeamName("backend")
		members := []entities.User{
			{
				ID:       value_objects.UserID("user1"),
				Username: "Alice",
				Team:     value_objects.TeamName("backend"),
				IsActive: true,
			},
		}

		txManager.On("Do", ctx, mock.AnythingOfType("func(context.Context) error")).Once().
			Return(errors.New("any error"))

		service := NewTeamService(userRepository, teamRepository, txManager)
		resultTeam, resultUsers, err := service.Create(ctx, teamName, members)

		assert.Error(t, err)
		assert.Equal(t, entities.Team{}, resultTeam)
		assert.Nil(t, resultUsers)
	})
}

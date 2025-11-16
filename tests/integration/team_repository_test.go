package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"pr-service/internal/domain"
	"pr-service/internal/domain/entities"
	"pr-service/internal/domain/value_objects"
	"pr-service/internal/infrastructure/postgres/repositories"
	"pr-service/tests/integration/helpers"
)

func TestTeamRepository_Create_Success(t *testing.T) {
	db := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, db)

	repository := repositories.NewTeamRepository(db)
	ctx := context.Background()

	team := entities.Team{
		Name: value_objects.TeamName("backend"),
	}

	err := repository.Create(ctx, team)

	assert.NoError(t, err)

	exists, err := helpers.TeamExists(db, "backend")
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestTeamRepository_Create_DuplicateTeam(t *testing.T) {
	db := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, db)

	repository := repositories.NewTeamRepository(db)
	ctx := context.Background()

	team := entities.Team{
		Name: value_objects.TeamName("backend"),
	}

	err := repository.Create(ctx, team)
	assert.NoError(t, err)

	err = repository.Create(ctx, team)
	assert.Error(t, err)
}

func TestTeamRepository_GetByName_Success(t *testing.T) {
	db := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, db)

	repository := repositories.NewTeamRepository(db)
	ctx := context.Background()

	teamName := value_objects.TeamName("payments")
	err := helpers.InsertTestTeam(db, "payments", "payments")
	require.NoError(t, err)

	team, err := repository.GetByName(ctx, teamName)

	assert.NoError(t, err)
	assert.Equal(t, teamName, team.Name)
}

func TestTeamRepository_GetByName_NotFound(t *testing.T) {
	db := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, db)

	repository := repositories.NewTeamRepository(db)
	ctx := context.Background()

	nonExistentName := value_objects.TeamName("non-existent-team")

	team, err := repository.GetByName(ctx, nonExistentName)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrTeamNotFound, err)
	assert.Equal(t, entities.Team{}, team)
}

func TestTeamRepository_GetAll_Success(t *testing.T) {
	db := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, db)

	repository := repositories.NewTeamRepository(db)
	ctx := context.Background()

	teams := []string{"backend", "frontend", "devops", "qa"}

	for _, teamName := range teams {
		err := helpers.InsertTestTeam(db, teamName, teamName)
		require.NoError(t, err)
	}

	allTeams, err := repository.GetAll(ctx)

	assert.NoError(t, err)
	assert.Len(t, allTeams, 4)

	teamMap := make(map[string]entities.Team)
	for _, team := range allTeams {
		teamMap[string(team.Name)] = team
	}

	assert.Contains(t, teamMap, "backend")
	assert.Contains(t, teamMap, "frontend")
	assert.Contains(t, teamMap, "devops")
	assert.Contains(t, teamMap, "qa")
}

func TestTeamRepository_GetAll_Empty(t *testing.T) {
	db := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, db)

	repository := repositories.NewTeamRepository(db)
	ctx := context.Background()

	allTeams, err := repository.GetAll(ctx)

	assert.NoError(t, err)
	assert.Empty(t, allTeams)
}

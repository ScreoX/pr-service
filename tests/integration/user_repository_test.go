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

func TestUserRepository_GetByID_Success(t *testing.T) {
	db := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, db)

	repository := repositories.NewUserRepository(db)
	ctx := context.Background()

	userID := value_objects.UserID("user-1")
	err := helpers.InsertTestUser(db, string(userID), "Test User", "backend", true)
	require.NoError(t, err)

	user, err := repository.GetByID(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, userID, user.ID)
	assert.Equal(t, "Test User", user.Username)
	assert.Equal(t, value_objects.TeamName("backend"), user.Team)
	assert.True(t, user.IsActive)
}

func TestUserRepository_GetByID_NotFound(t *testing.T) {
	db := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, db)

	repository := repositories.NewUserRepository(db)
	ctx := context.Background()

	nonExistentID := value_objects.UserID("non-existent-user")

	user, err := repository.GetByID(ctx, nonExistentID)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrUserNotFound, err)
	assert.Equal(t, entities.User{}, user)
}

func TestUserRepository_SetIsActive_Success(t *testing.T) {
	db := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, db)

	repository := repositories.NewUserRepository(db)
	ctx := context.Background()

	userID := value_objects.UserID("user-2")
	err := helpers.InsertTestUser(db, string(userID), "Test User 2", "backend", true)
	require.NoError(t, err)

	updatedUser, err := repository.SetIsActive(ctx, userID, false)

	assert.NoError(t, err)
	assert.Equal(t, userID, updatedUser.ID)
	assert.False(t, updatedUser.IsActive)

	isActive, err := helpers.GetUserActivity(db, string(userID))
	assert.NoError(t, err)
	assert.False(t, isActive)
}

func TestUserRepository_SetIsActive_UserNotFound(t *testing.T) {
	db := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, db)

	repository := repositories.NewUserRepository(db)
	ctx := context.Background()

	user, err := repository.SetIsActive(ctx, "non-existent", false)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrUserNotFound, err)
	assert.Equal(t, entities.User{}, user)
}

func TestUserRepository_UpsertMembers_Success(t *testing.T) {
	db := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, db)

	repository := repositories.NewUserRepository(db)
	ctx := context.Background()

	teamName := value_objects.TeamName("payments")
	members := []entities.User{
		{
			ID:       value_objects.UserID("u1"),
			Username: "Alice",
			Team:     teamName,
			IsActive: true,
		},
		{
			ID:       value_objects.UserID("u2"),
			Username: "Bob",
			Team:     teamName,
			IsActive: true,
		},
	}

	err := repository.UpsertMembers(ctx, teamName, members)

	assert.NoError(t, err)

	user1, err := repository.GetByID(ctx, "u1")
	assert.NoError(t, err)
	assert.Equal(t, "Alice", user1.Username)
	assert.Equal(t, teamName, user1.Team)

	user2, err := repository.GetByID(ctx, "u2")
	assert.NoError(t, err)
	assert.Equal(t, "Bob", user2.Username)
	assert.Equal(t, teamName, user2.Team)
}

func TestUserRepository_UpsertMembers_UpdateExisting(t *testing.T) {
	db := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, db)

	repository := repositories.NewUserRepository(db)
	ctx := context.Background()

	teamName := value_objects.TeamName("backend")

	err := helpers.InsertTestUser(db, "u1", "Old Name", "old-team", true)
	require.NoError(t, err)

	members := []entities.User{
		{
			ID:       value_objects.UserID("u1"),
			Username: "New Name",
			Team:     teamName,
			IsActive: false,
		},
	}

	err = repository.UpsertMembers(ctx, teamName, members)

	assert.NoError(t, err)

	user, err := repository.GetByID(ctx, "u1")
	assert.NoError(t, err)
	assert.Equal(t, "New Name", user.Username)
	assert.Equal(t, teamName, user.Team)
	assert.False(t, user.IsActive)
}

func TestUserRepository_GetUsersByTeam_Success(t *testing.T) {
	db := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, db)

	repository := repositories.NewUserRepository(db)
	ctx := context.Background()

	teamName := value_objects.TeamName("backend")

	err := helpers.InsertTestUser(db, "u1", "User 1", string(teamName), true)
	require.NoError(t, err)

	err = helpers.InsertTestUser(db, "u2", "User 2", string(teamName), true)
	require.NoError(t, err)

	err = helpers.InsertTestUser(db, "u3", "User 3", "frontend", true)
	require.NoError(t, err)

	users, err := repository.GetUsersByTeam(ctx, teamName)

	assert.NoError(t, err)
	assert.Len(t, users, 2)

	userNames := make([]string, 0, len(users))
	for _, user := range users {
		userNames = append(userNames, user.Username)
		assert.Equal(t, teamName, user.Team)
	}

	assert.Contains(t, userNames, "User 1")
	assert.Contains(t, userNames, "User 2")
	assert.NotContains(t, userNames, "User 3")
}

func TestUserRepository_GetUsersByTeam_EmptyTeam(t *testing.T) {
	db := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, db)

	repository := repositories.NewUserRepository(db)
	ctx := context.Background()

	teamName := value_objects.TeamName("empty-team")

	users, err := repository.GetUsersByTeam(ctx, teamName)

	assert.NoError(t, err)
	assert.Empty(t, users)
}

func TestUserRepository_GetAll_Success(t *testing.T) {
	db := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, db)

	repository := repositories.NewUserRepository(db)
	ctx := context.Background()

	users := []struct {
		id       string
		username string
		team     string
		isActive bool
	}{
		{"user-1", "Alice", "backend", true},
		{"user-2", "Bob", "backend", false},
		{"user-3", "Charlie", "frontend", true},
		{"user-4", "David", "frontend", true},
	}

	for _, user := range users {
		err := helpers.InsertTestUser(db, user.id, user.username, user.team, user.isActive)
		require.NoError(t, err)
	}

	allUsers, err := repository.GetAll(ctx)

	assert.NoError(t, err)
	assert.Len(t, allUsers, 4)

	userMap := make(map[string]entities.User)
	for _, user := range allUsers {
		userMap[string(user.ID)] = user
	}

	assert.Contains(t, userMap, "user-1")
	assert.Contains(t, userMap, "user-2")
	assert.Contains(t, userMap, "user-3")
	assert.Contains(t, userMap, "user-4")

	assert.Equal(t, "Alice", userMap["user-1"].Username)
	assert.Equal(t, value_objects.TeamName("backend"), userMap["user-1"].Team)
	assert.True(t, userMap["user-1"].IsActive)

	assert.Equal(t, "Bob", userMap["user-2"].Username)
	assert.Equal(t, value_objects.TeamName("backend"), userMap["user-2"].Team)
	assert.False(t, userMap["user-2"].IsActive)
}

func TestUserRepository_GetAll_Empty(t *testing.T) {
	db := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, db)

	repository := repositories.NewUserRepository(db)
	ctx := context.Background()

	allUsers, err := repository.GetAll(ctx)

	assert.NoError(t, err)
	assert.Empty(t, allUsers)
}

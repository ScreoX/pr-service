package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"pr-service/internal/domain"
	"pr-service/internal/domain/entities"
	"pr-service/internal/domain/value_objects"
	"pr-service/internal/infrastructure/postgres/repositories"
	"pr-service/tests/integration/helpers"
)

func TestPullRequestRepository_Create_Success(t *testing.T) {
	db := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, db)

	repository := repositories.NewPullRequestRepository(db)
	ctx := context.Background()

	err := helpers.InsertTestUser(db, "author-1", "Author", "team1", true)
	require.NoError(t, err)

	pullRequest := &entities.PullRequest{
		ID:       value_objects.PullRequestID("pull-request-1"),
		Name:     "Test PR",
		AuthorID: value_objects.UserID("author-1"),
		Status:   entities.StatusOpen,
	}

	err = repository.Create(ctx, pullRequest)

	assert.NoError(t, err)

	exists, err := helpers.PullRequestExists(db, "pull-request-1")
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestPullRequestRepository_Create_WithReviewers(t *testing.T) {
	db := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, db)

	repository := repositories.NewPullRequestRepository(db)
	ctx := context.Background()

	err := helpers.InsertTestUser(db, "author-1", "Author", "team1", true)
	require.NoError(t, err)

	err = helpers.InsertTestUser(db, "user1", "User 1", "team1", true)
	require.NoError(t, err)

	err = helpers.InsertTestUser(db, "user2", "User 2", "team1", true)
	require.NoError(t, err)

	pullRequest := &entities.PullRequest{
		ID:       value_objects.PullRequestID("pull-request-1"),
		Name:     "Test PR",
		AuthorID: value_objects.UserID("author-1"),
		Status:   entities.StatusOpen,
	}
	pullRequest.AddReviewers([]value_objects.UserID{"user1", "user2"})

	err = repository.Create(ctx, pullRequest)

	assert.NoError(t, err)

	reviewers, err := helpers.GetPullRequestReviewers(db, "pull-request-1")
	assert.NoError(t, err)
	assert.Len(t, reviewers, 2)
	assert.Contains(t, reviewers, "user1")
	assert.Contains(t, reviewers, "user2")
}

func TestPullRequestRepository_GetByID_Success(t *testing.T) {
	db := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, db)

	repository := repositories.NewPullRequestRepository(db)
	ctx := context.Background()

	err := helpers.InsertTestUser(db, "author-1", "Author", "team1", true)
	require.NoError(t, err)

	err = helpers.InsertTestPullRequest(db, "pull-request-1", "Test PR", "author-1", "OPEN")
	require.NoError(t, err)

	pullRequest, err := repository.GetByID(ctx, "pull-request-1")

	assert.NoError(t, err)
	assert.Equal(t, value_objects.PullRequestID("pull-request-1"), pullRequest.ID)
	assert.Equal(t, "Test PR", pullRequest.Name)
	assert.Equal(t, value_objects.UserID("author-1"), pullRequest.AuthorID)
	assert.Equal(t, entities.StatusOpen, pullRequest.Status)
}

func TestPullRequestRepository_GetByID_WithReviewers(t *testing.T) {
	db := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, db)

	repository := repositories.NewPullRequestRepository(db)
	ctx := context.Background()

	err := helpers.InsertTestUser(db, "author-1", "Author", "team1", true)
	require.NoError(t, err)

	err = helpers.InsertTestUser(db, "reviewer-1", "Reviewer 1", "team1", true)
	require.NoError(t, err)

	err = helpers.InsertTestUser(db, "reviewer-2", "Reviewer 2", "team1", true)
	require.NoError(t, err)

	err = helpers.InsertTestPullRequest(db, "pull-request-1", "Test PR", "author-1", "OPEN")
	require.NoError(t, err)

	err = helpers.AddReviewerToPullRequest(db, "pull-request-1", "reviewer-1")
	require.NoError(t, err)

	err = helpers.AddReviewerToPullRequest(db, "pull-request-1", "reviewer-2")
	require.NoError(t, err)

	pullRequest, err := repository.GetByID(ctx, "pull-request-1")

	assert.NoError(t, err)
	assert.Len(t, pullRequest.Reviewers(), 2)
	assert.Contains(t, pullRequest.Reviewers(), value_objects.UserID("reviewer-1"))
	assert.Contains(t, pullRequest.Reviewers(), value_objects.UserID("reviewer-2"))
}

func TestPullRequestRepository_GetByID_NotFound(t *testing.T) {
	db := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, db)

	repository := repositories.NewPullRequestRepository(db)
	ctx := context.Background()

	pullRequest, err := repository.GetByID(ctx, "non-existent-pr")

	assert.Error(t, err)
	assert.Equal(t, domain.ErrPRNotFound, err)
	assert.Nil(t, pullRequest)
}

func TestPullRequestRepository_Save_UpdateStatus(t *testing.T) {
	db := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, db)

	repository := repositories.NewPullRequestRepository(db)
	ctx := context.Background()

	err := helpers.InsertTestUser(db, "author-1", "Author", "team1", true)
	require.NoError(t, err)

	err = helpers.InsertTestPullRequest(db, "pull-request-1", "Test PR", "author-1", "OPEN")
	require.NoError(t, err)

	pullRequest, err := repository.GetByID(ctx, "pull-request-1")
	assert.NoError(t, err)

	mergedAt := time.Now()
	pullRequest.Status = entities.StatusMerged
	pullRequest.MergedAt = &mergedAt

	err = repository.Save(ctx, pullRequest)

	assert.NoError(t, err)

	status, err := helpers.GetPullRequestStatus(db, "pull-request-1")
	assert.NoError(t, err)
	assert.Equal(t, "MERGED", status)
}

func TestPullRequestRepository_Save_UpdateReviewers(t *testing.T) {
	db := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, db)

	repository := repositories.NewPullRequestRepository(db)
	ctx := context.Background()

	err := helpers.InsertTestUser(db, "author-1", "Author", "team1", true)
	require.NoError(t, err)

	err = helpers.InsertTestUser(db, "old-reviewer", "Old Reviewer", "team1", true)
	require.NoError(t, err)

	err = helpers.InsertTestUser(db, "new-reviewer-1", "New Reviewer 1", "team1", true)
	require.NoError(t, err)

	err = helpers.InsertTestUser(db, "new-reviewer-2", "New Reviewer 2", "team1", true)
	require.NoError(t, err)

	err = helpers.InsertTestPullRequest(db, "pull-request-1", "Test PR", "author-1", "OPEN")
	require.NoError(t, err)

	err = helpers.AddReviewerToPullRequest(db, "pull-request-1", "old-reviewer")
	require.NoError(t, err)

	pullRequest, err := repository.GetByID(ctx, "pull-request-1")
	assert.NoError(t, err)

	updatedPullRequest := &entities.PullRequest{
		ID:       pullRequest.ID,
		Name:     pullRequest.Name,
		AuthorID: pullRequest.AuthorID,
		Status:   entities.StatusMerged,
		MergedAt: pullRequest.MergedAt,
	}

	err = repository.Save(ctx, updatedPullRequest)
	assert.NoError(t, err)

	reviewers, err := helpers.GetPullRequestReviewers(db, "pull-request-1")
	assert.NoError(t, err)

	assert.Len(t, reviewers, 1)
	assert.Contains(t, reviewers, "old-reviewer")
	assert.NotContains(t, reviewers, "new-reviewer-1")
	assert.NotContains(t, reviewers, "new-reviewer-2")

	updatedPullRequestFromDB, err := repository.GetByID(ctx, "pull-request-1")
	assert.NoError(t, err)
	assert.Equal(t, entities.StatusMerged, updatedPullRequestFromDB.Status)
}

func TestPullRequestRepository_ReassignReviewer(t *testing.T) {
	db := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, db)

	repository := repositories.NewPullRequestRepository(db)
	ctx := context.Background()

	err := helpers.InsertTestUser(db, "author-1", "Author", "team1", true)
	require.NoError(t, err)

	err = helpers.InsertTestUser(db, "old-reviewer", "Old Reviewer", "team1", true)
	require.NoError(t, err)

	err = helpers.InsertTestUser(db, "new-reviewer", "New Reviewer", "team1", true)
	require.NoError(t, err)

	err = helpers.InsertTestPullRequest(db, "pull-request-1", "Test PR", "author-1", "OPEN")
	require.NoError(t, err)

	err = helpers.AddReviewerToPullRequest(db, "pull-request-1", "old-reviewer")
	require.NoError(t, err)

	err = repository.ReassignReviewer(ctx, "pull-request-1", "old-reviewer", "new-reviewer")
	assert.NoError(t, err)

	reviewers, err := helpers.GetPullRequestReviewers(db, "pull-request-1")
	assert.NoError(t, err)
	assert.Len(t, reviewers, 1)
	assert.Contains(t, reviewers, "new-reviewer")
	assert.NotContains(t, reviewers, "old-reviewer")
}

func TestPullRequestRepository_GetByReviewer(t *testing.T) {
	db := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, db)

	repository := repositories.NewPullRequestRepository(db)
	ctx := context.Background()

	err := helpers.InsertTestUser(db, "author-1", "Author 1", "team1", true)
	require.NoError(t, err)

	err = helpers.InsertTestUser(db, "author-2", "Author 2", "team1", true)
	require.NoError(t, err)

	err = helpers.InsertTestUser(db, "reviewer-1", "Reviewer 1", "team1", true)
	require.NoError(t, err)

	err = helpers.InsertTestPullRequest(db, "pull-request-1", "PR 1", "author-1", "OPEN")
	require.NoError(t, err)

	err = helpers.InsertTestPullRequest(db, "pull-request-2", "PR 2", "author-2", "OPEN")
	require.NoError(t, err)

	err = helpers.AddReviewerToPullRequest(db, "pull-request-1", "reviewer-1")
	require.NoError(t, err)

	err = helpers.AddReviewerToPullRequest(db, "pull-request-2", "reviewer-1")
	require.NoError(t, err)

	pullRequests, err := repository.GetByReviewer(ctx, "reviewer-1")

	assert.NoError(t, err)
	assert.Len(t, pullRequests, 2)

	pullRequestNames := make([]string, 0, len(pullRequests))

	for _, pullRequest := range pullRequests {
		pullRequestNames = append(pullRequestNames, pullRequest.Name)
	}

	assert.Contains(t, pullRequestNames, "PR 1")
	assert.Contains(t, pullRequestNames, "PR 2")
}

func TestPullRequestRepository_GetByReviewer_Empty(t *testing.T) {
	db := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, db)

	repository := repositories.NewPullRequestRepository(db)
	ctx := context.Background()

	pullRequests, err := repository.GetByReviewer(ctx, "non-reviewer")

	assert.NoError(t, err)
	assert.Empty(t, pullRequests)
}

func TestPullRequestRepository_GetAll_Success(t *testing.T) {
	db := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, db)

	repository := repositories.NewPullRequestRepository(db)
	ctx := context.Background()

	users := []struct {
		id       string
		username string
		team     string
		isActive bool
	}{
		{"author-1", "Alice", "backend", true},
		{"author-2", "Bob", "frontend", true},
	}

	for _, user := range users {
		err := helpers.InsertTestUser(db, user.id, user.username, user.team, user.isActive)
		require.NoError(t, err)
	}

	pullRequests := []struct {
		id     string
		name   string
		author string
		status string
	}{
		{"pull-request-1", "Feature A", "author-1", "OPEN"},
		{"pull-request-2", "Feature B", "author-2", "OPEN"},
		{"pull-request-3", "Bug Fix", "author-1", "MERGED"},
	}

	for _, pullRequest := range pullRequests {
		err := helpers.InsertTestPullRequest(db, pullRequest.id, pullRequest.name, pullRequest.author, pullRequest.status)
		require.NoError(t, err)
	}

	allPullRequests, err := repository.GetAll(ctx)

	assert.NoError(t, err)
	assert.Len(t, allPullRequests, 3)

	pullRequestMap := make(map[string]entities.PullRequest)
	for _, pullRequest := range allPullRequests {
		pullRequestMap[string(pullRequest.ID)] = pullRequest
	}

	assert.Contains(t, pullRequestMap, "pull-request-1")
	assert.Contains(t, pullRequestMap, "pull-request-2")
	assert.Contains(t, pullRequestMap, "pull-request-3")

	firstPullRequest := pullRequestMap["pull-request-1"]
	assert.Equal(t, "Feature A", firstPullRequest.Name)
	assert.Equal(t, value_objects.UserID("author-1"), firstPullRequest.AuthorID)
	assert.Equal(t, entities.StatusOpen, firstPullRequest.Status)

	thirdPullRequest := pullRequestMap["pull-request-3"]
	assert.Equal(t, "Bug Fix", thirdPullRequest.Name)
	assert.Equal(t, entities.StatusMerged, thirdPullRequest.Status)
}

func TestPullRequestRepository_GetAll_Empty(t *testing.T) {
	db := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, db)

	repository := repositories.NewPullRequestRepository(db)
	ctx := context.Background()

	allPullRequests, err := repository.GetAll(ctx)

	assert.NoError(t, err)
	assert.Empty(t, allPullRequests)
}

func TestPullRequestRepository_GetAll_WithMixedStatus(t *testing.T) {
	db := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, db)

	repository := repositories.NewPullRequestRepository(db)
	ctx := context.Background()

	err := helpers.InsertTestUser(db, "author-1", "Author", "team1", true)
	require.NoError(t, err)

	statuses := []string{"OPEN", "MERGED", "OPEN", "MERGED"}
	for i, status := range statuses {
		pullRequestID := fmt.Sprintf("pull-request-%d", i+1)
		insertErr := helpers.InsertTestPullRequest(db, pullRequestID, fmt.Sprintf("Pull Request %d", i+1), "author-1", status)
		require.NoError(t, insertErr)
	}

	allPullRequests, err := repository.GetAll(ctx)

	assert.NoError(t, err)
	assert.Len(t, allPullRequests, 4)

	openCount := 0
	mergedCount := 0
	for _, pullRequest := range allPullRequests {
		if pullRequest.Status == entities.StatusOpen {
			openCount++
		} else if pullRequest.Status == entities.StatusMerged {
			mergedCount++
		}
	}

	assert.Equal(t, 2, openCount)
	assert.Equal(t, 2, mergedCount)
}

func TestPullRequestRepository_GetAll_Order(t *testing.T) {
	db := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, db)

	repository := repositories.NewPullRequestRepository(db)
	ctx := context.Background()

	err := helpers.InsertTestUser(db, "author-1", "Author", "team1", true)
	require.NoError(t, err)

	pullRequests := []struct {
		id     string
		name   string
		status string
	}{
		{"pull-request-1", "First PR", "OPEN"},
		{"pull-request-2", "Second PR", "OPEN"},
		{"pull-request-3", "Third PR", "MERGED"},
	}

	for _, pullRequest := range pullRequests {
		insertErr := helpers.InsertTestPullRequest(db, pullRequest.id, pullRequest.name, "author-1", pullRequest.status)
		require.NoError(t, insertErr)
	}

	allPullRequests, err := repository.GetAll(ctx)

	assert.NoError(t, err)
	assert.Len(t, allPullRequests, 3)

	pullRequestIDs := make(map[string]bool)
	for _, pullRequest := range allPullRequests {
		pullRequestIDs[string(pullRequest.ID)] = true
	}

	assert.True(t, pullRequestIDs["pull-request-1"])
	assert.True(t, pullRequestIDs["pull-request-2"])
	assert.True(t, pullRequestIDs["pull-request-3"])
}

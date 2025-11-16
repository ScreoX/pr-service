package entities

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"pr-service/internal/domain"
	"pr-service/internal/domain/value_objects"
)

func TestNewPullRequest(t *testing.T) {
	now := time.Now()
	pullRequestID := value_objects.PullRequestID("pull-request-1")
	authorID := value_objects.UserID("Artem")

	pullRequest := NewPullRequest(pullRequestID, "feature_111", authorID, now)

	assert.Equal(t, pullRequestID, pullRequest.ID)
	assert.Equal(t, "feature_111", pullRequest.Name)
	assert.Equal(t, authorID, pullRequest.AuthorID)
	assert.Equal(t, StatusOpen, pullRequest.Status)
	assert.Empty(t, pullRequest.Reviewers())
	assert.Equal(t, now, pullRequest.CreatedAt)
	assert.Nil(t, pullRequest.MergedAt)
}

func TestPullRequest_AddReviewers(t *testing.T) {
	t.Run("add reviewers when pull request is open and has capacity", func(t *testing.T) {
		pullRequest := NewPullRequest("pull-request-1", "feature_111", "Artem", time.Now())
		candidates := []value_objects.UserID{"user1", "user2"}

		added := pullRequest.AddReviewers(candidates)

		assert.Len(t, added, 2)
		assert.Len(t, pullRequest.Reviewers(), 2)
		assert.Contains(t, pullRequest.Reviewers(), value_objects.UserID("user1"))
		assert.Contains(t, pullRequest.Reviewers(), value_objects.UserID("user2"))
	})

	t.Run("skip author when in candidates", func(t *testing.T) {
		pullRequest := NewPullRequest("pull-request-1", "feature_111", "Artem", time.Now())
		candidates := []value_objects.UserID{"Artem", "user1", "user2"}

		added := pullRequest.AddReviewers(candidates)

		assert.Len(t, added, 2)
		assert.NotContains(t, added, value_objects.UserID("Artem"))
		assert.Contains(t, added, value_objects.UserID("user1"))
		assert.Contains(t, added, value_objects.UserID("user2"))
	})

	t.Run("skip already assigned reviewers", func(t *testing.T) {
		pullRequest := NewPullRequest("pull-request-1", "feature_111", "Artem", time.Now())
		pullRequest.AddReviewers([]value_objects.UserID{"user1"})

		candidates := []value_objects.UserID{"user1", "user2", "user3"}

		added := pullRequest.AddReviewers(candidates)

		assert.Len(t, added, 1)
		assert.Equal(t, value_objects.UserID("user2"), added[0])
		assert.Len(t, pullRequest.Reviewers(), 2)
	})

	t.Run("return nil when pullRequest is merged", func(t *testing.T) {
		pullRequest := NewPullRequest("pull-request-1", "feature_111", "Artem", time.Now())
		pullRequest.Merge(time.Now())
		candidates := []value_objects.UserID{"user1", "user2"}

		added := pullRequest.AddReviewers(candidates)

		assert.Nil(t, added)
		assert.Empty(t, pullRequest.Reviewers())
	})

	t.Run("return nil when pullRequest is fully assigned", func(t *testing.T) {
		pullRequest := NewPullRequest("pull-request-1", "feature_111", "Artem", time.Now())
		pullRequest.AddReviewers([]value_objects.UserID{"user1", "user2"})
		candidates := []value_objects.UserID{"user3"}

		added := pullRequest.AddReviewers(candidates)

		assert.Nil(t, added)
		assert.Len(t, pullRequest.Reviewers(), 2)
	})

	t.Run("return nil when no candidates", func(t *testing.T) {
		pullRequest := NewPullRequest("pull-request-1", "feature_111", "Artem", time.Now())

		added := pullRequest.AddReviewers([]value_objects.UserID{})

		assert.Nil(t, added)
		assert.Empty(t, pullRequest.Reviewers())
	})

	t.Run("add available candidates when not enough for full capacity", func(t *testing.T) {
		pullRequest := NewPullRequest("pull-request-1", "feature_111", "Artem", time.Now())
		candidates := []value_objects.UserID{"user1"}

		added := pullRequest.AddReviewers(candidates)

		assert.Len(t, added, 1)
		assert.Equal(t, value_objects.UserID("user1"), added[0])
		assert.Len(t, pullRequest.Reviewers(), 1)
	})
}

func TestPullRequest_ReassignReviewer(t *testing.T) {
	t.Run("successfully reassign reviewer", func(t *testing.T) {
		pullRequest := NewPullRequest("pull-request-1", "feature_111", "Artem", time.Now())
		pullRequest.AddReviewers([]value_objects.UserID{"user1", "user2"})

		err := pullRequest.ReassignReviewer("user1", "user3")

		assert.NoError(t, err)
		reviewers := pullRequest.Reviewers()
		assert.Len(t, reviewers, 2)
		assert.Contains(t, reviewers, value_objects.UserID("user3"))
		assert.Contains(t, reviewers, value_objects.UserID("user2"))
		assert.NotContains(t, reviewers, value_objects.UserID("user1"))
	})

	t.Run("fail when pullRequest is merged", func(t *testing.T) {
		pullRequest := NewPullRequest("pull-request-1", "feature_111", "Artem", time.Now())
		pullRequest.AddReviewers([]value_objects.UserID{"user1"})
		pullRequest.Merge(time.Now())

		err := pullRequest.ReassignReviewer("user1", "user2")

		assert.Error(t, err)
		assert.Equal(t, domain.ErrPRMerged, err)
	})

	t.Run("fail when old reviewer not assigned", func(t *testing.T) {
		pullRequest := NewPullRequest("pull-request-1", "feature_111", "Artem", time.Now())
		pullRequest.AddReviewers([]value_objects.UserID{"user1"})

		err := pullRequest.ReassignReviewer("user2", "user3")

		assert.Error(t, err)
		assert.Equal(t, domain.ErrNotAssigned, err)
	})

	t.Run("fail when new reviewer is author", func(t *testing.T) {
		pullRequest := NewPullRequest("pull-request-1", "feature_111", "Artem", time.Now())
		pullRequest.AddReviewers([]value_objects.UserID{"user1"})

		err := pullRequest.ReassignReviewer("user1", "Artem")

		assert.Error(t, err)
		assert.Equal(t, domain.ErrNoCandidate, err)
	})

	t.Run("fail when new reviewer is already assigned", func(t *testing.T) {
		pullRequest := NewPullRequest("pull-request-1", "feature_111", "Artem", time.Now())
		pullRequest.AddReviewers([]value_objects.UserID{"user1", "user2"})

		err := pullRequest.ReassignReviewer("user1", "user2")

		assert.Error(t, err)
		assert.Equal(t, domain.ErrNoCandidate, err)
	})
}

func TestPullRequest_Merge(t *testing.T) {
	t.Run("successfully merge open pullRequest", func(t *testing.T) {
		now := time.Now()
		pullRequest := NewPullRequest("pull-request-1", "feature_111", "Artem", time.Now().Add(-time.Hour))
		pullRequest.AddReviewers([]value_objects.UserID{"user1"})

		pullRequest.Merge(now)

		assert.Equal(t, StatusMerged, pullRequest.Status)
		assert.NotNil(t, pullRequest.MergedAt)
		assert.Equal(t, now, *pullRequest.MergedAt)
		assert.Len(t, pullRequest.Reviewers(), 1)
	})

	t.Run("idempotent merge - no change when already merged", func(t *testing.T) {
		firstMergeTime := time.Now()
		secondMergeTime := firstMergeTime.Add(time.Hour)

		pullRequest := NewPullRequest("pull-request-1", "feature_111", "Artem", time.Now().Add(-time.Hour))
		pullRequest.Merge(firstMergeTime)

		initialMergedAt := pullRequest.MergedAt

		pullRequest.Merge(secondMergeTime)

		assert.Equal(t, StatusMerged, pullRequest.Status)
		assert.Equal(t, initialMergedAt, pullRequest.MergedAt)
		assert.Equal(t, firstMergeTime, *pullRequest.MergedAt)
	})
}

func TestPullRequest_IsReviewer(t *testing.T) {
	pullRequest := NewPullRequest("pull-request-1", "feature_111", "Artem", time.Now())
	pullRequest.AddReviewers([]value_objects.UserID{"user1", "user2"})

	assert.True(t, pullRequest.IsReviewer("user1"))
	assert.True(t, pullRequest.IsReviewer("user2"))
	assert.False(t, pullRequest.IsReviewer("user3"))
	assert.False(t, pullRequest.IsReviewer("Artem"))
}

func TestPullRequest_IsMerged(t *testing.T) {
	pullRequest := NewPullRequest("pull-request-1", "feature_111", "Artem", time.Now())

	assert.False(t, pullRequest.IsMerged())

	pullRequest.Merge(time.Now())

	assert.True(t, pullRequest.IsMerged())
}

func TestPullRequest_IsFullyAssigned(t *testing.T) {
	pullRequest := NewPullRequest("pull-request-1", "feature_111", "Artem", time.Now())

	assert.False(t, pullRequest.IsFullyAssigned())

	pullRequest.AddReviewers([]value_objects.UserID{"user1"})
	assert.False(t, pullRequest.IsFullyAssigned())

	pullRequest.AddReviewers([]value_objects.UserID{"user2"})
	assert.True(t, pullRequest.IsFullyAssigned())
}

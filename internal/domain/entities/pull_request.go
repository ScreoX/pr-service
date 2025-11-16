package entities

import (
	"time"

	"pr-service/internal/domain"
	"pr-service/internal/domain/value_objects"
)

type PullRequestStatus string

const (
	StatusOpen   PullRequestStatus = "OPEN"
	StatusMerged PullRequestStatus = "MERGED"
)

const maxReviewers = 2

type PullRequest struct {
	ID       value_objects.PullRequestID
	Name     string
	AuthorID value_objects.UserID
	Status   PullRequestStatus

	reviewers []value_objects.UserID

	CreatedAt time.Time
	MergedAt  *time.Time
}

func NewPullRequest(id value_objects.PullRequestID, name string, authorID value_objects.UserID, createdAt time.Time) *PullRequest {
	return &PullRequest{
		ID:        id,
		Name:      name,
		AuthorID:  authorID,
		Status:    StatusOpen,
		reviewers: make([]value_objects.UserID, 0, maxReviewers),
		CreatedAt: createdAt,
	}
}

func (pr *PullRequest) Reviewers() []value_objects.UserID {
	reviewersCopy := make([]value_objects.UserID, len(pr.reviewers))
	copy(reviewersCopy, pr.reviewers)

	return reviewersCopy
}

func (pr *PullRequest) IsReviewer(id value_objects.UserID) bool {
	for _, reviewerID := range pr.reviewers {
		if reviewerID == id {
			return true
		}
	}

	return false
}

func (pr *PullRequest) IsMerged() bool {
	return pr.Status == StatusMerged
}

func (pr *PullRequest) IsFullyAssigned() bool {
	return len(pr.reviewers) == maxReviewers
}

func (pr *PullRequest) AddReviewers(candidates []value_objects.UserID) []value_objects.UserID {
	if pr.IsMerged() || pr.IsFullyAssigned() || len(candidates) == 0 {
		return nil
	}

	reviewersLimit := maxReviewers - len(pr.reviewers)
	addedReviewers := make([]value_objects.UserID, 0, reviewersLimit)

	for _, candidate := range candidates {
		if len(addedReviewers) == reviewersLimit {
			break
		}
		if candidate == pr.AuthorID || pr.IsReviewer(candidate) {
			continue
		}

		pr.reviewers = append(pr.reviewers, candidate)
		addedReviewers = append(addedReviewers, candidate)
	}

	if len(addedReviewers) == 0 {
		return nil
	}

	return addedReviewers
}

func (pr *PullRequest) SetReviewers(reviewers []value_objects.UserID) {
	pr.reviewers = reviewers
}

func (pr *PullRequest) ReassignReviewer(oldID, newID value_objects.UserID) error {
	if pr.IsMerged() {
		return domain.ErrPRMerged
	}
	if !pr.IsReviewer(oldID) {
		return domain.ErrNotAssigned
	}
	if newID == pr.AuthorID {
		return domain.ErrNoCandidate
	}
	for _, reviewerID := range pr.reviewers {
		if reviewerID == newID {
			return domain.ErrNoCandidate
		}
	}

	for i, reviewerID := range pr.reviewers {
		if reviewerID == oldID {
			pr.reviewers[i] = newID

			return nil
		}
	}

	return domain.ErrNotAssigned
}

func (pr *PullRequest) Merge(mergedAt time.Time) {
	if pr.Status == StatusMerged {
		return
	}

	pr.Status = StatusMerged
	pr.MergedAt = &mergedAt
}

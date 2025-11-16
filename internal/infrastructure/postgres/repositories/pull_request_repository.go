package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"

	"pr-service/internal/app"
	"pr-service/internal/domain"
	"pr-service/internal/domain/entities"
	"pr-service/internal/domain/value_objects"
	"pr-service/internal/infrastructure/db_mappers"
	"pr-service/internal/infrastructure/db_models"
)

type pullRequestRepository struct {
	db *sql.DB
	sb squirrel.StatementBuilderType
}

func NewPullRequestRepository(db *sql.DB) app.PullRequestRepository {
	return &pullRequestRepository{
		db: db,
		sb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *pullRequestRepository) Create(ctx context.Context, pullRequest *entities.PullRequest) error {
	dbPullRequest := db_mappers.ToPullRequestDBModel(*pullRequest)

	query, args, err := r.sb.Insert("pull_requests").
		Columns("id", "pull_request_name", "author_id", "status", "created_at", "merged_at").
		Values(dbPullRequest.ID, dbPullRequest.Name, dbPullRequest.AuthorID, dbPullRequest.Status, dbPullRequest.CreatedAt, dbPullRequest.MergedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build insert query: %v", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to insert pull request: %v", err)
	}

	if len(pullRequest.Reviewers()) > 0 {
		for _, reviewerID := range pullRequest.Reviewers() {
			reviewerQuery, reviewerArgs, err := r.sb.Insert("pull_request_reviewers").
				Columns("pull_request_id", "user_id").
				Values(dbPullRequest.ID, reviewerID).
				ToSql()
			if err != nil {
				return fmt.Errorf("failed to build insert query for reviewers: %v", err)
			}

			_, err = r.db.ExecContext(ctx, reviewerQuery, reviewerArgs...)
			if err != nil {
				return fmt.Errorf("failed to insert reviewer: %v", err)
			}
		}
	}

	return nil
}

func (r *pullRequestRepository) Save(ctx context.Context, pullRequest *entities.PullRequest) error {
	dbPullRequest := db_mappers.ToPullRequestDBModel(*pullRequest)

	query, args, err := r.sb.Update("pull_requests").
		Set("status", dbPullRequest.Status).
		Set("merged_at", dbPullRequest.MergedAt).
		Where(squirrel.Eq{"id": dbPullRequest.ID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update query: %v", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update pull request: %v", err)
	}

	return nil
}

func (r *pullRequestRepository) GetByID(ctx context.Context, id value_objects.PullRequestID) (*entities.PullRequest, error) {
	var dbPullRequest db_models.PullRequest

	query, args, err := r.sb.Select("id", "pull_request_name", "author_id", "status", "created_at", "merged_at").
		From("pull_requests").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %v", err)
	}

	err = r.db.QueryRowContext(ctx, query, args...).Scan(&dbPullRequest.ID, &dbPullRequest.Name, &dbPullRequest.AuthorID, &dbPullRequest.Status, &dbPullRequest.CreatedAt, &dbPullRequest.MergedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrPRNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pull request: %v", err)
	}

	pullRequest := db_mappers.FromPullRequestDBModel(dbPullRequest)

	var reviewers []value_objects.UserID
	reviewersQuery, reviewersArgs, err := r.sb.Select("user_id").
		From("pull_request_reviewers").
		Where(squirrel.Eq{"pull_request_id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build reviewers query: %v", err)
	}

	rows, err := r.db.QueryContext(ctx, reviewersQuery, reviewersArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch reviewers: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var reviewerID value_objects.UserID
		if err := rows.Scan(&reviewerID); err != nil {
			return nil, fmt.Errorf("failed to scan reviewer: %v", err)
		}

		reviewers = append(reviewers, reviewerID)
	}

	pullRequest.SetReviewers(reviewers)

	return &pullRequest, nil
}

func (r *pullRequestRepository) GetByReviewer(ctx context.Context, reviewerID value_objects.UserID) ([]entities.PullRequest, error) {
	var pullRequests []entities.PullRequest

	query, args, err := r.sb.Select("pr.id", "pr.pull_request_name", "pr.author_id", "pr.status", "pr.created_at", "pr.merged_at").
		From("pull_requests AS pr").
		Join("pull_request_reviewers AS prr ON pr.id = prr.pull_request_id").
		Where(squirrel.Eq{"prr.user_id": reviewerID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %v", err)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pull requests: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var dbPullRequest db_models.PullRequest

		if err := rows.Scan(&dbPullRequest.ID, &dbPullRequest.Name, &dbPullRequest.AuthorID, &dbPullRequest.Status, &dbPullRequest.CreatedAt, &dbPullRequest.MergedAt); err != nil {
			return nil, fmt.Errorf("failed to scan pull request: %v", err)
		}

		pullRequest := db_mappers.FromPullRequestDBModel(dbPullRequest)

		pullRequests = append(pullRequests, pullRequest)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %v", err)
	}

	return pullRequests, nil
}

func (r *pullRequestRepository) GetAll(ctx context.Context) ([]entities.PullRequest, error) {
	query, args, err := r.sb.Select("id", "pull_request_name", "author_id", "status", "created_at", "merged_at").
		From("pull_requests").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %v", err)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pull requests: %v", err)
	}
	defer rows.Close()

	var pullRequests []entities.PullRequest

	for rows.Next() {
		var dbPullRequest db_models.PullRequest
		if err := rows.Scan(&dbPullRequest.ID, &dbPullRequest.Name, &dbPullRequest.AuthorID, &dbPullRequest.Status, &dbPullRequest.CreatedAt, &dbPullRequest.MergedAt); err != nil {
			return nil, fmt.Errorf("failed to scan pull request: %v", err)
		}

		pullRequests = append(pullRequests, db_mappers.FromPullRequestDBModel(dbPullRequest))
	}

	return pullRequests, nil
}

func (r *pullRequestRepository) ReassignReviewer(ctx context.Context, pullRequestID value_objects.PullRequestID, oldReviewerID value_objects.UserID, newReviewerID value_objects.UserID) error {
	query, args, err := r.sb.Update("pull_request_reviewers").
		Set("user_id", newReviewerID).
		Where(squirrel.Eq{"pull_request_id": pullRequestID}).
		Where(squirrel.Eq{"user_id": oldReviewerID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update query: %v", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to reassign reviewer: %v", err)
	}

	return nil
}

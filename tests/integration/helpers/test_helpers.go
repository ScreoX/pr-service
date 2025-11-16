package helpers

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/lib/pq"

	"pr-service/config"
	"pr-service/internal/infrastructure/db"
)

func SetupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	cfg := config.LoadTest()

	testDB, err := db.Init(cfg)
	if err != nil {
		t.Fatalf("Failed to connect to test DB: %v", err)
	}

	return testDB
}

func CleanupTestDB(t *testing.T, db *sql.DB) {
	t.Helper()

	if db != nil {
		tables := []string{
			"pull_request_reviewers",
			"pull_requests",
			"users",
			"teams",
		}

		for _, table := range tables {
			_, err := db.Exec(fmt.Sprintf("DELETE FROM %s", table))
			if err != nil {
				t.Logf("Failed to clean table %s: %v", table, err)
			}
		}
	}
}

func InsertTestUser(db *sql.DB, id, username, teamName string, isActive bool) error {
	_, err := db.Exec(
		"INSERT INTO users (id, username, team_name, is_active) VALUES ($1, $2, $3, $4)",
		id, username, teamName, isActive,
	)

	return err
}

func GetUserActivity(db *sql.DB, userID string) (bool, error) {
	var isActive bool

	err := db.QueryRow("SELECT is_active FROM users WHERE id = $1", userID).Scan(&isActive)

	return isActive, err
}

func InsertTestTeam(db *sql.DB, id, teamName string) error {
	_, err := db.Exec(
		"INSERT INTO teams (id, team_name) VALUES ($1, $2)",
		id, teamName,
	)

	return err
}

func TeamExists(db *sql.DB, teamName string) (bool, error) {
	var exists bool

	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM teams WHERE team_name = $1)", teamName).Scan(&exists)

	return exists, err
}

func InsertTestPullRequest(db *sql.DB, id, name, authorID, status string) error {
	_, err := db.Exec(
		"INSERT INTO pull_requests (id, pull_request_name, author_id, status) VALUES ($1, $2, $3, $4)",
		id, name, authorID, status,
	)

	return err
}

func AddReviewerToPullRequest(db *sql.DB, pullRequestID, userID string) error {
	_, err := db.Exec(
		"INSERT INTO pull_request_reviewers (pull_request_id, user_id) VALUES ($1, $2)",
		pullRequestID, userID,
	)

	return err
}

func PullRequestExists(db *sql.DB, pullRequestID string) (bool, error) {
	var exists bool

	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM pull_requests WHERE id = $1)", pullRequestID).Scan(&exists)

	return exists, err
}

func GetPullRequestStatus(db *sql.DB, pullRequestID string) (string, error) {
	var status string

	err := db.QueryRow("SELECT status FROM pull_requests WHERE id = $1", pullRequestID).Scan(&status)

	return status, err
}

func GetPullRequestReviewers(db *sql.DB, pullRequestID string) ([]string, error) {
	var reviewers []string

	rows, err := db.Query("SELECT user_id FROM pull_request_reviewers WHERE pull_request_id = $1", pullRequestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var reviewer string
		if err := rows.Scan(&reviewer); err != nil {
			return nil, err
		}
		reviewers = append(reviewers, reviewer)
	}

	return reviewers, nil
}

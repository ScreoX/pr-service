package domain

import "errors"

var (
	ErrTeamExists      = errors.New("TEAM_EXISTS")
	ErrPRExists        = errors.New("PR_EXISTS")
	ErrPRMerged        = errors.New("PR_MERGED")
	ErrNotAssigned     = errors.New("NOT_ASSIGNED")
	ErrNoCandidate     = errors.New("NO_CANDIDATE")
	ErrUserNotFound    = errors.New("USER_NOT_FOUND")
	ErrTeamNotFound    = errors.New("TEAM_NOT_FOUND")
	ErrPRNotFound      = errors.New("PR_NOT_FOUND")
	ErrAuthorNotActive = errors.New("AUTHOR_NOT_ACTIVE")
)

package apierrors

const (
	InvalidRequestBody = "INVALID_REQUEST_BODY"
	DuplicateUserIDs   = "DUPLICATE_USER_IDS"
	MissingUserID      = "MISSING_USER_ID"
	MissingTeamName    = "MISSING_TEAM_NAME"
	PRExists           = "PR_EXISTS"
	TeamExists         = "TEAM_EXISTS"
	PRMerged           = "PR_MERGED"
	NoCandidate        = "NO_CANDIDATE"
	NotAssigned        = "NOT_ASSIGNED"
	AuthorNotActive    = "AUTHOR_NOT_ACTIVE"
	NotFound           = "NOT_FOUND"
	InternalError      = "INTERNAL_ERROR"
)

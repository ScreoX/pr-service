package apierrors

const (
	InvalidRequestBodyMessage = "invalid request body"
	DuplicateUserIDsMessage   = "team contains duplicate user_ids"
	MissingUserIDMessage      = "user ID is required"
	MissingTeamNameMessage    = "team name is required"
	PRExistsMessage           = "PR id already exists"
	TeamExistsMessage         = "team_name already exists"
	PRMergedMessage           = "cannot reassign on merged PR"
	NoCandidateMessage        = "no active replacement candidate in team"
	NotAssignedMessage        = "reviewer is not assigned to this PR"
	AuthorNotActiveMessage    = "user can not create PR with false active status"
	NotFoundMessage           = "resource not found"
	InternalErrorMessage      = "internal server error"
)

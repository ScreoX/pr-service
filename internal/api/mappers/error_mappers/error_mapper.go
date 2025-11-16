package error_mappers

import (
	"errors"
	"net/http"

	"pr-service/internal/api/apierrors"
	"pr-service/internal/api/dto"
	"pr-service/internal/domain"
)

func ToHTTPError(domainErr error) (int, dto.ErrorResponse) {
	switch {
	case errors.Is(domainErr, domain.ErrPRExists):
		return http.StatusConflict, dto.ErrorResponse{
			Error: dto.Error{
				Code:    apierrors.PRExists,
				Message: apierrors.PRExistsMessage,
			},
		}

	case errors.Is(domainErr, domain.ErrTeamExists):
		return http.StatusConflict, dto.ErrorResponse{
			Error: dto.Error{
				Code:    apierrors.TeamExists,
				Message: apierrors.TeamExistsMessage,
			},
		}

	case errors.Is(domainErr, domain.ErrPRMerged):
		return http.StatusConflict, dto.ErrorResponse{
			Error: dto.Error{
				Code:    apierrors.PRMerged,
				Message: apierrors.PRMergedMessage,
			},
		}

	case errors.Is(domainErr, domain.ErrNoCandidate):
		return http.StatusConflict, dto.ErrorResponse{
			Error: dto.Error{
				Code:    apierrors.NoCandidate,
				Message: apierrors.NoCandidateMessage,
			},
		}

	case errors.Is(domainErr, domain.ErrNotAssigned):
		return http.StatusConflict, dto.ErrorResponse{
			Error: dto.Error{
				Code:    apierrors.NotAssigned,
				Message: apierrors.NotAssignedMessage,
			},
		}

	case errors.Is(domainErr, domain.ErrAuthorNotActive):
		return http.StatusConflict, dto.ErrorResponse{
			Error: dto.Error{
				Code:    apierrors.AuthorNotActive,
				Message: apierrors.AuthorNotActiveMessage,
			},
		}

	case errors.Is(domainErr, domain.ErrUserNotFound),
		errors.Is(domainErr, domain.ErrTeamNotFound),
		errors.Is(domainErr, domain.ErrPRNotFound):
		return http.StatusNotFound, dto.ErrorResponse{
			Error: dto.Error{
				Code:    apierrors.NotFound,
				Message: apierrors.NotFoundMessage,
			},
		}

	default:
		return http.StatusInternalServerError, dto.ErrorResponse{
			Error: dto.Error{
				Code:    apierrors.InternalError,
				Message: apierrors.InternalErrorMessage,
			},
		}
	}
}

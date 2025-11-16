package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"pr-service/internal/api/apierrors"
	"pr-service/internal/api/dto"
	"pr-service/internal/api/mappers/dto_mappers"
	"pr-service/internal/api/mappers/error_mappers"
	"pr-service/internal/app/services"
	"pr-service/internal/domain/value_objects"
)

type PullRequestHandler struct {
	pullRequestService services.PullRequestService
}

func NewPullRequestHandler(pullRequestService services.PullRequestService) *PullRequestHandler {
	return &PullRequestHandler{
		pullRequestService: pullRequestService,
	}
}

func (h *PullRequestHandler) CreatePullRequest(c *gin.Context) {
	var request dto.CreatePullRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.Error{
				Code:    apierrors.InvalidRequestBody,
				Message: apierrors.InvalidRequestBodyMessage,
			},
		})
		return
	}

	pullRequestID, pullRequestName, authorID := dto_mappers.FromCreatePullRequestDTO(request)
	pullRequest, err := h.pullRequestService.Create(c, pullRequestID, pullRequestName, authorID)
	if err != nil {
		statusCode, errorResponse := error_mappers.ToHTTPError(err)
		c.JSON(statusCode, errorResponse)
		return
	}

	c.JSON(http.StatusCreated, dto_mappers.ToPullRequestResponseDTO(*pullRequest))
}

func (h *PullRequestHandler) MergePullRequest(c *gin.Context) {
	var request dto.MergePullRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.Error{
				Code:    apierrors.InvalidRequestBody,
				Message: apierrors.InvalidRequestBodyMessage,
			},
		})
		return
	}

	pullRequestID := value_objects.PullRequestID(request.PullRequestID)
	pullRequest, err := h.pullRequestService.Merge(c, pullRequestID)
	if err != nil {
		statusCode, errorResponse := error_mappers.ToHTTPError(err)
		c.JSON(statusCode, errorResponse)
		return
	}

	c.JSON(http.StatusOK, dto_mappers.ToPullRequestResponseDTO(*pullRequest))
}

func (h *PullRequestHandler) ReassignReviewer(c *gin.Context) {
	var request dto.ReassignReviewerRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.Error{
				Code:    apierrors.InvalidRequestBody,
				Message: apierrors.InvalidRequestBodyMessage,
			},
		})
		return
	}

	pullRequestID, oldReviewerID := dto_mappers.FromReassignReviewerRequestDTO(request)
	pullRequest, newReviewerID, err := h.pullRequestService.ReassignReviewer(c, pullRequestID, oldReviewerID)
	if err != nil {
		statusCode, errorResponse := error_mappers.ToHTTPError(err)
		c.JSON(statusCode, errorResponse)
		return
	}

	c.JSON(http.StatusOK, dto_mappers.ToPullRequestReassignResponseDTO(*pullRequest, newReviewerID))
}

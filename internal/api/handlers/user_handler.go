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

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) SetActiveStatus(c *gin.Context) {
	var request dto.UserStatusRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.Error{
				Code:    apierrors.InvalidRequestBody,
				Message: apierrors.InvalidRequestBodyMessage,
			},
		})
		return
	}

	userID := value_objects.UserID(request.UserID)
	updatedUser, err := h.userService.SetActiveStatus(c, userID, request.IsActive)
	if err != nil {
		statusCode, errorResponse := error_mappers.ToHTTPError(err)
		c.JSON(statusCode, errorResponse)
		return
	}

	c.JSON(http.StatusOK, dto_mappers.ToUserStatusResponseDTO(updatedUser))
}

func (h *UserHandler) GetUserReviews(c *gin.Context) {
	userID := c.DefaultQuery("user_id", "")
	if userID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.Error{
				Code:    apierrors.MissingUserID,
				Message: apierrors.MissingUserIDMessage,
			},
		})
		return
	}

	parsedUserID := value_objects.UserID(userID)
	pullRequests, err := h.userService.GetUserReviews(c, parsedUserID)
	if err != nil {
		statusCode, errorResponse := error_mappers.ToHTTPError(err)
		c.JSON(statusCode, errorResponse)
		return
	}

	c.JSON(http.StatusOK, dto_mappers.ToUserReviewsResponseDTO(parsedUserID, pullRequests))
}

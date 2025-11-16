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

type TeamHandler struct {
	teamService services.TeamService
}

func NewTeamHandler(teamService services.TeamService) *TeamHandler {
	return &TeamHandler{
		teamService: teamService,
	}
}

func (h *TeamHandler) CreateTeam(c *gin.Context) {
	var request dto.CreateTeamRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.Error{
				Code:    apierrors.InvalidRequestBody,
				Message: apierrors.InvalidRequestBodyMessage,
			},
		})
		return
	}

	if hasDuplicateUserIDs(request.Members) {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.Error{
				Code:    apierrors.DuplicateUserIDs,
				Message: apierrors.DuplicateUserIDsMessage,
			},
		})
		return
	}

	teamName, members := dto_mappers.FromCreateTeamRequestDTO(request)
	team, members, err := h.teamService.Create(c, teamName, members)
	if err != nil {
		statusCode, errorResponse := error_mappers.ToHTTPError(err)
		c.JSON(statusCode, errorResponse)
		return
	}

	c.JSON(http.StatusCreated, dto_mappers.ToTeamResponseDTO(team, members))
}

func (h *TeamHandler) GetTeam(c *gin.Context) {
	teamName := c.DefaultQuery("team_name", "")
	if teamName == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.Error{
				Code:    apierrors.MissingTeamName,
				Message: apierrors.MissingTeamNameMessage,
			},
		})
		return
	}

	parsedTeamName := value_objects.TeamName(teamName)
	team, members, err := h.teamService.GetByName(c, parsedTeamName)
	if err != nil {
		statusCode, errorResponse := error_mappers.ToHTTPError(err)
		c.JSON(statusCode, errorResponse)
		return
	}

	c.JSON(http.StatusOK, dto_mappers.ToTeamResponseDTO(team, members))
}

func hasDuplicateUserIDs(members []dto.TeamMember) bool {
	seen := make(map[string]bool)

	for _, member := range members {
		if seen[member.UserID] {
			return true
		}
		seen[member.UserID] = true
	}

	return false
}

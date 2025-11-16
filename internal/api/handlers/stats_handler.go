package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"pr-service/internal/api/apierrors"
	"pr-service/internal/api/dto"
	"pr-service/internal/app/services"
)

type StatsHandler struct {
	statsService services.StatsService
}

func NewStatsHandler(statsService services.StatsService) *StatsHandler {
	return &StatsHandler{
		statsService: statsService,
	}
}

func (h *StatsHandler) GetStats(c *gin.Context) {
	stats, err := h.statsService.GetStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: dto.Error{
				Code:    apierrors.InternalError,
				Message: apierrors.InternalErrorMessage,
			},
		})
		return
	}

	c.JSON(http.StatusOK, stats)
}

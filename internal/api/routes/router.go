package routes

import (
	"github.com/gin-gonic/gin"

	"pr-service/internal/api/handlers"
	"pr-service/internal/api/middleware"
)

func Setup(userHandler *handlers.UserHandler, teamHandler *handlers.TeamHandler, pullRequestHandler *handlers.PullRequestHandler, statsHandler *handlers.StatsHandler) *gin.Engine {
	router := gin.Default()

	err := router.SetTrustedProxies(nil)
	if err != nil {
		return nil
	}

	router.Use(middleware.LoggerMiddleware())

	router.POST("/users/setIsActive", userHandler.SetActiveStatus)
	router.GET("/users/getReview", userHandler.GetUserReviews)

	router.POST("/team/add", teamHandler.CreateTeam)
	router.GET("/team/get", teamHandler.GetTeam)

	router.POST("/pullRequest/create", pullRequestHandler.CreatePullRequest)
	router.POST("/pullRequest/merge", pullRequestHandler.MergePullRequest)
	router.POST("/pullRequest/reassign", pullRequestHandler.ReassignReviewer)

	router.GET("/stats", statsHandler.GetStats)

	return router
}

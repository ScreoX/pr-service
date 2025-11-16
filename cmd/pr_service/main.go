package main

import (
	"log"

	"pr-service/config"
	"pr-service/internal/api/handlers"
	"pr-service/internal/api/routes"
	"pr-service/internal/app/services"
	"pr-service/internal/infrastructure/db"
	"pr-service/internal/infrastructure/postgres/repositories"
	"pr-service/internal/infrastructure/providers"
)

func main() {
	cfg := config.Load()

	database, err := db.Init(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	txManager := db.NewTxManager(database)
	timeProvider := providers.NewCurrentTime()
	randomProvider := providers.NewRealRandom()

	userRepository := repositories.NewUserRepository(database)
	teamRepository := repositories.NewTeamRepository(database)
	pullRequestRepository := repositories.NewPullRequestRepository(database)

	userService := services.NewUserService(userRepository, pullRequestRepository)
	teamService := services.NewTeamService(userRepository, teamRepository, txManager)
	pullRequestService := services.NewPullRequestService(userRepository, teamRepository, pullRequestRepository, txManager, timeProvider, randomProvider)
	statsService := services.NewStatsService(userRepository, teamRepository, pullRequestRepository)

	userHandler := handlers.NewUserHandler(userService)
	teamHandler := handlers.NewTeamHandler(teamService)
	pullRequestHandler := handlers.NewPullRequestHandler(pullRequestService)
	statsHandler := handlers.NewStatsHandler(statsService)

	router := routes.Setup(userHandler, teamHandler, pullRequestHandler, statsHandler)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

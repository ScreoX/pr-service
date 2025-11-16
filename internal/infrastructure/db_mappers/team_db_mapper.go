package db_mappers

import (
	"pr-service/internal/domain/entities"
	"pr-service/internal/domain/value_objects"
	"pr-service/internal/infrastructure/db_models"
)

func FromTeamDBModel(dbTeam db_models.Team) entities.Team {
	return entities.Team{
		Name: value_objects.TeamName(dbTeam.Name),
	}
}

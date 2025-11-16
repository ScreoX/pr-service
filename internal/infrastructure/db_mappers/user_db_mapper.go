package db_mappers

import (
	"pr-service/internal/domain/entities"
	"pr-service/internal/domain/value_objects"
	"pr-service/internal/infrastructure/db_models"
)

func ToUserDBModel(user entities.User) db_models.User {
	return db_models.User{
		ID:       string(user.ID),
		Username: user.Username,
		Team:     string(user.Team),
		IsActive: user.IsActive,
	}
}

func FromUserDBModel(dbUser db_models.User) entities.User {
	return entities.User{
		ID:       value_objects.UserID(dbUser.ID),
		Username: dbUser.Username,
		Team:     value_objects.TeamName(dbUser.Team),
		IsActive: dbUser.IsActive,
	}
}

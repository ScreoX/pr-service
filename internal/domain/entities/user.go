package entities

import (
	"pr-service/internal/domain/value_objects"
)

type User struct {
	ID       value_objects.UserID
	Username string
	Team     value_objects.TeamName
	IsActive bool
}

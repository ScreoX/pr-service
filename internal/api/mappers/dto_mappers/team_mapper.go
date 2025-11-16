package dto_mappers

import (
	"pr-service/internal/api/dto"
	"pr-service/internal/domain/entities"
	"pr-service/internal/domain/value_objects"
)

func FromCreateTeamRequestDTO(dto dto.CreateTeamRequest) (value_objects.TeamName, []entities.User) {
	teamName := value_objects.TeamName(dto.TeamName)
	users := make([]entities.User, len(dto.Members))

	for i, member := range dto.Members {
		users[i] = entities.User{
			ID:       value_objects.UserID(member.UserID),
			Username: member.Username,
			IsActive: member.IsActive,
			Team:     teamName,
		}
	}

	return teamName, users
}

func ToTeamResponseDTO(team entities.Team, members []entities.User) dto.TeamResponse {
	var memberDTOs []dto.TeamMember

	for _, member := range members {
		memberDTOs = append(memberDTOs, dto.TeamMember{
			UserID:   string(member.ID),
			Username: member.Username,
			IsActive: member.IsActive,
		})
	}

	return dto.TeamResponse{
		TeamName: string(team.Name),
		Members:  memberDTOs,
	}
}

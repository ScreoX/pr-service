package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"

	"pr-service/internal/app"
	"pr-service/internal/domain"
	"pr-service/internal/domain/entities"
	"pr-service/internal/domain/value_objects"
	"pr-service/internal/infrastructure/db_mappers"
	"pr-service/internal/infrastructure/db_models"
)

type teamRepository struct {
	db *sql.DB
	sb squirrel.StatementBuilderType
}

func NewTeamRepository(db *sql.DB) app.TeamRepository {
	return &teamRepository{
		db: db,
		sb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *teamRepository) Create(ctx context.Context, team entities.Team) error {
	query, args, err := r.sb.Insert("teams").
		Columns("id", "team_name").
		Values(string(team.Name), string(team.Name)).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build insert query: %v", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to create team: %v", err)
	}

	return nil
}

func (r *teamRepository) GetByName(ctx context.Context, name value_objects.TeamName) (entities.Team, error) {
	query, args, err := r.sb.Select("id", "team_name").
		From("teams").
		Where(squirrel.Eq{"team_name": string(name)}).
		ToSql()
	if err != nil {
		return entities.Team{}, fmt.Errorf("failed to build query: %v", err)
	}

	var dbTeam db_models.Team

	err = r.db.QueryRowContext(ctx, query, args...).Scan(&dbTeam.ID, &dbTeam.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.Team{}, domain.ErrTeamNotFound
		}
		return entities.Team{}, fmt.Errorf("failed to fetch team: %v", err)
	}

	return db_mappers.FromTeamDBModel(dbTeam), nil
}

func (r *teamRepository) GetAll(ctx context.Context) ([]entities.Team, error) {
	query, args, err := r.sb.Select("id", "team_name").
		From("teams").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %v", err)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch teams: %v", err)
	}
	defer rows.Close()

	var teams []entities.Team

	for rows.Next() {
		var dbTeam db_models.Team
		if err := rows.Scan(&dbTeam.ID, &dbTeam.Name); err != nil {
			return nil, fmt.Errorf("failed to scan team: %v", err)
		}

		teams = append(teams, db_mappers.FromTeamDBModel(dbTeam))
	}

	return teams, nil
}

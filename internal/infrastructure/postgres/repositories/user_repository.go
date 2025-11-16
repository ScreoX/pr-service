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

type userRepository struct {
	db *sql.DB
	sb squirrel.StatementBuilderType
}

func NewUserRepository(db *sql.DB) app.UserRepository {
	return &userRepository{
		db: db,
		sb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *userRepository) GetByID(ctx context.Context, id value_objects.UserID) (entities.User, error) {
	var dbUser db_models.User

	query, args, err := r.sb.Select("id", "username", "team_name", "is_active").
		From("users").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return entities.User{}, fmt.Errorf("failed to build query: %v", err)
	}

	err = r.db.QueryRowContext(ctx, query, args...).Scan(&dbUser.ID, &dbUser.Username, &dbUser.Team, &dbUser.IsActive)
	if errors.Is(err, sql.ErrNoRows) {
		return entities.User{}, domain.ErrUserNotFound
	}
	if err != nil {
		return entities.User{}, fmt.Errorf("failed to fetch user: %v", err)
	}

	return db_mappers.FromUserDBModel(dbUser), nil
}

func (r *userRepository) GetUsersByTeam(ctx context.Context, teamName value_objects.TeamName) ([]entities.User, error) {
	var dbUsers []db_models.User

	query, args, err := r.sb.Select("id", "username", "team_name", "is_active").
		From("users").
		Where(squirrel.Eq{"team_name": string(teamName)}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %v", err)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var dbUser db_models.User
		if err := rows.Scan(&dbUser.ID, &dbUser.Username, &dbUser.Team, &dbUser.IsActive); err != nil {
			return nil, fmt.Errorf("failed to scan user: %v", err)
		}
		dbUsers = append(dbUsers, dbUser)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %v", err)
	}

	var users []entities.User

	for _, dbUser := range dbUsers {
		users = append(users, db_mappers.FromUserDBModel(dbUser))
	}

	return users, nil
}

func (r *userRepository) GetAll(ctx context.Context) ([]entities.User, error) {
	query, args, err := r.sb.Select("id", "username", "team_name", "is_active").
		From("users").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %v", err)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %v", err)
	}
	defer rows.Close()

	var users []entities.User

	for rows.Next() {
		var dbUser db_models.User
		if err := rows.Scan(&dbUser.ID, &dbUser.Username, &dbUser.Team, &dbUser.IsActive); err != nil {
			return nil, fmt.Errorf("failed to scan user: %v", err)
		}

		users = append(users, db_mappers.FromUserDBModel(dbUser))
	}

	return users, nil
}

func (r *userRepository) UpsertMembers(ctx context.Context, teamName value_objects.TeamName, members []entities.User) error {
	for _, member := range members {
		dbUser := db_mappers.ToUserDBModel(member)

		query, args, err := r.sb.Insert("users").
			Columns("id", "username", "team_name", "is_active").
			Values(dbUser.ID, dbUser.Username, string(teamName), dbUser.IsActive).
			Suffix("ON CONFLICT (id) DO UPDATE SET username = EXCLUDED.username, team_name = EXCLUDED.team_name, is_active = EXCLUDED.is_active").
			ToSql()

		if err != nil {
			return fmt.Errorf("failed to build upsert query: %v", err)
		}

		_, err = r.db.ExecContext(ctx, query, args...)
		if err != nil {
			return fmt.Errorf("failed to execute upsert query: %v", err)
		}
	}

	return nil
}

func (r *userRepository) SetIsActive(ctx context.Context, id value_objects.UserID, isActive bool) (entities.User, error) {
	query, args, err := r.sb.Update("users").
		Set("is_active", isActive).
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return entities.User{}, fmt.Errorf("failed to build update query: %v", err)
	}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return entities.User{}, fmt.Errorf("failed to execute update: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return entities.User{}, fmt.Errorf("failed to get rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return entities.User{}, domain.ErrUserNotFound
	}

	return r.GetByID(ctx, id)
}

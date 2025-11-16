package db_models

type User struct {
	ID       string `db:"id"`
	Username string `db:"username"`
	Team     string `db:"team_name"`
	IsActive bool   `db:"is_active"`
}

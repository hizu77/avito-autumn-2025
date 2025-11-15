package dbmodel

type User struct {
	ID       string `db:"id"`
	Name     string `db:"name"`
	TeamName string `db:"team_name"`
	IsActive bool   `db:"is_active"`
}

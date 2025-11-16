package dbmodel

type Row struct {
	TName     string  `db:"team_name"`
	UID       *string `db:"user_id"`
	UName     *string `db:"user_name"`
	UIsActive *bool   `db:"user_is_active"`
}

package dbmodel

type Admin struct {
	ID           string `db:"id"`
	PasswordHash string `db:"password"`
}

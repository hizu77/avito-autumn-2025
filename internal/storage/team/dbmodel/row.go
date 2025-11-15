package dbmodel

type Row struct {
	TName     string `json:"team_name"`
	UID       string `json:"user_id"`
	UName     string `json:"user_name"`
	UIsActive bool   `json:"user_is_active"`
}

package dbmodel

type Row struct {
	TName string `json:"t.name"`

	UID       string `json:"u.id"`
	UName     string `json:"u.name"`
	UTeamName string `json:"u.team_name"`
	UIsActive bool   `json:"u.is_active"`
}

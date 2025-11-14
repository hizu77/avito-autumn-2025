package team

const (
	teamTableName = "teams"
	userTableName = "users"

	teamColumnName = "name"

	userColumnID       = "id"
	userColumnName     = "name"
	userColumnTeamName = "team_name"
	userColumnIsActive = "is_active"
)

var allTeamColumns = []string{
	teamColumnName,
}

var allUserColumns = []string{
	userColumnID,
	userColumnName,
	userColumnTeamName,
	userColumnIsActive,
}

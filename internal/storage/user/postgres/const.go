package user

const (
	tableName = "users"

	columnID       = "id"
	columnName     = "name"
	columnTeamName = "team_name"
	columnIsActive = "is_active"
)

var allColumns = []string{
	columnID,
	columnName,
	columnTeamName,
	columnIsActive,
}

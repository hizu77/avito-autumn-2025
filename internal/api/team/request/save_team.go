package request

type SaveTeam struct {
	Name    string       `json:"team_name"`
	Members []TeamMember `json:"members"`
}

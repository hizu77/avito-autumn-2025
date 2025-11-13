package request

type TeamMember struct {
	ID       string `json:"user_id"`
	Name     string `json:"username"`
	IsActive bool   `json:"is_active"`
}

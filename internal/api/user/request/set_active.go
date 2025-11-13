package request

type SetActive struct {
	ID       string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

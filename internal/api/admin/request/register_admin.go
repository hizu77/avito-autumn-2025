package request

type RegisterAdmin struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}

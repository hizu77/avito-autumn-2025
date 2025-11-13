package request

type LoginAdmin struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}

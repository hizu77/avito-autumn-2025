package httperr

type ErrorBody struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

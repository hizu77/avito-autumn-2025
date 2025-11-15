package common

import (
	"net/http"

	"github.com/go-chi/render"
)

type Error struct {
	Body ErrorBody `json:"error"`
}

func NewError(code ErrorCode, msg ...string) Error {
	message := code.DefaultMessage()
	if len(msg) > 0 {
		message = msg[0]
	}

	return Error{
		Body: ErrorBody{
			Code:    code,
			Message: message,
		},
	}
}

func WriteError(
	w http.ResponseWriter,
	r *http.Request,
	code ErrorCode,
	msg ...string,
) {
	errResp := NewError(code, msg...)

	render.Status(r, code.HTTPStatus())
	render.JSON(w, r, errResp)
}

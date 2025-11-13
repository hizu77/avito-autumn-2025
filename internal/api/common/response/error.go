package response

type Error struct {
	Body ErrorBody `json:"error"`
}

func NewError(code Code, message string) Error {
	err := ErrorBody{
		Code:    code,
		Message: message,
	}

	return Error{
		Body: err,
	}
}

func NewBadRequestError(message string) Error {
	return NewError(CodeBadRequest, message)
}

func NewInternalServerError() Error {
	return NewError(CodeInternal, "internal server error")
}

func NewTeamExistsError() Error {
	return NewError(CodeTeamExists, "team already exists")
}

func NewNotFoundError() Error {
	return NewError(CodeNotFound, "resource not found")
}

func NewPRExistsError() Error {
	return NewError(CodePrExists, "PR already exists")
}

func NewPRMergedError() Error {
	return NewError(CodePrMerged, "PR is merged")
}

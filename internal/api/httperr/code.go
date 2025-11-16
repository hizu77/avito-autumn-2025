package httperr

import "net/http"

type ErrorCode string

const (
	CodeBadRequest         ErrorCode = "BAD_REQUEST"
	CodeInternal           ErrorCode = "INTERNAL"
	CodeTeamExists         ErrorCode = "TEAM_EXISTS"
	CodeNotFound           ErrorCode = "NOT_FOUND"
	CodePrExists           ErrorCode = "PR_EXISTS"
	CodePrMerged           ErrorCode = "PR_MERGED"
	CodeInvalidCredentials ErrorCode = "INVALID_CREDENTIALS"
	CodeAdminExists        ErrorCode = "ADMIN_EXISTS"
	CodeUnauthorized       ErrorCode = "UNAUTHORIZED"
	CodeNotAssigned        ErrorCode = "NOT_ASSIGNED"
	CodeNoCandidate        ErrorCode = "NO_CANDIDATE"
)

func (c ErrorCode) HTTPStatus() int {
	switch c {
	case CodeBadRequest, CodeTeamExists:
		return http.StatusBadRequest
	case CodeNotFound, CodeNotAssigned:
		return http.StatusNotFound
	case CodeUnauthorized, CodeInvalidCredentials:
		return http.StatusUnauthorized
	case CodePrExists, CodePrMerged,
		CodeNoCandidate, CodeAdminExists:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

func (c ErrorCode) DefaultMessage() string {
	switch c {
	case CodeBadRequest:
		return "bad request"
	case CodeNotFound:
		return "resource not found"
	case CodeUnauthorized:
		return "invalid token"
	case CodeInvalidCredentials:
		return "invalid id or password"
	case CodeTeamExists:
		return "team already exists"
	case CodePrExists:
		return "pull request already exists"
	case CodePrMerged:
		return "cannot reassign on merged PR"
	case CodeAdminExists:
		return "admin already exists"
	case CodeNotAssigned:
		return "reviewer not assigned"
	case CodeNoCandidate:
		return "no candidate to reassign"
	default:
		return "internal server error"
	}
}

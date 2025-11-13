package model

import "github.com/pkg/errors"

var (
	ErrAdminAlreadyExists = errors.New("admin already exists")
	ErrAdminDoesNotExist  = errors.New("admin does not exist")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type Admin struct {
	ID           string
	PasswordHash string
}

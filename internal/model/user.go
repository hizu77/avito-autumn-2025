package model

import "github.com/pkg/errors"

var (
	ErrUserDoesNotExist = errors.New("user does not exist")
)

type User struct {
	ID       string
	Name     string
	TeamName string
	IsActive bool
}

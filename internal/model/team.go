package model

import "github.com/pkg/errors"

var (
	ErrTeamAlreadyExists = errors.New("team already exists")
	ErrTeamDoesNotExist  = errors.New("team does not exist")
)

type Team struct {
	Name    string
	Members []User
}

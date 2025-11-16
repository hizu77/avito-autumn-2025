package model

import "errors"

var (
	ErrUserDoesNotExist = errors.New("user does not exist")
)

type User struct {
	ID       string
	Name     string
	TeamName string
	IsActive bool
}

// I use this in collection.Map

func (u User) GetID() string {
	return u.ID
}

func (u User) GetName() string {
	return u.Name
}

func (u User) GetTeamName() string {
	return u.TeamName
}

func (u User) GetIsActive() bool {
	return u.IsActive
}

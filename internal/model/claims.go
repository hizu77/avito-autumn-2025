package model

import "github.com/golang-jwt/jwt/v4"

type Claims struct {
	AdminID string
	jwt.RegisteredClaims
}

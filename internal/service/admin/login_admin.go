package admin

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/hizu77/avito-autumn-2025/internal/model"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

const (
	jwtExpiration = time.Hour

	adminIDPayloadKey         = "admin_id"
	tokenExpirationPayloadKey = "exp"
)

func (s *Service) LoginAdmin(
	ctx context.Context,
	id string,
	password string,
) (string, error) {
	admin, err := s.storage.GetAdmin(ctx, id)
	if err != nil {
		return "", errors.Wrap(err, "getting admin")
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(admin.PasswordHash),
		[]byte(password),
	); err != nil {
		return "", model.ErrInvalidAdminPassword
	}

	claims := jwt.MapClaims{
		adminIDPayloadKey:         id,
		tokenExpirationPayloadKey: time.Now().Add(jwtExpiration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", errors.Wrap(err, "signing token")
	}

	return signedToken, nil
}

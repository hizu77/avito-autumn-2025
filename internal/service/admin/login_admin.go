package admin

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/hizu77/avito-autumn-2025/internal/model"
	"github.com/pkg/errors"
	"go.uber.org/zap"
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
		s.logger.Error("getting admin", zap.String("id", id), zap.Error(err))
		return "", errors.Wrap(err, "getting admin")
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(admin.PasswordHash),
		[]byte(password),
	); err != nil {
		s.logger.Error("invalid credentials", zap.String("id", id), zap.Error(err))
		return "", model.ErrInvalidAdminPassword
	}

	claims := jwt.MapClaims{
		adminIDPayloadKey:         id,
		tokenExpirationPayloadKey: time.Now().Add(jwtExpiration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(s.jwtSecret)
	if err != nil {
		s.logger.Error("signing token", zap.String("id", id), zap.Error(err))
		return "", errors.Wrap(err, "signing token")
	}

	return signedToken, nil
}

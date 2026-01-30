package jwt

import (
	"fmt"
	"time"

	"titan/common/logger"
	"titan/internal/auth-service/config"

	"github.com/golang-jwt/jwt/v5"
)

type AccessTokenClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type JWTManager struct {
	secret     []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
	issuer     string
	cfg        *config.Config
}

func New(
	cfg *config.Config) *JWTManager {
	return &JWTManager{
		secret:    []byte(cfg.JWT.Secret),
		accessTTL: cfg.JWT.AccessTLL,
		issuer:    cfg.JWT.Issuer,
		cfg:       cfg,
	}
}

func (j *JWTManager) GenerateToken(userID string, email string) (string, error) {
	now := time.Now()

	claims := AccessTokenClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.accessTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(j.secret)
	if err != nil {
		logger.Logger.Error("failed to sign JWT-token: %v", err)
		return "", fmt.Errorf("failed to sign JWT-token: %v", err)
	}

	return tokenString, nil
}

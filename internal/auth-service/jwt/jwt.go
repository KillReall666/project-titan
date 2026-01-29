package jwt

import (
	"time"

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
}

func NewJWTManager(secret string, accessTTL, refreshTTL time.Duration, issuer string) *JWTManager {
	return &JWTManager{
		secret:     []byte(secret),
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
		issuer:     issuer,
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
	return token.SignedString(j.secret)
}

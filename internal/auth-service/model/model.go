package model

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"createdAt"`
	Password     string    `json:"password"`
}

type RegisterRequest struct {
	UserMail string `json:"user_mail" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	UserMail string `json:"user_mail" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type ValidateRequest struct {
	Token string `json:"token" binding:"required"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

package service

import (
	"context"
	"errors"
	"time"

	"titan/internal/auth-service/config"
	"titan/internal/auth-service/model"
	"titan/internal/auth-service/storage"
	"titan/pkg/logger"

	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	cfg config.Config
	db  storage.AuthRepository
}

func NewAuthService(cfg config.Config, db storage.AuthRepository) *authService {
	return &authService{
		cfg: cfg,
		db:  db,
	}
}

// HashPassword — bcrypt hash
func (a *authService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword - verify pass
func (a *authService) CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (a *authService) SetUser(ctx context.Context, user model.RegisterRequest) error {
	hash, err := a.HashPassword(user.Password)
	if err != nil {
		logger.Logger.Error("err hash password", err)
		return errors.New("auth service error") //Пока такая заглушка, подумать что отдавать пользователю.
	}

	newUser := model.User{
		UserName:     user.Username,
		PasswordHash: hash,
		CreatedAt:    time.Now(),
	}

	err = a.db.SetUser(ctx, newUser)
	if err != nil {
		return err
	}

	return nil
}

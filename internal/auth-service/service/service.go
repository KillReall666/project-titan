package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"titan/common/logger"
	"titan/internal/auth-service/config"
	"titan/internal/auth-service/jwt"
	"titan/internal/auth-service/model"
	"titan/internal/auth-service/storage"

	jwt2 "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	cfg        config.Config
	db         storage.AuthRepository
	jwtManager jwt.JWTManager
}

func New(cfg config.Config, db storage.AuthRepository, jwt jwt.JWTManager) *authService {
	return &authService{
		cfg:        cfg,
		db:         db,
		jwtManager: jwt,
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
		ID:           uuid.New().String(),
		Email:        user.UserMail,
		PasswordHash: hash,
		CreatedAt:    time.Now(),
	}

	err = a.db.SetUser(ctx, newUser)
	if err != nil {
		return err
	}
	logger.Logger.Info("new user set success")

	return nil
}

func (a *authService) GetUser(ctx context.Context, username model.LoginRequest) (string, error) {
	user, err := a.db.GetUser(ctx, username.UserMail)
	if err != nil {
		return "", err //подумать над ошибкой
	}

	if !a.CheckPassword(username.Password, user.PasswordHash) {
		return "", errors.New("invalid username or password")
	}

	accessToken, err := a.jwtManager.GenerateToken(user.ID, user.Email)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (a *authService) Validate(token string) error {
	t, err := jwt2.ParseWithClaims(token, &jwt2.MapClaims{}, func(token *jwt2.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt2.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(a.cfg.JWT.Secret), nil
	})

	if err != nil {
		logger.Logger.Warn("invalid token attempt", err)
		return err
	}

	if t == nil || !t.Valid {
		return errors.New("token validation failed")
	}

	return nil
}

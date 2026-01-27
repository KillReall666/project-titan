package service

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"titan/internal/auth-service/model"
)

type authorizationService struct {
	cfg ?
	db 
}

// HashPassword — bcrypt hash
func (a *authorizationService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword - verify pass
func (a *authorizationService) CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// RegisterHandler - POST /register
func RegisterHandler(c *gin.Context) {
	var req model.RegisterRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hash, err := a.HashPassword(req.Password)
}

package register

import (
	"net/http"

	"titan/internal/auth-service/model"

	"github.com/gin-gonic/gin"
)

type registrator interface {
	// сюда интерфейс
}

type registrationHandler struct {
	//сюда интерфейс
}

func NewRegistrationHandler(register registrator) {
	return &registrationHandler{}
}

// RegisterHandler - POST /register
func RegistrationHandler(c *gin.Context) {
	var req model.RegisterRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hash, err := a.HashPassword(req.Password)
}

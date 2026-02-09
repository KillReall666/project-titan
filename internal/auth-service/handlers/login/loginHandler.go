package login

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"titan/internal/auth-service/model"
)

type loginer interface {
	GetUser(ctx context.Context, user model.LoginRequest) (string, error)
}

type loginHandler struct {
	loginer loginer
}

func NewLoginHandler(loginer loginer) *loginHandler {
	return &loginHandler{loginer: loginer}
}

// LoginHandler - POST /sign in
func (l *loginHandler) LoginHandler(c *gin.Context) {
	var req model.LoginRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := l.loginer.GetUser(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	resp := model.TokenResponse{Token: token}

	c.JSON(http.StatusOK, resp)
}

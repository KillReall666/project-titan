package register

import (
	"context"
	"net/http"

	"titan/common/utils/validation"
	"titan/internal/auth-service/model"

	"github.com/gin-gonic/gin"
)

type registrator interface {
	SetUser(ctx context.Context, user model.RegisterRequest) error
}

type registrationHandler struct {
	registrator registrator
}

func NewRegistrationHandler(register registrator) *registrationHandler {
	return &registrationHandler{registrator: register}
}

// RegistrationHandler - POST /sign up
func (r *registrationHandler) RegistrationHandler(c *gin.Context) {
	var req model.RegisterRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Можно будет в дальнейшем делать вайт-листы на домены или сами адреса.
	valid, err := validation.IsValidEmail(req.UserMail, false)
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = r.registrator.SetUser(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User successfully registered",
		"email":   req.UserMail},
	)

}

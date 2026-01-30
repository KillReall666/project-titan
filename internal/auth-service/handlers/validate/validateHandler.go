package validate

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type validator interface {
	Validate(token string) error
}

type validateHandler struct {
	validator validator
}

func NewValidationHandler(validator validator) *validateHandler {
	return &validateHandler{validator: validator}
}

func (v *validateHandler) ValidateHandler(c *gin.Context) {
	var req struct {
		Token string `json:"token" binding:"required"`
	}

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := v.validator.Validate(req.Token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid": true,
		//"username": claims.Username,
		//"expires":  claims.ExpiresAt.Time,
	})
}

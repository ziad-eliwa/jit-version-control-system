package api

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ziad-eliwa/jit-version-control-system/internal/middleware"
	"github.com/ziad-eliwa/jit-version-control-system/internal/models"
	"github.com/ziad-eliwa/jit-version-control-system/internal/services"
	"github.com/ziad-eliwa/jit-version-control-system/internal/utils"
)


type AuthHandler struct {
	AuthenticatonService services.AuthService
	TokenGenerator *middleware.JWTGenerator
	Logger *slog.Logger	
}
func (ah *AuthHandler) HandleRegister(c *gin.Context) {
	var req models.RegisterRequest		

	err := c.BindJSON(&req)

	if err != nil {
		ah.Logger.Error("Error converting the request from json")
		return
	}
	if req.Username == "" || req.Password == "" || req.EmailAddress == "" {
		c.JSON(http.StatusBadRequest, utils.Envelope{"error":"Missing Credentials"})
		ah.Logger.Error("Invalid Request Body")
		return
	}
	tokenRes, err := ah.AuthenticatonService.Register(req.Username,req.Password,req.FullName,req.EmailAddress)

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Envelope{"error":err})
		ah.Logger.Error("Error Registering User")
		return
	}

	c.JSON(http.StatusCreated,tokenRes)
}

func (ah *AuthHandler) HandleLogin(c *gin.Context) {

}

package api

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ziad-eliwa/jit-version-control-system/internal/middleware"
	"github.com/ziad-eliwa/jit-version-control-system/internal/models"
	"github.com/ziad-eliwa/jit-version-control-system/internal/services"
)

type AuthHandler struct {
	AuthenticatonService services.AuthService
	Logger               *slog.Logger
}

func (ah *AuthHandler) HandleRegister(c *gin.Context) {
	var req models.RegisterRequest

	err := c.BindJSON(&req)

	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid JSON Format")
		ah.Logger.Error("Error converting the request from json")
		return
	}

	if req.Username == "" || req.Password == "" || req.EmailAddress == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing Credentials"})
		ah.Logger.Error("Invalid Request Body")
		return
	}
	tokenRes, err := ah.AuthenticatonService.Register(req.Username, req.Password, req.FullName, req.EmailAddress)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		ah.Logger.Error(fmt.Sprintf("Error Registering User, %v", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, tokenRes)
}

func (ah *AuthHandler) HandleLogin(c *gin.Context) {
	var req models.LoginRequest
	err := c.BindJSON(&req)

	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid JSON Format")
		ah.Logger.Error("Error Converting the request from JSON")
		return
	}

	if req.Username == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing Credentials"})
		ah.Logger.Error("Invalid Request Body")
		return
	}

	tokenRes, err := ah.AuthenticatonService.Login(req.Username, req.Password)

	if err != nil {
		if err == services.ErrIncorrectPassword {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to Login"})
		}
		ah.Logger.Error(fmt.Sprintf("Error Login User, %v", err))
		return
	}

	c.JSON(http.StatusCreated, tokenRes)
}

func (ah *AuthHandler) HandleLogout(c *gin.Context) {
	username, err := ah.AuthenticatonService.Authentication.ExtractUserFromContext(c)

	if err != nil {
		ah.Logger.Error(fmt.Sprintf("Username does not exist in Context, %v", err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": err})
	}

	err = ah.AuthenticatonService.TokenStore.RevokeAllTokens(username)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Logged out Successfully",
	})
}

func (ah *AuthHandler) HandleRefresh(c *gin.Context) {
	req, err := middleware.ExtractRefreshTokenFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No refresh token found in request"})
		return
	}

	token, err := ah.AuthenticatonService.TokenStore.GetRefreshToken(req.RefreshToken)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token does not exist. Please log in again"})
		return
	}

	if token.Revoked {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token revoked. Please log in again"})
		return
	}

	expired := token.CreatedAt.Add(ah.AuthenticatonService.Authentication.MaxRefresh).After(time.Now())
	if expired {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token expired. Please log in again"})
		return
	}

	// User
	username, err := ah.AuthenticatonService.UserStore.GetUserByToken(token.Token)

	newTokenPair, err := ah.AuthenticatonService.GenerateAccessTokens(username)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"refresh":       "ok",
		"access_token":  newTokenPair.AccessToken,
		"refresh_token": newTokenPair.RefreshToken,
	})
}

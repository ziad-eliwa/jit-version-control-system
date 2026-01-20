package api

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ziad-eliwa/jit-version-control-system/internal/database"
	"github.com/ziad-eliwa/jit-version-control-system/internal/middleware"
)

type UserHandler struct {
	Authentication middleware.AuthenticationMiddleware
	UserStore      database.UserStore
	Logger         *slog.Logger
}

func (uh *UserHandler) HandleGetProfile(c *gin.Context) {
	usernameCtx, err := uh.Authentication.ExtractUserFromContext(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no username found in context"})
		return
	}

	usernameParam := c.Param("username")

	if usernameParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no username found in url"})
		return
	}

	_, err = uh.UserStore.GetUserbyUsername(usernameParam)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	userProfile := &database.UserProfile{}

	if usernameCtx == usernameParam {
		userProfile, err = uh.UserStore.GetUserSelfProfile(usernameParam)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		uh.Logger.Info(fmt.Sprintf("User Profile retrieved for %v by himself", usernameParam))
	} else {
		userProfile, err = uh.UserStore.GetUserProfile(usernameParam)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		uh.Logger.Info(fmt.Sprintf("User Profile retrieved for %v by %v", usernameParam, usernameCtx))
	}

	c.JSON(http.StatusFound, userProfile)
}

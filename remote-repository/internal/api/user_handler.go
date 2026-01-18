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
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "no username found in context"})
	}

	usernameParam := c.Param("username")

	if usernameParam == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no username found in url"})
	}

	_, err = uh.UserStore.GetUserbyUsername(usernameParam)

	if err != nil {
		if err == sql.ErrNoRows {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "user not found"})
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}

	userProfile := &database.UserProfile{}

	if usernameCtx == usernameParam {
		userProfile, err = uh.UserStore.GetUserSelfProfile(usernameParam)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		uh.Logger.Info(fmt.Sprintf("User Profile retrieved for %v by himself", usernameParam))
	} else {
		userProfile, err = uh.UserStore.GetUserProfile(usernameParam)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		uh.Logger.Info(fmt.Sprintf("User Profile retrieved for %v by %v", usernameParam, usernameCtx))
	}

	c.AbortWithStatusJSON(http.StatusFound, userProfile)
}

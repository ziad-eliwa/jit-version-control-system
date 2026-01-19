package api

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ziad-eliwa/jit-version-control-system/internal/database"
	"github.com/ziad-eliwa/jit-version-control-system/internal/middleware"
	"github.com/ziad-eliwa/jit-version-control-system/internal/services"
)

type RepoHandler struct {
	RepoStore   *database.PostgresRepoStore
	Authorizer  *middleware.AuthenticationMiddleware
	Logger      *slog.Logger
	PushService *services.PushService
	PullService *services.PullService
}

func (rh *RepoHandler) HandleGetRepo(c *gin.Context) {
	privacy, ok := c.Get("PRIVACY")

	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "privacy was not found in context"})
	}

	contributor, ok := c.Get("CONTRIBUTOR")

	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "contributor state was not found in context"})
	}

	if privacy == "PRIVATE" && contributor == false {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "you do not have access to this reposoitory"})
	}

	repository := &database.Repository{}

	c.JSON(http.StatusFound, repository)
}

func (rh *RepoHandler) HandleCreateRepo(c *gin.Context) {
	repo := &database.Repository{}

	err := c.BindJSON(&repo)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "repository was not found in request"})
	}

	_, err = rh.RepoStore.CreateRepo(repo)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}

	c.JSON(http.StatusCreated, gin.H{"message": "repository created succecssfully"})
}

func (rh *RepoHandler) HandleAddRemoteRepo(c *gin.Context) {
	repoOwner := c.GetString("REPOOWNER")
	repoName := c.GetString("REPONAME")

	contributor, ok := c.Get("CONTRIBUTOR")
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "contributor state was not found in context"})
	}

	if contributor == true {
		secret, err := rh.RepoStore.GetRepoSecret(repoOwner, repoName)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		c.JSON(http.StatusAccepted, gin.H{"secret": secret})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "you can not add this repository as remote"})
	}
}

type GrantOrRevokeRequest struct {
	TargetUsername string `json:"target"`
}

func (rh *RepoHandler) HandleGrantAccessOnRepo(c *gin.Context) {
	repoOwner := c.GetString("REPOOWNER")
	repoName := c.GetString("REPONAME")

	var req GrantOrRevokeRequest

	err := c.BindJSON(&req)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no target was specified in request"})
	}

	err = rh.RepoStore.GrantAccessOnRepo(repoOwner, repoName, req.TargetUsername)

	if err != nil {
		if err == sql.ErrNoRows {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "this user already does not have access already"})
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "success granting access"})
}

func (rh *RepoHandler) HandleRevokeAccessOnRepo(c *gin.Context) {
	repoOwner := c.GetString("REPOOWNER")
	repoName := c.GetString("REPONAME")

	var req GrantOrRevokeRequest

	err := c.BindJSON(&req)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no target was specified in request"})
	}

	err = rh.RepoStore.RevokeAccessOnRepo(repoOwner, repoName, req.TargetUsername)

	if err != nil {
		if err == sql.ErrNoRows {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "this user already does not have access already"})
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "success revoking access"})
}

func (rh *RepoHandler) HandleGetAllRepos(c *gin.Context) {
	currentUser, err := rh.Authorizer.ExtractUserFromContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "please log in"})
	}

	user := c.Param("username")
	if user == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing username in request"})
	}

	repos, err := rh.RepoStore.GetAllReposbyUsername(user, currentUser)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}

	c.JSON(http.StatusFound, repos)
}

func (rh *RepoHandler) HandlePush(c *gin.Context) {

}

func (rh *RepoHandler) HandlePull(c *gin.Context) {

}

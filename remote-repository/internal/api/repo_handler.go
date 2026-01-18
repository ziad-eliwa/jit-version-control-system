package api

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/ziad-eliwa/jit-version-control-system/internal/database"
	"github.com/ziad-eliwa/jit-version-control-system/internal/services"
)

type RepoHandler struct {
	RepoStore   *database.PostgresRepoStore
	Logger      *slog.Logger
	PushService *services.PushService
	PullService *services.PullService
}

func (rh *RepoHandler) HandleCreateRepo(c *gin.Context) {

}

func (rh *RepoHandler) HandleAddRemoteRepo(c *gin.Context) {

}

func (rh *RepoHandler) HandleGrantAccessOnRepo(c *gin.Context) {

}

func (rh *RepoHandler) HandleRevokeAccessOnRepo(c *gin.Context) {

}

func (rh *RepoHandler) HandleGetAllRepos(c *gin.Context) {

}

func (rh *RepoHandler) HandleGetRepo(c *gin.Context) {

}

func (rh *RepoHandler) HandlePush(c *gin.Context) {

}

func (rh *RepoHandler) HandlePull(c *gin.Context) {

}

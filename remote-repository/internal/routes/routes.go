package route

import (
	"github.com/gin-gonic/gin"
	"github.com/ziad-eliwa/jit-version-control-system/internal/app"
)

func SetupRoutes(app *app.Application) *gin.Engine {
	r := gin.Default()
	gin.SetMode(gin.DebugMode)
	r.GET("/health", app.CheckHealth)

	r.GET("/", app.Main)

	auth := r.Group("/auth")
	auth.POST("/login", app.AuthHandler.HandleLogin)       // Done
	auth.POST("/register", app.AuthHandler.HandleRegister) // Done

	auth.POST("/refresh", app.AuthHandler.HandleRefresh)                                 // Done
	auth.POST("/logout", app.AuthMiddleware.Autheticate(), app.AuthHandler.HandleLogout) // Done

	user := r.Group("/:username", app.AuthMiddleware.Autheticate())
	user.GET("/", app.UserHandler.HandleGetProfile) // Get Profile

	repo := user.Group("/repo")
	repo.GET("/", app.RepoHandler.HandleGetAllRepos) // Get All user repos
	repo.POST("/", app.RepoHandler.HandleCreateRepo) // Create Repository

	reponame := repo.Group("/:reponame", app.AuthMiddleware.AuthorizePrivacy())
	reponame.GET("/", app.AuthMiddleware.AuthorizeEditAccess(), app.RepoHandler.HandleGetRepo) // Get Repo Details if public

	reponame.GET("/remote", app.AuthMiddleware.AuthorizeEditAccess(), app.RepoHandler.HandleAddRemoteRepo) // Add remote if have access

	reponame.POST("/grant", app.AuthMiddleware.AuthorizeOwnership(), app.RepoHandler.HandleGrantAccessOnRepo)   // Grant Access to a user if you are owner --> Authorization
	reponame.POST("/revoke", app.AuthMiddleware.AuthorizeOwnership(), app.RepoHandler.HandleRevokeAccessOnRepo) // Revoke Access from a user if you are owner --> Authorization

	reponame.POST("/push", app.AuthMiddleware.AuthorizeEditAccess(), app.RepoHandler.HandlePush) // Push if have access
	reponame.GET("/pull", app.AuthMiddleware.AuthorizeEditAccess(), app.RepoHandler.HandlePull)  // Pull if have access

	r.NoRoute(app.NotFound)

	return r
}

/*
1- Handlers and Repostory Store -- Done
2- Testing using Insomnia
3- C++ Add commands
5- Push/Pull Service with AWS S3
6- Add redis for caching and logout 
*/

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

	r.Group("/auth", func(ctx *gin.Context) {
		r.POST("/login", app.AuthHandler.HandleLogin)       // Done
		r.POST("/register", app.AuthHandler.HandleRegister) // Done

		r.GET("/:provider") // Google - GitHub
		r.GET("/:provider/callback")

		r.POST("/refresh", app.AuthHandler.HandleRefresh)                                 // Done
		r.POST("/logout", app.AuthMiddleware.Autheticate(), app.AuthHandler.HandleLogout) // Done
	})

	r.Group("/:username", app.AuthMiddleware.Autheticate(), func(ctx *gin.Context) {
		r.GET("/", app.UserHandler.HandleGetProfile) // Get Profile

		r.Group("/repo", func(ctx *gin.Context) {
			r.GET("/", app.RepoHandler.HandleGetAllRepos) // Get All user repos
			r.POST("/", app.RepoHandler.HandleCreateRepo) // Create Repository
			
			r.Group("/:reponame", app.AuthMiddleware.AuthorizePrivacy(), func(ctx *gin.Context) {
				r.GET("/", app.AuthMiddleware.AuthorizeEditAccess(), app.RepoHandler.HandleGetRepo) // Get Repo Details if public

				r.GET("/remote", app.AuthMiddleware.AuthorizeEditAccess(), app.RepoHandler.HandleAddRemoteRepo)     // Add remote if have access

				r.POST("/grant", app.AuthMiddleware.AuthorizeOwnership(), app.RepoHandler.HandleGrantAccessOnRepo)   // Grant Access to a user if you are owner --> Authorization
				r.POST("/revoke", app.AuthMiddleware.AuthorizeOwnership(), app.RepoHandler.HandleRevokeAccessOnRepo) // Revoke Access from a user if you are owner --> Authorization

				r.POST("/push", app.AuthMiddleware.AuthorizeEditAccess(), app.RepoHandler.HandlePush) // Push if have access
				r.GET("/pull", app.AuthMiddleware.AuthorizeEditAccess(), app.RepoHandler.HandlePull)  // Pull if have access
			})
		})
	})

	r.NoRoute(app.NotFound)

	return r
}

/*
1- Handlers and Repostory Store 
2- Testing using Insomnia
3- C++ Add commands
4- OAuth
5- Push/Pull Service with AWS S3
6- Add redis for caching
*/
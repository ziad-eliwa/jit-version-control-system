package route

import (
	"github.com/gin-gonic/gin"
	"github.com/ziad-eliwa/jit-version-control-system/internal/app"
)

func SetupRoutes(app *app.Application) *gin.Engine {
	r := gin.Default()
	r.GET("/health", app.CheckHealth)

	r.GET("/",app.Main)

	r.Group("/auth", func(ctx *gin.Context) {
		r.POST("/login", app.AuthHandler.HandleLogin) // Done
		r.POST("/register", app.AuthHandler.HandleRegister) // Done

		r.GET("/:provider") // Google - GitHub
		r.GET("/:provider/callback")

		r.POST("/refresh", app.AuthHandler.HandleRefresh) // Done
		r.POST("/logout", app.AuthMiddleware.Autheticate(), app.AuthHandler.HandleLogout) // Done
	})

	// Authenticate checks if a user is logged in not who he is
	// Everyone can access the profiles of any user
	// AuthorizePrivacy Gives a trigger to the handler as if he should display all the repos or only public ones
	// AuthorizeCreate Repo checks if the logged-in user is the one creating for his own
	// AuthorizeOwnership that he is the owner of the repo to revoke or grant access
	// AuthorizeAccess that he has access to the repo to push pull or add remote.

	r.Group("/:username", app.AuthMiddleware.Autheticate(), func(ctx *gin.Context) {
		r.GET("/", app.UserHandler.HandleGetProfile) // Get Profile

		r.Group("/repo", func(ctx *gin.Context) { 
			r.GET("/", app.RepoHandler.HandleGetAllRepos) // Get All user repos
			r.POST("/", app.RepoHandler.HandleCreateRepo)       // Create Repository

			r.Group("/:reponame", app.AuthMiddleware.AuthorizePrivacy(), func(ctx *gin.Context) {
				r.GET("/", app.RepoHandler.HandleGetRepo) // Get Repo Details if public

				r.POST("/remote", app.AuthMiddleware.AuthorizeAccess(), app.RepoHandler.HandleAddRemoteRepo)         // Add remote if have access
				r.POST("/grant", app.AuthMiddleware.AuthorizeOwnership(), app.RepoHandler.HandleGrantAccessOnRepo)   // Grant Access to a user if you are owner --> Authorization
				r.POST("/revoke", app.AuthMiddleware.AuthorizeOwnership(), app.RepoHandler.HandleRevokeAccessOnRepo) // Revoke Access from a user if you are owner --> Authorization

				r.POST("/push", app.AuthMiddleware.AuthorizeAccess(), app.RepoHandler.HandlePush) // Push if have access
				r.GET("/pull", app.AuthMiddleware.AuthorizeAccess(), app.RepoHandler.HandlePull)  // Pull if have access
			})
		})
	})

	r.NoRoute(app.NotFound)

	return r
}

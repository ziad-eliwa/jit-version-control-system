package route

import (
	"github.com/gin-gonic/gin"
	"github.com/ziad-eliwa/jit-version-control-system/internal/app"
)

func SetupRoutes(app *app.Application) *gin.Engine {
	r := gin.Default()
	r.GET("/health", app.CheckHealth)

	r.GET("/")

	r.Group("/auth", func(ctx *gin.Context) {
		r.POST("/login")
		r.POST("/register")

		r.GET("/:provider")
		r.GET("/:provider/callback")

		r.POST("/logout")

		r.POST("/refresh")
		r.GET("/me")
	})

	r.Group("/:username", func(ctx *gin.Context) {
		r.GET("/",)  // Get Profile

		r.Group("/repo", func(ctx *gin.Context) {
			r.GET("/")  // Get All user repos
			r.POST("/") // Create Repository
			r.Group("/:reponame", func(ctx *gin.Context) {
				r.GET("/") // Get Repo Details if public

				r.POST("/remote") // Add remote if have access
				r.POST("/grant")  // Grant Access to a user if you are owner --> Authorization
				r.POST("/revoke") // Revoke Access from a user if you are owner --> Authorization

				r.POST("/push") // Push if have access
				r.GET("/pull")  // Pull if have access
			})
		})
	})

	r.NoRoute(app.NotFound)

	return r
}

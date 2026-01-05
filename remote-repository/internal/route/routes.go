package route

import (
	"github.com/gin-gonic/gin"
	"github.com/ziad-eliwa/jit-version-control-system/internal/app"
)

func SetupRoutes(app *app.Application) *gin.Engine {
	r := gin.Default()
	r.GET("/health", app.CheckHealth)

	r.POST("/register")
	r.POST("/login")

	r.GET("/home")
	
	r.Group("/user",func(ctx *gin.Context) {
		r.GET("/profile")

		r.Group("/repo",func(ctx *gin.Context) {
			r.GET("/") // Get All users repos
			r.POST("/") // Create Repository
			r.Group("/{reponame}", func(ctx *gin.Context) {
				r.GET("/") // Get Repo Details if public
				
				r.GET("/remote") // Add remote if have access				
				r.POST("/grant") // Grant Access to a user if you are owner --> Authorization

				r.POST("/push") // Push if have access
				r.GET("/pull") // Pull if have access
			})			
		})
	})
	r.NoRoute(app.NotFound)
	return r
}
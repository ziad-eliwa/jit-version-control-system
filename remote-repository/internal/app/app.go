package app

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ziad-eliwa/jit-version-control-system/internal/api"
	"github.com/ziad-eliwa/jit-version-control-system/internal/database"
	"github.com/ziad-eliwa/jit-version-control-system/internal/middleware"
	"github.com/ziad-eliwa/jit-version-control-system/migrations"
)

type Application struct {
	Logger *slog.Logger
	DB     *sql.DB
	AuthHandler *api.AuthHandler
	UserHandler *api.UserHandler
	RepoHandler *api.RepoHandler
}

func NewApplication(logger *slog.Logger) (*Application, error) {
	logCtx := context.WithValue(context.Background(), "logger", logger)

	pgDB, err := database.OpenPostgresDB(logCtx)
	if err != nil {
		return nil, err
	}
	err = database.MigrateFS(pgDB, migrations.FS, ".")
	if err != nil {
		panic(err)
	}
	// Stores

	// Services
	
	// Middleware
	authMiddleware := middleware.InitAuthMiddleware()	
	TokenGenerator := middleware.InitTokenGenerator(authMiddleware)
	// Handlers
	authHandler := &api.AuthHandler{
		Logger: logger,
		TokenGenerator: &TokenGenerator,
	}
	return &Application{
		Logger: logger,
		DB:     pgDB,
		AuthHandler: authHandler,
	}, nil
}

func (app *Application) CheckHealth(c *gin.Context) {
	app.Logger.Info("CHECK HEALTH: Kolo Zay El Fol")
}

func (app *Application) NotFound(c *gin.Context) {
	app.Logger.Error("COMMAND NOT FOUND: Mesh Hna Ya 5oya")
	c.JSON(http.StatusNotFound,"Invalid Ya A5oya")
}


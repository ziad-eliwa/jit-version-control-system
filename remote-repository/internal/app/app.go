package app

import (
	"context"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/ziad-eliwa/jit-version-control-system/internal/api"
	"github.com/ziad-eliwa/jit-version-control-system/internal/database"
	"github.com/ziad-eliwa/jit-version-control-system/internal/middleware"
	"github.com/ziad-eliwa/jit-version-control-system/internal/services"
	"github.com/ziad-eliwa/jit-version-control-system/migrations"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type Application struct {
	Logger *slog.Logger
	DB     *sql.DB

	AuthHandler *api.AuthHandler
	UserHandler *api.UserHandler
	RepoHandler *api.RepoHandler

	AuthMiddleware *middleware.AuthenticationMiddleware
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
	userStore := &database.PostgresUserStore{
		DB:     pgDB,
		Logger: logger,
	}
	tokenStore := &database.PostgresTokenStore{
		DB: pgDB,
	}
	repoStore := &database.PostgresRepoStore{
		DB:     pgDB,
		Logger: logger,
	}
	// Middleware
	authMiddleware := &middleware.AuthenticationMiddleware{
		TokenStore:  tokenStore,
		JWTSecret:   os.Getenv("JWT_SECRET"),
		Timeout:     15 * time.Minute,
		MaxRefresh:  24 * 7 * time.Hour,
		Logger:      logger,
		IdentityKey: "username",
	}
	// Services
	authService := services.NewAuthService(userStore, tokenStore)
	pushService := &services.PushService{}
	pullService := &services.PullService{}
	// Handlers
	authHandler := &api.AuthHandler{
		Logger:               logger,
		AuthenticatonService: *authService,
	}
	repoHandler := &api.RepoHandler{
		Logger:      logger,
		RepoStore:   repoStore,
		PushService: pushService,
		PullService: pullService,
	}

	return &Application{
		Logger:         logger,
		DB:             pgDB,
		AuthHandler:    authHandler,
		RepoHandler:    repoHandler,
		AuthMiddleware: authMiddleware,
	}, nil
}

func (app *Application) CheckHealth(c *gin.Context) {
	app.Logger.Info("CHECK HEALTH: Kolo Zay El Fol")
}

func (app *Application) NotFound(c *gin.Context) {
	app.Logger.Error("COMMAND NOT FOUND: Mesh Hna Ya 5oya")
	c.JSON(http.StatusNotFound, gin.H{"message": "Invalid Ya A5oya"})
}

func (app *Application) Main(c *gin.Context) {
	c.JSON(http.StatusFound, gin.H{"message":"Ahlan fe JitHub"})
}

package app

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/ziad-eliwa/jit-version-control-system/internal/database"
	"github.com/ziad-eliwa/jit-version-control-system/migrations"
)

type Application struct {
	Logger *slog.Logger
	DB     *sql.DB
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
	return &Application{
		Logger: logger,
		DB:     pgDB,
	}, nil
}

func (app *Application) CheckHealth(c *gin.Context) {
	app.Logger.Info("CHECK HEALTH: Kolo Zay El Fol")
}

func (app *Application) NotFound(c *gin.Context) {
	app.Logger.Error("PAGE NOT FOUND: Mesh Hna Ya 5oya")
}


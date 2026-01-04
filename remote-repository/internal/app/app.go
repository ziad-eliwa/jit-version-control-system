package app

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/ziad-eliwa/jit-version-control-system/internal/database"
)

type Application struct {
	Logger *slog.Logger
	DB 	*pgx.Conn
}

func NewApplication(logger *slog.Logger) (*Application, error) {
	ctx := context.WithValue(context.Background(),"logger",logger)

	pgDB, err := database.OpenPostgresDB(ctx)

	if err != nil {
		return nil, err	
	}

	return &Application{
		Logger: logger,
		DB: pgDB,
	}, nil
}

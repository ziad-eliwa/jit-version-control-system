package database

import (
	"context"
	"fmt"
	"log/slog"
	"github.com/jackc/pgx/v5"
	"github.com/ziad-eliwa/jit-version-control-system/internal/utils"
)


func OpenPostgresDB(ctx context.Context) (*pgx.Conn, error){
	logger := ctx.Value("logger").(*slog.Logger)
	conn, err := pgx.Connect(ctx,utils.GetConnectionString())

	if err != nil {
		return nil, fmt.Errorf("Database not connected due to an error: %v", err)
	}

	if err = conn.Ping(ctx); err != nil {
		return nil, fmt.Errorf("Database not connected due to an error: %v", err)
	}

	logger.Info("Connected to Database Successfully.")
	return conn, nil
}

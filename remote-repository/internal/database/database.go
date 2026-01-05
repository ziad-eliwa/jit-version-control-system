package database

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"log/slog"

	_"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/ziad-eliwa/jit-version-control-system/internal/utils"
)


func OpenPostgresDB(ctx context.Context) (*sql.DB, error){
	logger := ctx.Value("logger").(*slog.Logger)
	conn, err := sql.Open("pgx",utils.GetConnectionString())

	if err != nil {
		return nil, fmt.Errorf("Database not connected due to an error: %v", err)
	}

	if err = conn.Ping(); err != nil {
		return nil, fmt.Errorf("Database not connected due to an error: %v", err)
	}

	logger.Info("Connected to Database Successfully.")
	return conn, nil
}

func MigrateFS(db *sql.DB, migrationsFS fs.FS, dir string) error {
	goose.SetBaseFS(migrationsFS)
	defer func() {
		goose.SetBaseFS(nil)
	}()
	return Migrate(db, dir)
}

func Migrate(db *sql.DB, dir string) error {
	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	err = goose.Up(db, dir)
	if err != nil {
		return fmt.Errorf("goose up: %w", err)
	}
	return nil
}

package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	migratemysql "github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	"blog-server/migrations"
)

func Open(ctx context.Context, dsn string) (*sql.DB, error) {
	cfg, err := parseDSN(dsn)
	if err != nil {
		return nil, err
	}

	pool, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("open mysql: %w", err)
	}
	pool.SetMaxOpenConns(20)
	pool.SetMaxIdleConns(10)
	pool.SetConnMaxLifetime(30 * time.Minute)

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := pool.PingContext(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping mysql: %w", err)
	}
	return pool, nil
}

// Migrate uses its own connection because migration files need multiStatements,
// which should stay disabled on the application pool.
func Migrate(dsn string) error {
	cfg, err := parseDSN(dsn)
	if err != nil {
		return err
	}
	cfg.MultiStatements = true

	conn, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return fmt.Errorf("open mysql for migration: %w", err)
	}
	defer conn.Close()

	driver, err := migratemysql.WithInstance(conn, &migratemysql.Config{})
	if err != nil {
		return fmt.Errorf("migration driver: %w", err)
	}
	source, err := iofs.New(migrations.FS, ".")
	if err != nil {
		return fmt.Errorf("migration source: %w", err)
	}
	m, err := migrate.NewWithInstance("iofs", source, "mysql", driver)
	if err != nil {
		return fmt.Errorf("init migrate: %w", err)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("apply migrations: %w", err)
	}
	return nil
}

func parseDSN(dsn string) (*mysql.Config, error) {
	cfg, err := mysql.ParseDSN(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse DATABASE_DSN: %w", err)
	}
	cfg.ParseTime = true
	cfg.Loc = time.UTC
	// Report matched rather than changed rows, so an UPDATE with identical values is not mistaken for a missing row.
	cfg.ClientFoundRows = true
	if cfg.Params == nil {
		cfg.Params = map[string]string{}
	}
	cfg.Params["charset"] = "utf8mb4"
	// DATETIME defaults (CURRENT_TIMESTAMP) follow the session zone; it must match Loc.
	cfg.Params["time_zone"] = "'+00:00'"
	cfg.Collation = "utf8mb4_unicode_ci"
	return cfg, nil
}

// Command import-cloudbase loads a CloudBase console export into the MySQL schema.
//
// Usage:
//
//	go run ./cmd/import-cloudbase -dir ./cloudbase-export [-truncate] [-admin-email you@example.com]
//
// The directory holds one file per collection named after it (articles.json, classes.json,
// allComments.json, ...), each either a JSON array or JSON Lines.
package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strings"

	"blog-server/internal/db"
)

var contentTables = []string{
	"comments", "article_tags", "articles", "categories", "tags",
	"moments", "friend_links", "changelogs", "projects",
}

func main() {
	dir := flag.String("dir", "", "directory containing the CloudBase export files")
	truncate := flag.Bool("truncate", false, "delete existing content before importing")
	adminEmail := flag.String("admin-email", "", "comments with this email are marked as written by the admin")
	dryRun := flag.Bool("dry-run", false, "run the import in a transaction and roll it back")
	flag.Parse()

	if *dir == "" {
		fmt.Fprintln(os.Stderr, "-dir is required")
		flag.Usage()
		os.Exit(2)
	}
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		fmt.Fprintln(os.Stderr, "DATABASE_DSN is required")
		os.Exit(2)
	}

	if err := run(context.Background(), dsn, *dir, *truncate, *dryRun, strings.TrimSpace(*adminEmail)); err != nil {
		fmt.Fprintln(os.Stderr, "import failed:", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, dsn, dir string, truncate, dryRun bool, adminEmail string) error {
	fmt.Println("reading export from", dir)
	exp, err := loadExport(dir)
	if err != nil {
		return err
	}

	if err := db.Migrate(dsn); err != nil {
		return err
	}
	pool, err := db.Open(ctx, dsn)
	if err != nil {
		return err
	}
	defer pool.Close()

	if truncate && !dryRun {
		if err := truncateContent(ctx, pool); err != nil {
			return err
		}
	} else if err := ensureEmpty(ctx, pool); err != nil {
		return err
	}

	tx, err := pool.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	imp := &importer{tx: tx, adminEmail: strings.ToLower(adminEmail)}
	if err := imp.run(ctx, exp); err != nil {
		return err
	}
	imp.report()

	if dryRun {
		fmt.Println("dry run: rolled back, nothing was written")
		return nil
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}
	fmt.Println("import committed")
	return nil
}

func ensureEmpty(ctx context.Context, pool *sql.DB) error {
	for _, t := range contentTables {
		var n int
		if err := pool.QueryRowContext(ctx, "SELECT COUNT(*) FROM "+t).Scan(&n); err != nil {
			return err
		}
		if n > 0 {
			return fmt.Errorf("table %s already has %d rows; rerun with -truncate to replace existing content", t, n)
		}
	}
	return nil
}

// TRUNCATE commits implicitly in MySQL, so it runs on a dedicated connection before the import transaction.
func truncateContent(ctx context.Context, pool *sql.DB) error {
	conn, err := pool.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	if _, err := conn.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS = 0"); err != nil {
		return err
	}
	defer conn.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS = 1")

	for _, t := range contentTables {
		if _, err := conn.ExecContext(ctx, "TRUNCATE TABLE "+t); err != nil {
			return fmt.Errorf("truncate %s: %w", t, err)
		}
	}
	fmt.Println("existing content truncated")
	return nil
}

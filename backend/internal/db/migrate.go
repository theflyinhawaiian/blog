package db

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"sort"
	"strings"

	"github.com/jmoiron/sqlx"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func Migrate(database *sqlx.DB) error {
	_, err := database.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			name       VARCHAR(255) NOT NULL PRIMARY KEY,
			applied_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`)
	if err != nil {
		return fmt.Errorf("create schema_migrations: %w", err)
	}

	var applied []string
	if err := database.Select(&applied, `SELECT name FROM schema_migrations ORDER BY name`); err != nil {
		return fmt.Errorf("read schema_migrations: %w", err)
	}
	appliedSet := make(map[string]bool, len(applied))
	for _, name := range applied {
		appliedSet[name] = true
	}

	entries, err := fs.ReadDir(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("read migrations dir: %w", err)
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	for _, entry := range entries {
		name := entry.Name()
		if !strings.HasSuffix(name, ".sql") || appliedSet[name] {
			continue
		}

		sql, err := fs.ReadFile(migrationsFS, "migrations/"+name)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", name, err)
		}

		log.Printf("applying migration: %s", name)
		for _, stmt := range strings.Split(string(sql), ";") {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}
			if _, err := database.Exec(stmt); err != nil {
				return fmt.Errorf("migration %s: %w", name, err)
			}
		}

		if _, err := database.Exec(`INSERT INTO schema_migrations (name) VALUES (?)`, name); err != nil {
			return fmt.Errorf("record migration %s: %w", name, err)
		}
		log.Printf("applied migration: %s", name)
	}

	return nil
}

package migration

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// RunMigrations executes .sql files found in migrationsPath (transactional, sorted).
func RunMigrations(db *sql.DB, migrationsPath string) error {
	pattern := filepath.Join(migrationsPath, "*.sql")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("glob migrations: %w", err)
	}
	if len(files) == 0 {
		// nothing to do
		return nil
	}
	sort.Strings(files)
	for _, f := range files {
		// Read file
		content, err := os.ReadFile(f)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", f, err)
		}

		// Start tx
		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("begin tx for %s: %w", f, err)
		}

		if _, err := tx.Exec(string(content)); err != nil {
			tx.Rollback()
			return fmt.Errorf("exec migration %s: %w", f, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit migration %s: %w", f, err)
		}
		// small delay to avoid hammering DB in CI loops (optional)
		time.Sleep(50 * time.Millisecond)
	}
	return nil
}

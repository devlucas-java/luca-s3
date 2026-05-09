package cassandra

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gocql/gocql"
)

func runMigrations(s *gocql.Session, migrationsDir string) error {

	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".cql") {
			continue
		}

		path := filepath.Join(migrationsDir, entry.Name())
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", entry.Name(), err)
		}

		for _, query := range strings.Split(string(content), ";") {
			query = strings.TrimSpace(query)
			// skip blank lines and pure-comment blocks
			if query == "" || strings.HasPrefix(query, "--") {
				continue
			}
			if err := s.Query(query).Exec(); err != nil {
				return fmt.Errorf("migration %s failed on query [%.80s]: %w", entry.Name(), query, err)
			}
		}

		log.Debugf("Executed migration: %s", entry.Name())
	}

	return nil
}

// RunMigrations is exported for use in tests or manual tooling.
func RunMigrations(migrationsDir string) error {
	if session == nil {
		return fmt.Errorf("session is not initialized")
	}
	return runMigrations(session, migrationsDir)
}

func GetSession() *gocql.Session {
	return session
}

func CloseSession() {
	if session != nil {
		session.Close()
		log.Debug("Cassandra session closed")
	}
}

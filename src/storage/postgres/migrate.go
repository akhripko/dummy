package postgres

import (
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // nolint
	_ "github.com/lib/pq"                                // nolint
)

func (s *Storage) MakeMigration() error {
	driver, err := postgres.WithInstance(s.db, &postgres.Config{
		StatementTimeout: 30 * time.Second,
	})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		s.dbName, driver)
	if err != nil {
		return err
	}
	err = m.Up()
	if err == migrate.ErrNoChange {
		return nil
	}
	return err
}

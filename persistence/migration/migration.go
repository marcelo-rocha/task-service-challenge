package migration

import (
	"embed"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed sql/*.sql
var fs embed.FS

func newInstance(url string) (*migrate.Migrate, error) {
	drv, err := iofs.New(fs, "sql")
	if err != nil {
		return nil, err
	}
	m, err := migrate.NewWithSourceInstance("iofs", drv, url)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func Up(url string) error {
	m, err := newInstance(url)
	if err != nil {
		return err
	}
	return m.Up()
}

func Down(url string) error {
	m, err := newInstance(url)
	if err != nil {
		return err
	}
	return m.Down()
}

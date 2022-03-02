package migration

import (
	"embed"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed sql/*.sql
var fs embed.FS

func Up(url string) error {
	drv, err := iofs.New(fs, "sql")
	if err != nil {
		return err
	}
	m, err := migrate.NewWithSourceInstance("iofs", drv, url)
	if err != nil {
		return err
	}
	err = m.Up()
	return err
}

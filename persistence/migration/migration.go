package migration

import (
	"database/sql"
	"embed"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed sql/*.sql
var fs embed.FS

func newInstance(db *sql.DB) (*migrate.Migrate, error) {
	src, err := iofs.New(fs, "sql")
	if err != nil {
		return nil, err
	}
	driver, _ := mysql.WithInstance(db, &mysql.Config{})

	m, err := migrate.NewWithInstance("embed sql", src, "mysql", driver)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func Up(db *sql.DB) error {
	m, err := newInstance(db)
	if err != nil {
		return err
	}
	return m.Up()
}

func Down(db *sql.DB) error {
	m, err := newInstance(db)
	if err != nil {
		return err
	}
	return m.Down()
}

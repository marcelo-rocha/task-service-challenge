package persistence

import (
	"context"
	"database/sql"
	"time"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	_ "github.com/go-sql-driver/mysql"
)

type Connection struct {
	Driver *sql.DB
	db     *goqu.Database
}

func NewConnection(ctx context.Context, url string) (*Connection, error) {
	drv, err := sql.Open("mysql", url)
	if err != nil {
		return nil, err
	}

	drv.SetConnMaxLifetime(time.Minute * 5)
	drv.SetMaxOpenConns(10)
	drv.SetMaxIdleConns(10)

	db := goqu.New("mysql", drv)
	return &Connection{Driver: drv, db: db}, nil

}

func (c *Connection) Close() {
	if c != nil {
		c.Driver.Close()
	}
}

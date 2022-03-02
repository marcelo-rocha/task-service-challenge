package persistence

import (
	"database/sql"

	"github.com/doug-martin/goqu/v9"
	"github.com/marcelo-rocha/task-service-challenge/domain/entities"
	"go.uber.org/zap"
)

type Users struct {
	conn   *Connection
	logger *zap.Logger
	ds     *goqu.SelectDataset
}

func scanUser(rows *sql.Rows) (entities.User, error) {
	var u entities.User
	err := rows.Scan(&u.Id, &u.Login, &u.Name, &u.Kind, &u.Active)
	if err != nil {
		return entities.User{}, err
	}
	return u, err
}

func NewUsers(conn *Connection, logger *zap.Logger) *Tasks {
	return &Tasks{
		conn:   conn,
		logger: logger,
		ds:     conn.db.From("users"),
	}
}

package persistence

import (
	"context"

	"github.com/marcelo-rocha/task-service-challenge/domain/entities"

	"github.com/doug-martin/goqu/v9"
	"go.uber.org/zap"
)

const DefaultAdminUserId = 1 // This user is inserted on database migration

type Users struct {
	conn   *Connection
	logger *zap.Logger
	ds     *goqu.SelectDataset
}

func NewUsers(conn *Connection, logger *zap.Logger) *Users {
	return &Users{
		conn:   conn,
		logger: logger,
		ds:     conn.db.From("users"),
	}
}

func (u *Users) InsertUser(ctx context.Context, login string, name string,
	kind entities.UserKind, active bool, managerID int64) (int64, error) {
	stmt := u.ds.Insert().Cols("login", "name", "kind", "active", "manager_id").Vals(
		goqu.Vals{login, name, string(kind), active, managerID})
	r, err := stmt.Executor().ExecContext(ctx)
	if err != nil {
		return 0, err
	}
	var newId int64
	if newId, err = r.LastInsertId(); err != nil {
		u.logger.Error("failed to get new user Id", zap.Error(err))
	}
	return newId, nil
}

func (u *Users) Truncate(ctx context.Context) error {
	stmt := u.ds.Truncate()
	_, err := stmt.Executor().ExecContext(ctx)
	return err
}

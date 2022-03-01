package persistence

import (
	"context"
	"database/sql"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/marcelo-rocha/task-service-challenge/domain/entities"
	"go.uber.org/zap"
)

type Tasks struct {
	conn   *Connection
	logger *zap.Logger
	ds     *goqu.SelectDataset
}

func scanTask(rows *sql.Rows) (entities.Task, error) {
	var t entities.Task
	err := rows.Scan(&t.Id, &t.Name, &t.Summary, &t.CreationDate, &t.FinishDate)
	return t, err
}

func NewTasks(conn *Connection, logger *zap.Logger) *Tasks {
	return &Tasks{
		conn:   conn,
		logger: logger,
		ds:     conn.db.From("tasks"),
	}
}

func (t *Tasks) InsertTask(ctx context.Context, name string, summary string, creationDate time.Time) (entities.Task, error) {
	stmt := t.ds.Insert().Cols("name", "summary", "creation_date").Vals(
		goqu.Vals{name},
		goqu.Vals{summary},
		goqu.Vals{creationDate},
	)
	sql, params, _ := stmt.ToSQL()
	r, err := t.conn.Driver.ExecContext(ctx, sql, params...)
	if err != nil {
		return entities.Task{}, err
	}
	var newId int64
	if newId, err = r.LastInsertId(); err != nil {
		t.logger.Error("failed to get new task Id", zap.Error(err))
	}
	return entities.Task{Id: newId}, nil
}

func (t *Tasks) FinalizeTask(ctx context.Context, id int64, finish_date time.Time) error {
	stmt := t.ds.Update().Set(goqu.Record{"finish_date": finish_date}).Where(goqu.Ex{"id": id})
	sql, params, _ := stmt.ToSQL()

	_, err := t.conn.Driver.ExecContext(ctx, sql, params...)
	if err != nil {
		return err
	}

	return nil
}

func (t *Tasks) GetTask(ctx context.Context, id int64) (entities.Task, error) {

	stmt := t.ds.Where(goqu.Ex{"id": id})
	sql, params, _ := stmt.ToSQL()

	rows, err := t.conn.Driver.QueryContext(ctx, sql, params...)
	if err != nil {
		return entities.Task{}, err
	}
	if !rows.Next() {
		return entities.Task{}, entities.ErrUnknownTaskID
	}
	task, err := scanTask(rows)
	return task, err
}

func (t *Tasks) GetTasks(ctx context.Context, lastId int64, limit uint) ([]entities.Task, error) {

	stmt := t.ds.Where(goqu.Ex{"id": goqu.Op{"gt": lastId}}).Order(goqu.C("id").Asc()).Limit(limit)
	sql, params, _ := stmt.ToSQL()

	rows, err := t.conn.Driver.QueryContext(ctx, sql, params...)
	if err != nil {
		return []entities.Task{}, err
	}
	var result []entities.Task
	for rows.Next() {
		task, err := scanTask(rows)
		if err != nil {
			return []entities.Task{}, err
		}
		result = append(result, task)
	}
	return result, nil
}

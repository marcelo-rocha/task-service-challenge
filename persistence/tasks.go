package persistence

import (
	"context"
	"database/sql"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/marcelo-rocha/task-service-challenge/domain"
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
	var nullDate sql.NullTime
	err := rows.Scan(&t.Id, &t.Name, &t.Summary, &t.CreationDate, &nullDate, &t.UserId)
	if err != nil {
		return entities.Task{}, err
	}
	if nullDate.Valid {
		t.FinishDate = nullDate.Time
	}
	return t, err
}

func NewTasks(conn *Connection, logger *zap.Logger) *Tasks {
	return &Tasks{
		conn:   conn,
		logger: logger,
		ds:     conn.db.From("tasks"),
	}
}

func (t *Tasks) InsertTask(ctx context.Context, name string, summary string, creationDate time.Time, userId int64) (int64, error) {
	stmt := t.ds.Insert().Cols("name", "summary", "creation_date", "user_id").Vals(
		goqu.Vals{name, summary, creationDate, userId})
	r, err := stmt.Executor().ExecContext(ctx)
	if err != nil {
		return 0, err
	}
	var newId int64
	if newId, err = r.LastInsertId(); err != nil {
		t.logger.Error("failed to get new task Id", zap.Error(err))
	}
	return newId, nil
}

func (t *Tasks) FinalizeTask(ctx context.Context, id int64, finish_date time.Time) error {
	stmt := t.ds.Update().Set(goqu.Record{"finish_date": finish_date}).Where(
		goqu.Ex{"id": id, "finish_date": nil})
	r, err := stmt.Executor().ExecContext(ctx)
	if err != nil {
		return err
	}
	var count int64
	if count, err = r.RowsAffected(); err != nil {
		return err
	} else if count == 0 {
		_, err := t.GetTask(ctx, id)
		if err != nil {
			return err
		}
		return domain.ErrTaskAlreadyFinalized
	}
	return nil
}

func (t *Tasks) GetTask(ctx context.Context, id int64) (entities.Task, error) {
	stmt := t.ds.Where(goqu.Ex{"id": id})
	rows, err := stmt.Executor().QueryContext(ctx)
	if err != nil {
		return entities.Task{}, err
	}
	if !rows.Next() {
		return entities.Task{}, domain.ErrTaskNotFound
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

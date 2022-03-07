package persistence_test

import (
	"context"
	"testing"
	"time"

	"github.com/marcelo-rocha/task-service-challenge/domain"
	"github.com/marcelo-rocha/task-service-challenge/domain/entities"
	"github.com/marcelo-rocha/task-service-challenge/persistence"
	"github.com/stretchr/testify/require"
)

func TestInsertTask(t *testing.T) {
	repo := persistence.NewTasks(db, logger)

	aTask := entities.Task{
		Name:         "create spreasheet",
		Summary:      "create a spreasheet with project costs",
		CreationDate: time.Date(2022, 1, 2, 12, 0, 0, 0, time.UTC),
	}
	taskId, err := repo.InsertTask(context.Background(), aTask.Name, aTask.Summary, aTask.CreationDate, DemoUserId)
	require.NoError(t, err)

	readTask, err := repo.GetTask(context.Background(), taskId)
	require.NoError(t, err)

	require.Equal(t, aTask.Name, readTask.Name)
	require.Equal(t, aTask.Summary, readTask.Summary)
	require.Equal(t, aTask.CreationDate, readTask.CreationDate)

	tasks, err := repo.GetTasksByUser(context.Background(), 0, 10, DemoUserId)
	require.NoError(t, err)

	require.Len(t, tasks, 1)
}

func TestFinishTask(t *testing.T) {
	repo := persistence.NewTasks(db, logger)

	aTask := entities.Task{
		Name:         "create spreasheet",
		Summary:      "create a spreasheet with project costs",
		CreationDate: time.Date(2022, 1, 2, 12, 0, 0, 0, time.UTC),
	}
	taskId, err := repo.InsertTask(context.Background(), aTask.Name, aTask.Summary, aTask.CreationDate, DemoUserId)
	require.NoError(t, err)

	finishDate := time.Date(2022, 3, 1, 12, 0, 0, 0, time.UTC)
	err = repo.FinalizeTask(context.Background(), taskId, finishDate)
	require.NoError(t, err)

	readTask, err := repo.GetTask(context.Background(), taskId)
	require.NoError(t, err)

	require.Equal(t, aTask.Name, readTask.Name)
	require.Equal(t, aTask.Summary, readTask.Summary)
	require.Equal(t, aTask.CreationDate, readTask.CreationDate)
	require.Equal(t, finishDate, *readTask.FinishDate)

}

func TestFinishTaskTwice(t *testing.T) {
	repo := persistence.NewTasks(db, logger)

	aTask := entities.Task{
		Name:         "create spreasheet",
		Summary:      "create a spreasheet with project costs",
		CreationDate: time.Date(2022, 1, 2, 12, 0, 0, 0, time.UTC),
	}
	taskId, err := repo.InsertTask(context.Background(), aTask.Name, aTask.Summary, aTask.CreationDate, DemoUserId)
	require.NoError(t, err)

	finishDate := time.Date(2022, 3, 1, 12, 0, 0, 0, time.UTC)
	err = repo.FinalizeTask(context.Background(), taskId, finishDate)
	require.NoError(t, err)

	err = repo.FinalizeTask(context.Background(), taskId, finishDate)
	require.ErrorIs(t, err, domain.ErrTaskAlreadyFinalized)

}

package persistence

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestInsertTask(t *testing.T) {
	tasks := NewTasks(db, logger)

	aTask, err := tasks.InsertTask(context.Background(), "create spreasheet", "create a spreasheet with project costs", time.Now())
	require.NoError(t, err)

	readTask, err := tasks.GetTask(context.Background(), aTask.Id)
	require.NoError(t, err)

	require.Equal(t, aTask.Name, readTask.Name)
	require.Equal(t, aTask.Summary, readTask.Summary)
}

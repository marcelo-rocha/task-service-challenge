package task

import (
	"context"

	"github.com/marcelo-rocha/task-service-challenge/domain/entities"
)

type GetTaskUseCase struct {
	persistence GetTaskPersister
}

type GetTaskPersister interface {
	GetTask(ctx context.Context, id int64) (entities.Task, error)
}

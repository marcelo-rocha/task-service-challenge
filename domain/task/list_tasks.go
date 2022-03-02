package task

import (
	"context"

	"github.com/marcelo-rocha/task-service-challenge/domain/entities"
	"github.com/marcelo-rocha/task-service-challenge/domain/user"
)

//go:generate moq -stub -pkg mocks -out mocks/list_tasks_persister.go . ListTasksPersister

type ListTasksUseCase struct {
	Persistence ListTasksPersister
}

type ListTasksPersister interface {
	GetTasks(ctx context.Context, lastId int64, limit uint) ([]entities.Task, error)
}

func (c *ListTasksUseCase) ListTasks(ctx context.Context, lastTaskId int64, limit uint) ([]entities.Task, error) {

	userInfo, err := user.GetAuthenticatedUserInfo(ctx)
	if err != nil {
		return []entities.Task{}, err
	}

	if userInfo.Kind != entities.Manager {
		return []entities.Task{}, ErrNotAllowed
	}

	aTasks, err := c.Persistence.GetTasks(ctx, lastTaskId, limit)
	if err != nil {
		return []entities.Task{}, err
	}

	return aTasks, nil
}

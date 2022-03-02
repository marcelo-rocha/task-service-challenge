package task

import (
	"context"

	"github.com/marcelo-rocha/task-service-challenge/domain/entities"
	"github.com/marcelo-rocha/task-service-challenge/domain/user"
)

//go:generate moq -stub -pkg mocks -out mocks/get_task_persister.go . GetTaskPersister

type GetTaskUseCase struct {
	persistence GetTaskPersister
}

type GetTaskPersister interface {
	GetTask(ctx context.Context, id int64) (entities.Task, error)
}

func (u *GetTaskUseCase) GetTask(ctx context.Context, id int64) (entities.Task, error) {

	userInfo, err := user.GetAuthenticatedUserInfo(ctx)
	if err != nil {
		return entities.Task{}, err
	}

	aTask, err := u.persistence.GetTask(ctx, id)
	if err != nil {
		return entities.Task{}, err
	}
	if userInfo.Kind == entities.Manager || userInfo.Id == aTask.UserId {
		return aTask, nil
	}

	return entities.Task{}, ErrNotAllowed
}

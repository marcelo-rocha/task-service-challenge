package task

import (
	"context"
	"time"

	"github.com/marcelo-rocha/task-service-challenge/domain/entities"
	"github.com/marcelo-rocha/task-service-challenge/domain/user"
)

//go:generate moq -stub -pkg mocks -out mocks/finalize_task_persister.go . FinalizeTaskPersister

type FinalizeTaskUseCase struct {
	Persistence FinalizeTaskPersister
}

type FinalizeTaskPersister interface {
	FinalizeTask(ctx context.Context, id int64, finishDate time.Time) error
	GetTask(ctx context.Context, id int64) (entities.Task, error)
}

func (u *FinalizeTaskUseCase) FinalizeTask(ctx context.Context, id int64) error {

	userInfo, err := user.GetAuthenticatedUserInfo(ctx)
	if err != nil {
		return err
	}

	aTask, err := u.Persistence.GetTask(ctx, id)
	if err != nil {
		return err
	}
	if userInfo.Kind != entities.Technician || userInfo.Id != aTask.UserId {
		return ErrNotAllowed
	}

	return u.Persistence.FinalizeTask(ctx, id, time.Now())
}

package task

import (
	"context"
	"time"

	"github.com/marcelo-rocha/task-service-challenge/domain/user"
)

//go:generate moq -stub -pkg mocks -out mocks/new_task_persister.go . NewTaskPersister

type NewTaskUseCase struct {
	Persistence NewTaskPersister
}

type NewTaskPersister interface {
	InsertTask(ctx context.Context, name string, summary string, creationDate time.Time, userId int64) (int64, error)
}

func (c *NewTaskUseCase) NewTask(ctx context.Context, name string, summary string) (int64, error) {
	userInfo, err := user.GetAuthenticatedUserInfo(ctx)
	if err != nil {
		return 0, err
	}
	return c.Persistence.InsertTask(ctx, name, summary, time.Now(), userInfo.Id)

}

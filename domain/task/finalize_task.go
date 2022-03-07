package task

import (
	"context"
	"fmt"
	"time"

	"github.com/marcelo-rocha/task-service-challenge/domain/entities"
	"github.com/marcelo-rocha/task-service-challenge/domain/user"
)

//go:generate moq -stub -pkg mocks -out mocks/finalize_task_persister.go . FinalizeTaskPersister

type FinalizeTaskUseCase struct {
	Persistence FinalizeTaskPersister
	MQ          Publisher
	MsgSubject  string
}

type FinalizeTaskPersister interface {
	FinalizeTask(ctx context.Context, id int64, finishDate time.Time) error
	GetTask(ctx context.Context, id int64) (entities.Task, error)
}

type Publisher interface {
	Publish(subject string, msg string) error
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
	if userInfo.Id != aTask.UserId {
		return ErrNotAllowed
	}

	now := time.Now()
	err = u.Persistence.FinalizeTask(ctx, id, now)
	if err != nil {
		return err
	}
	finishDate, _ := now.UTC().MarshalText()

	u.MQ.Publish(u.MsgSubject, fmt.Sprintf("The tech %v performed the task '%v' on %v", userInfo.Id, aTask.Name, string(finishDate)))

	return nil
}

// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mocks

import (
	"context"
	"github.com/marcelo-rocha/task-service-challenge/domain/entities"
	"github.com/marcelo-rocha/task-service-challenge/domain/task"
	"sync"
	"time"
)

// Ensure, that FinalizeTaskPersisterMock does implement task.FinalizeTaskPersister.
// If this is not the case, regenerate this file with moq.
var _ task.FinalizeTaskPersister = &FinalizeTaskPersisterMock{}

// FinalizeTaskPersisterMock is a mock implementation of task.FinalizeTaskPersister.
//
// 	func TestSomethingThatUsesFinalizeTaskPersister(t *testing.T) {
//
// 		// make and configure a mocked task.FinalizeTaskPersister
// 		mockedFinalizeTaskPersister := &FinalizeTaskPersisterMock{
// 			FinalizeTaskFunc: func(ctx context.Context, id int64, finishDate time.Time) error {
// 				panic("mock out the FinalizeTask method")
// 			},
// 			GetTaskFunc: func(ctx context.Context, id int64) (entities.Task, error) {
// 				panic("mock out the GetTask method")
// 			},
// 		}
//
// 		// use mockedFinalizeTaskPersister in code that requires task.FinalizeTaskPersister
// 		// and then make assertions.
//
// 	}
type FinalizeTaskPersisterMock struct {
	// FinalizeTaskFunc mocks the FinalizeTask method.
	FinalizeTaskFunc func(ctx context.Context, id int64, finishDate time.Time) error

	// GetTaskFunc mocks the GetTask method.
	GetTaskFunc func(ctx context.Context, id int64) (entities.Task, error)

	// calls tracks calls to the methods.
	calls struct {
		// FinalizeTask holds details about calls to the FinalizeTask method.
		FinalizeTask []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// ID is the id argument value.
			ID int64
			// FinishDate is the finishDate argument value.
			FinishDate time.Time
		}
		// GetTask holds details about calls to the GetTask method.
		GetTask []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// ID is the id argument value.
			ID int64
		}
	}
	lockFinalizeTask sync.RWMutex
	lockGetTask      sync.RWMutex
}

// FinalizeTask calls FinalizeTaskFunc.
func (mock *FinalizeTaskPersisterMock) FinalizeTask(ctx context.Context, id int64, finishDate time.Time) error {
	callInfo := struct {
		Ctx        context.Context
		ID         int64
		FinishDate time.Time
	}{
		Ctx:        ctx,
		ID:         id,
		FinishDate: finishDate,
	}
	mock.lockFinalizeTask.Lock()
	mock.calls.FinalizeTask = append(mock.calls.FinalizeTask, callInfo)
	mock.lockFinalizeTask.Unlock()
	if mock.FinalizeTaskFunc == nil {
		var (
			errOut error
		)
		return errOut
	}
	return mock.FinalizeTaskFunc(ctx, id, finishDate)
}

// FinalizeTaskCalls gets all the calls that were made to FinalizeTask.
// Check the length with:
//     len(mockedFinalizeTaskPersister.FinalizeTaskCalls())
func (mock *FinalizeTaskPersisterMock) FinalizeTaskCalls() []struct {
	Ctx        context.Context
	ID         int64
	FinishDate time.Time
} {
	var calls []struct {
		Ctx        context.Context
		ID         int64
		FinishDate time.Time
	}
	mock.lockFinalizeTask.RLock()
	calls = mock.calls.FinalizeTask
	mock.lockFinalizeTask.RUnlock()
	return calls
}

// GetTask calls GetTaskFunc.
func (mock *FinalizeTaskPersisterMock) GetTask(ctx context.Context, id int64) (entities.Task, error) {
	callInfo := struct {
		Ctx context.Context
		ID  int64
	}{
		Ctx: ctx,
		ID:  id,
	}
	mock.lockGetTask.Lock()
	mock.calls.GetTask = append(mock.calls.GetTask, callInfo)
	mock.lockGetTask.Unlock()
	if mock.GetTaskFunc == nil {
		var (
			taskOut entities.Task
			errOut  error
		)
		return taskOut, errOut
	}
	return mock.GetTaskFunc(ctx, id)
}

// GetTaskCalls gets all the calls that were made to GetTask.
// Check the length with:
//     len(mockedFinalizeTaskPersister.GetTaskCalls())
func (mock *FinalizeTaskPersisterMock) GetTaskCalls() []struct {
	Ctx context.Context
	ID  int64
} {
	var calls []struct {
		Ctx context.Context
		ID  int64
	}
	mock.lockGetTask.RLock()
	calls = mock.calls.GetTask
	mock.lockGetTask.RUnlock()
	return calls
}

package persistence

import (
	"context"
	"testing"

	"github.com/marcelo-rocha/task-service-challenge/domain/entities"
	"github.com/stretchr/testify/require"
)

func TestInsertUser(t *testing.T) {
	repo := NewUsers(db, logger)

	aUser := entities.User{
		Login:     "bill",
		Name:      "William Smith",
		Kind:      entities.Technician,
		Active:    true,
		ManagerId: 1,
	}
	_, err := repo.InsertUser(context.Background(), aUser.Login, aUser.Name, aUser.Kind, aUser.Active, aUser.ManagerId)
	require.NoError(t, err)
}

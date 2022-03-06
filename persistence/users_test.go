package persistence_test

import (
	"context"
	"testing"

	"github.com/marcelo-rocha/task-service-challenge/domain/entities"
	"github.com/marcelo-rocha/task-service-challenge/persistence"
	"github.com/stretchr/testify/require"
)

func TestInsertUser(t *testing.T) {
	repo := persistence.NewUsers(db, logger)

	aUser := entities.User{
		Login:     "bill",
		Name:      "William Smith",
		Kind:      entities.Technician,
		Active:    true,
		ManagerId: persistence.DefaultAdminUserId,
	}
	_, err := repo.InsertUser(context.Background(), aUser.Login, aUser.Name, aUser.Kind, aUser.Active, &aUser.ManagerId)
	require.NoError(t, err)
}

package server

import (
	"context"
	"net/http"
	"testing"

	"go.uber.org/zap"

	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/assert"
)

func TestGetTask(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	srv := New(DefaultCfg, logger)
	srv.Init(context.Background())

	apitest.New().
		Handler(srv.Router).
		Get("/api/tasks").
		Expect(t).
		Status(http.StatusUnauthorized).
		End()

	token, err := getAdminToken()
	assert.Nil(t, err)

	apitest.New().
		Handler(srv.Router).
		Get("/api/tasks").
		Header("Authorization", "Bearer "+token).
		Expect(t).
		Body(`{"tasks":[]}`).
		Status(http.StatusOK).
		End()

}

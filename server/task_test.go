package server

import (
	"context"
	"net/http"
	"testing"

	"go.uber.org/zap"

	"github.com/steinfletcher/apitest"
)

func TestGetTask(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	srv := New(DefaultCfg, logger)
	srv.Init(context.Background())
	apitest.New().
		Handler(srv.Router).
		Get("/api/tasks").
		Expect(t).
		Status(http.StatusOK).
		End()
}

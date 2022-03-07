package server

import (
	"context"
	"net/http"
	"testing"

	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
	"github.com/stretchr/testify/require"
)

func TestGetTask(t *testing.T) {
	apitest.New().
		Handler(srv.Router).
		Get("/api/tasks").
		Expect(t).
		Status(http.StatusUnauthorized).
		End()

	token, err := getAdminToken()
	require.NoError(t, err)

	apitest.New().
		Handler(srv.Router).
		Get("/api/tasks").
		Header("Authorization", "Bearer "+token).
		Expect(t).
		Body(`{"tasks":[]}`).
		Status(http.StatusOK).
		End()

}

func TestPostTask(t *testing.T) {
	ctx := context.Background()
	defer tasks.Truncate(ctx)

	demoToken, err := getDemoToken()
	require.NoError(t, err)
	operatorToken, err := getOperatorToken()
	require.NoError(t, err)
	adminToken, err := getAdminToken()
	require.NoError(t, err)

	apitest.New().
		Handler(srv.Router).
		Post("/api/tasks").
		Header("Authorization", "Bearer "+demoToken).
		JSON(`{"name": "setup", "summary": "create spreasheet"}`).
		Expect(t).
		Status(http.StatusCreated).
		End()

	apitest.New().
		Handler(srv.Router).
		Post("/api/tasks").
		Header("Authorization", "Bearer "+operatorToken).
		JSON(`{"name": "prepare", "summary": "prepare workspace"}`).
		Expect(t).
		Status(http.StatusCreated).
		End()

	apitest.New().
		Handler(srv.Router).
		Get("/api/tasks").
		Header("Authorization", "Bearer "+demoToken).
		Expect(t).
		Status(http.StatusOK).
		Assert(jsonpath.Len("$.tasks", 1)).
		End()

	apitest.New().
		Handler(srv.Router).
		Get("/api/tasks").
		Header("Authorization", "Bearer "+adminToken).
		Expect(t).
		Status(http.StatusOK).
		Assert(jsonpath.Len("$.tasks", 2)).
		End()

}

func TestFinishTask(t *testing.T) {
	ctx := context.Background()
	defer tasks.Truncate(ctx)

	demoToken, err := getDemoToken()
	require.NoError(t, err)

	apitest.New().
		Handler(srv.Router).
		Post("/api/tasks").
		Header("Authorization", "Bearer "+demoToken).
		JSON(`{"name": "setup", "summary": "create spreasheet"}`).
		Expect(t).
		Status(http.StatusCreated).
		End()

	apitest.New().
		Handler(srv.Router).
		Post("/api/tasks/1/finishing").
		Header("Authorization", "Bearer "+demoToken).
		Expect(t).
		Status(http.StatusOK).
		End()

}

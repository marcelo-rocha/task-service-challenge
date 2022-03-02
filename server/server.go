package server

import (
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/marcelo-rocha/task-service-challenge/domain/task"
	"go.uber.org/zap"
)

const DefaultSecretKey = "1863a2dfefc2f276e7e164ed2f2f7e975180f2ad7d22c3349f39ded08c11d7f7"

type ServerCfg struct {
	Addr      string `conf:"default:0.0.0.0:8080"`
	SecretKey string
}

type UserCases struct {
	task.NewTaskUseCase
	task.ListTasksUseCase
	task.FinalizeTaskUseCase
}

func Run(cfg *ServerCfg, logger *zap.Logger, uc *UserCases) {

	key := cfg.SecretKey
	if key == "" {
		key = DefaultSecretKey
	}
	var err error
	AuthencationSecretKey, err = hex.DecodeString(key)
	if err != nil {
		panic("decode secret key failed")
	}

	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)

	tasksHandler := TasksHandler{
		NewTaskUseCase:   uc.NewTaskUseCase,
		ListTasksUseCase: uc.ListTasksUseCase,
		Logger:           logger,
	}

	finishHandler := TaskFinishHandler{
		FinalizeTaskUseCase: uc.FinalizeTaskUseCase,
		Logger:              logger,
	}

	s := r.PathPrefix("/api").Subrouter()
	s.Handle("/tasks", &tasksHandler).Methods(http.MethodPost, http.MethodGet)
	//s.HandleFunc("/tasks/{id}", handleGetTask).Methods(http.MethodGet)
	s.Handle("/tasks/{id}", &finishHandler).Methods(http.MethodPatch)
	s.Use(authenticationMiddleware)

	http.Handle("/", r)
	srv := &http.Server{
		Handler:      r,
		Addr:         cfg.Addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	srv.ListenAndServe()
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
}

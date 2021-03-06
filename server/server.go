package server

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/marcelo-rocha/task-service-challenge/domain/task"
	"github.com/marcelo-rocha/task-service-challenge/gateways/msgpublisher"
	"github.com/marcelo-rocha/task-service-challenge/persistence"

	"go.uber.org/zap"
)

const DefaultSecretKey = "1863a2dfefc2f276e7e164ed2f2f7e975180f2ad7d22c3349f39ded08c11d7f7"
const AdminNotificationSubject = "admin_notify"

type ServerCfg struct {
	Addr      string `conf:"env:SV_ADDR,default:0.0.0.0:8080"`
	SecretKey string `conf:"env:SV_SECRET_KEY"`
	DBUrl     string `conf:"env:SV_DB_URL"`
	NATSUrl   string `conf:"env:SV_NATS_URL,default:localhost:4222"`
}

type UseCases struct {
	task.NewTaskUseCase
	task.ListTasksUseCase
	task.FinalizeTaskUseCase
}

type Server struct {
	UseCases
	cfg          *ServerCfg
	logger       *zap.Logger
	srv          *http.Server
	dbConnection *persistence.Connection
	msgPublisher *msgpublisher.MsgPublisher
	Router       *mux.Router
}

func New(cfg *ServerCfg, logger *zap.Logger) *Server {
	return &Server{
		cfg:    cfg,
		logger: logger,
	}
}

func (s *Server) Init(ctx context.Context) error {
	connection, err := persistence.NewConnection(ctx, s.cfg.DBUrl)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	s.dbConnection = connection

	s.msgPublisher = msgpublisher.New(s.logger)
	err = s.msgPublisher.Connect(s.cfg.NATSUrl)
	if err != nil {
		s.logger.Warn("failed on connect to msg publisher", zap.Error(err))
	}

	tasksRepository := persistence.NewTasks(connection, s.logger)

	s.NewTaskUseCase = task.NewTaskUseCase{Persistence: tasksRepository}
	s.ListTasksUseCase = task.ListTasksUseCase{Persistence: tasksRepository}
	s.FinalizeTaskUseCase = task.FinalizeTaskUseCase{
		Persistence: tasksRepository,
		MQ:          s.msgPublisher,
		MsgSubject:  AdminNotificationSubject,
	}

	key := s.cfg.SecretKey
	if key == "" {
		return errors.New("secret key is empty")
	}
	AuthencationSecretKey, err = hex.DecodeString(key)
	if err != nil {
		return fmt.Errorf("failed to decode authentication secret key: %w", err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/", HomeHandler)

	listTasksHandler := &ListTasksHandler{
		ListTasksUseCase: &s.ListTasksUseCase,
		Logger:           s.logger,
	}

	newTasksHandler := &NewTasksHandler{
		NewTaskUseCase: &s.NewTaskUseCase,
		Logger:         s.logger,
	}

	finishHandler := &TaskFinishHandler{
		FinalizeTaskUseCase: &s.FinalizeTaskUseCase,
		Logger:              s.logger,
	}

	sub := router.PathPrefix("/api").Subrouter()
	sub.Handle("/tasks", newTasksHandler).Methods(http.MethodPost)
	sub.Handle("/tasks", listTasksHandler).Methods(http.MethodGet)
	sub.Handle("/tasks/{id}/finishing", finishHandler).Methods(http.MethodPost)
	sub.Use(authenticationMiddleware)

	s.Router = router
	s.srv = &http.Server{
		Handler:      router,
		Addr:         s.cfg.Addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return nil
}

func (s *Server) Run() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) {
	if s.srv != nil {
		s.srv.Shutdown(ctx)
	}
	if s.dbConnection != nil {
		s.dbConnection.Close()
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from Task service"))
}

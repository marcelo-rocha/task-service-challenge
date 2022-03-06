package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/marcelo-rocha/task-service-challenge/domain"

	"github.com/marcelo-rocha/task-service-challenge/domain/entities"
	"github.com/marcelo-rocha/task-service-challenge/domain/task"

	"go.uber.org/zap"
)

type NewTaskRequest struct {
	Name    string `json:"name"`
	Summary string `json:"summary"`
}

type NewTaskResponse struct {
	TaskId int64 `json:"task_id"`
}

type ListTasksResponse struct {
	Tasks []entities.Task `json:"tasks"`
}

type TasksHandler struct {
	*task.NewTaskUseCase
	*task.ListTasksUseCase
	*zap.Logger
}

const DefaultPageLength = 30

func (h *TasksHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var content NewTaskRequest
		err := json.NewDecoder(r.Body).Decode(&content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		taskId, err := h.NewTaskUseCase.NewTask(r.Context(), content.Name, content.Summary)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusCreated)
		responseBody := NewTaskResponse{TaskId: taskId}
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "")
		if err := encoder.Encode(&responseBody); err != nil {
			h.Logger.Warn("write response failed", zap.Error(err))
		}
	} else if r.Method == http.MethodGet {
		lastId := 0
		if s, found := mux.Vars(r)["last_id"]; found {
			var err error
			if lastId, err = strconv.Atoi(s); err != nil {
				lastId = 0
			}
		}
		limit := DefaultPageLength
		if s, found := mux.Vars(r)["limit"]; found {
			var err error
			if limit, err = strconv.Atoi(s); err != nil {
				lastId = DefaultPageLength
			}
		}
		list, err := h.ListTasksUseCase.ListTasks(r.Context(), int64(lastId), uint(limit))
		if err != nil {
			if err == task.ErrNotAllowed {
				http.Error(w, err.Error(), http.StatusForbidden)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		body := ListTasksResponse{Tasks: list}
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "")
		if err := encoder.Encode(&body); err != nil {
			h.Logger.Warn("write response failed", zap.Error(err))
		}
	} else {
		http.Error(w, "unexpected method", http.StatusMethodNotAllowed)
	}
}

type TaskFinishHandler struct {
	*task.FinalizeTaskUseCase
	*zap.Logger
}

func (h *TaskFinishHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var taskId int64
	param := mux.Vars(r)["id"]
	var err error
	if taskId, err = strconv.ParseInt(param, 10, 64); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if err := h.FinalizeTaskUseCase.FinalizeTask(r.Context(), taskId); err != nil {
		if err == domain.ErrTaskNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else if err == task.ErrNotAllowed {
			http.Error(w, err.Error(), http.StatusForbidden)
		} else if err == domain.ErrTaskAlreadyFinalized {
			http.Error(w, err.Error(), http.StatusConflict)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

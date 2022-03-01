package server

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type ServerCfg struct {
	Addr string `conf:"default:0.0.0.0:8080"`
}

func Run(cfg *ServerCfg) {

	//ctx := context.Background()

	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/tasks", handleTasks).Methods("POST", "GET")
	r.HandleFunc("/tasks/{id}", handleGetTask).Methods("GET")
	r.HandleFunc("/tasks/{id}", handleUpdateTask).Methods("PUT")
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
	//r.Write()
}

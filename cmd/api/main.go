package main

import (
	"JobTracker/internal/handlers"
	"JobTracker/internal/repository"
	"log"
	"net/http"
)

func main() {

	jobStore := repository.NewJobStore()
	h := &handlers.TaskHandler{JobTack: jobStore}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /jobs", h.List)
	mux.HandleFunc("POST /jobs", h.Create)
	mux.HandleFunc("GET /jobs/{id}", h.Get)
	mux.HandleFunc("DELETE /jobs/{id}", h.Delete)
	mux.HandleFunc("PUT /jobs/{id}", h.Update)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

}

package handlers

import (
	"JobTracker/internal/dto"
	"JobTracker/internal/repository"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type TaskHandler struct {
	JobTack *repository.JobTrack
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	log.Printf("chamado Create ")

	var in dto.JobInput
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	jobs := h.JobTack.Add(in.Cargo, in.Empresa, in.Status)
	writeJson(w, http.StatusCreated, jobs)
}

func (h *TaskHandler) List(w http.ResponseWriter, r *http.Request) {
	writeJson(w, http.StatusOK, h.JobTack.ListAll())
}

func (h *TaskHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := parserId(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "Id inválido")
		return
	}
	job, found := h.JobTack.Get(id)
	if !found {
		writeError(w, http.StatusNotFound, "Id not found")
		return
	}
	writeJson(w, http.StatusOK, job)
}

func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := parserId(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "Id invalido")
		return
	}
	h.JobTack.DeleteById(id)
	w.WriteHeader(http.StatusNoContent)

}

func writeJson(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJson(w, status, map[string]string{"error": msg})
}

func parserId(r *http.Request) (int, bool) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		return 0, false
	}
	return id, true

}

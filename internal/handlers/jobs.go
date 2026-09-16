package handlers

import (
	"JobTracker/internal/dto"
	"JobTracker/internal/repository"
	"encoding/json"
	"net/http"
	"strconv"
)

type JobHandler struct {
	JobTack *repository.JobTrack
}

func NewHandler(jobStore *repository.JobTrack) *http.ServeMux {
	h := &JobHandler{JobTack: jobStore}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /jobs", h.List)
	mux.HandleFunc("POST /jobs", h.Create)
	mux.HandleFunc("GET /jobs/{id}", h.Get)
	mux.HandleFunc("DELETE /jobs/{id}", h.Delete)
	mux.HandleFunc("PUT /jobs/{id}", h.Update)

	return mux

}

func (h *JobHandler) Create(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var in dto.JobInput
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := in.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	jobs := h.JobTack.Add(in.Cargo, in.Empresa, in.Status)
	writeJson(w, http.StatusCreated, jobs)
}

func (h *JobHandler) List(w http.ResponseWriter, r *http.Request) {

	writeJson(w, http.StatusOK, h.JobTack.ListAll())
}

func (h *JobHandler) Get(w http.ResponseWriter, r *http.Request) {
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

func (h *JobHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := parserId(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "Id invalido")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var in dto.JobInput
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := in.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	jobs := h.JobTack.Update(id, in.Status)
	writeJson(w, http.StatusOK, jobs)
}

func (h *JobHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

package handlers

import (
	"JobTracker/internal/repository"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateJob(t *testing.T) {
	cases := []struct {
		name string
		body string
		want int
	}{
		{"valido", `{"cargo":"Dev","empresa":"I-Gaming", "status":"candidatura enviada"}`, http.StatusCreated},
		{"sem campo", `{"status":"candidatura enviada"}`, http.StatusBadRequest},
		{"campo desconhecido", `{"cargo":"Dev","empresa":"I-Gaming", "status":1}`, http.StatusBadRequest},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/jobs", strings.NewReader(tc.body))
			rec := httptest.NewRecorder()
			NewHandler(repository.NewJobStore()).ServeHTTP(rec, req)
			if rec.Code != tc.want {
				t.Fatalf("Status: %d, Expected: %d", rec.Code, tc.want)
			}
		})
	}
}

func TestGetJob(t *testing.T) {
	cases := []struct {
		name   string
		target string
		want   int
	}{
		{"Found", "/jobs/1", http.StatusOK},
		{"Not Found", "/jobs/2", http.StatusNotFound},
		{"List", "/jobs", http.StatusOK},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tc.target, nil)
			rec := httptest.NewRecorder()
			NewHandler(repository.NewJobStore()).ServeHTTP(rec, req)
			if rec.Code != tc.want {
				t.Fatalf("Status: %d, Expected: %d", rec.Code, tc.want)
			}
		})
	}
}

package handlers

import (
	"JobTracker/internal/repository"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateTask(t *testing.T) {

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

package handlers

import (
	"JobTracker/internal/repository"
	"context"
	"database/sql"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
)

var testPool *pgxpool.Pool

func TestMain(m *testing.M) {
	var err error
	dsn := os.Getenv("DATABSE_URL")
	if dsn == "" {
		dsn = "postgresql://myuser:mypassword@localhost:5430/mytestdatabase"
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}
	if err := goose.SetDialect("postgres"); err != nil {

		log.Fatal(err)
	}

	if err := goose.Up(db, "../../migrations"); err != nil {
		db.Close()
		log.Fatal(err)
	}
	db.Close()

	testPool, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Error on dataBase: %v", err)
	}
	code := m.Run()
	testPool.Close()
	os.Exit(code)

}

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
			NewHandler(repository.NewRepository(testPool)).ServeHTTP(rec, req)
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
		{"Not Found", "/jobs/100", http.StatusNotFound},
		{"List", "/jobs", http.StatusOK},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tc.target, nil)
			rec := httptest.NewRecorder()
			NewHandler(repository.NewRepository(testPool)).ServeHTTP(rec, req)
			if rec.Code != tc.want {
				t.Fatalf("Status: %d, Expected: %d", rec.Code, tc.want)
			}
		})
	}
}

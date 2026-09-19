package main

import (
	"JobTracker/internal/handlers"
	"JobTracker/internal/middleware"
	"JobTracker/internal/repository"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"go.uber.org/fx"
)

func ConnectDb(lc fx.Lifecycle) *pgxpool.Pool {
	ctx := context.Background()
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgresql://myuser:mypassword@localhost:5432/mydatabase"
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}
	if err := goose.SetDialect("postgres"); err != nil {

		log.Fatal(err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		db.Close()
		log.Fatal(err)
	}
	db.Close()

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatal(err)
	}
	cfg.MaxConns = 18
	cfg.MinConns = 2
	cfg.MaxConnLifetime = time.Hour
	cfg.MaxConnIdleTime = 30 * time.Minute
	cfg.HealthCheckPeriod = time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}

	lc.Append(fx.Hook{OnStop: func(ctx context.Context) error {
		log.Println("Closing conections pgxpool...")
		pool.Close()
		return nil
	}})

	pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database connected ! ")
	return pool
}

func NewHttpServer(lc fx.Lifecycle, mux *http.ServeMux) *http.Server {
	srv := &http.Server{Addr: ":8081", Handler: middleware.Logging(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			listener, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			log.Printf("Server running at port: %s", srv.Addr)
			go srv.Serve(listener)
			return nil

		},
		OnStop: func(ctx context.Context) error {
			log.Printf("Stopping server...")
			return srv.Shutdown(ctx)
		},
	})

	return srv

}

func main() {

	fx.New(
		fx.Provide(
			handlers.NewHandler, NewHttpServer, repository.NewRepository, ConnectDb,
		),
		fx.Invoke(func(*http.Server) {}),
	).Run()

}

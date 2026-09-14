package main

import (
	"JobTracker/internal/handlers"
	"JobTracker/internal/middleware"
	"JobTracker/internal/repository"
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"go.uber.org/fx"
)

func NewHttpServer(lc fx.Lifecycle, mux *http.ServeMux) *http.Server {
	srv := &http.Server{Addr: ":8080", Handler: middleware.Logging(mux),
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
			handlers.NewHandler, NewHttpServer, repository.NewJobStore,
		),
		fx.Invoke(func(*http.Server) {}),
	).Run()

}

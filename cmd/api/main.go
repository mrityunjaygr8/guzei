package main

import (
	"context"
	"github.com/go-chi/httplog/v2"
	"github.com/mrityunjaygr8/guzei/internal/postgres_store"
	"github.com/mrityunjaygr8/guzei/internal/server"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	logger := httplog.NewLogger("guzei", httplog.Options{
		// JSON:             true,
		LogLevel:         slog.LevelDebug,
		Concise:          true,
		RequestHeaders:   true,
		MessageFieldName: "message",
		// TimeFieldFormat: time.RFC850,
		Tags: map[string]string{
			"ye": "wo",
			"wo": "ye",
		},
		QuietDownRoutes: []string{
			"/",
			"/ping",
		},
		QuietDownPeriod: 10 * time.Second,
		//SourceFieldName: "source",
	})
	listenAddr := ":3000"
	logger.Info("Starting server on port: ", listenAddr)

	pgStore, err := postgres_store.NewPostgresStore(os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Error("error creating postgres store: ", err)
		os.Exit(1)
	}

	srv, err := server.NewServer(pgStore, logger)
	if err != nil {
		logger.Error("error creating postgres store: ", err)
		os.Exit(1)

	}

	httpServer := http.Server{
		Addr:    listenAddr,
		Handler: srv,
	}

	idleConnectionsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		if err := httpServer.Shutdown(context.Background()); err != nil {
			logger.Error("HTTP Server Shutdown Error: %v", err)
		}
		close(idleConnectionsClosed)
	}()

	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		logger.Error("HTTP server ListenAndServe Error: %v", err)
	}

	<-idleConnectionsClosed

	logger.Info("Bye bye")

}

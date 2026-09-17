package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"library-app-search/internal/bootstrap"
	"library-app-search/internal/config"
)

func main() {
	cfg := config.Load()

	app, err := bootstrap.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	server := &http.Server{
		Addr:              ":" + app.Port,
		Handler:           app.Router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	serverErr := make(chan error, 1)

	go func() {
		log.Printf("Search API listening on :%s", app.Port)

		if err := server.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	signalCtx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	select {
	case <-signalCtx.Done():
		log.Println("Shutdown signal received")

	case err := <-serverErr:
		log.Fatalf("HTTP server failed: %v", err)
	}

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown failed: %v", err)
	}

	if err := app.RedisClient.Close(); err != nil {
		log.Printf("Redis shutdown failed: %v", err)
	}

	log.Println("Search API stopped")
}

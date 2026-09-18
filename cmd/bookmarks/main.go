package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/GoreeCloud/goreecloud-bookmarks/internal/database/postgres"
	"github.com/GoreeCloud/goreecloud-bookmarks/internal/httpapi"
)

const defaultListenAddress = "127.0.0.1:8080"

func main() {
	logger := log.New(os.Stderr, "goreecloud-bookmarks: ", log.Ldate|log.Ltime|log.LUTC)

	handler := httpapi.NewHandler()
	var database *postgres.Database

	if databaseURL := strings.TrimSpace(os.Getenv("GOREECLOUD_BOOKMARKS_DATABASE_URL")); databaseURL != "" {
		var err error
		database, err = postgres.Open(context.Background(), databaseURL)
		if err != nil {
			// Database URLs may contain credentials. Do not emit parser/config
			// details to logs when initialization fails.
			logger.Fatal("database configuration is invalid or could not initialize")
		}
		defer database.Close()

		startupCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		if err := database.CheckStartupCompatibility(startupCtx); err != nil {
			cancel()
			logger.Fatal("database schema is incompatible with this Bookmarks revision")
		}
		cancel()

		handler = httpapi.NewHandlerWithReadiness(httpapi.ReadinessFunc(func(ctx context.Context) httpapi.ReadinessResult {
			result := database.CheckReadiness(ctx)
			return httpapi.ReadinessResult{
				Ready: result.Ready,
				Checks: map[string]string{
					"bookmarks-data": result.State,
				},
			}
		}))
	}

	server := &http.Server{
		Addr:              listenAddress(),
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	shutdownContext, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		<-shutdownContext.Done()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			logger.Printf("graceful shutdown failed: %v", err)
		}
	}()

	logger.Printf("starting GoreeCloud Bookmarks on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Fatalf("server failed: %v", err)
	}
}

func listenAddress() string {
	if value := strings.TrimSpace(os.Getenv("GOREECLOUD_BOOKMARKS_LISTEN_ADDR")); value != "" {
		return value
	}
	return defaultListenAddress
}

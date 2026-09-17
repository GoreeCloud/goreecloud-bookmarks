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

	"github.com/GoreeCloud/goreecloud-bookmarks/internal/httpapi"
)

const defaultListenAddress = "127.0.0.1:8080"

func main() {
	logger := log.New(os.Stderr, "goreecloud-bookmarks: ", log.Ldate|log.Ltime|log.LUTC)

	server := &http.Server{
		Addr:              listenAddress(),
		Handler:           httpapi.NewHandler(),
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

	logger.Printf("starting Experimental service foundation on %s", server.Addr)
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

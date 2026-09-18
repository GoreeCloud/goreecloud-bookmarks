package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/GoreeCloud/goreecloud-bookmarks/internal/database/postgres"
)

func main() {
	logger := log.New(os.Stderr, "goreecloud-bookmarks-migrate: ", log.Ldate|log.Ltime|log.LUTC)
	databaseURL := strings.TrimSpace(os.Getenv("GOREECLOUD_BOOKMARKS_DATABASE_URL"))
	if databaseURL == "" {
		logger.Fatal("GOREECLOUD_BOOKMARKS_DATABASE_URL is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	database, err := postgres.Open(ctx, databaseURL)
	if err != nil {
		logger.Fatal("database configuration is invalid or could not initialize")
	}
	defer database.Close()

	result, err := database.Migrate(ctx)
	if err != nil {
		logger.Fatalf("migration failed: %v", err)
	}

	logger.Printf("migration complete: applied=%d schema_version=%d", result.Applied, result.CurrentVersion)
}

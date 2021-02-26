package main

import (
	"fmt"
	"goschedule/internal/accounts"
	"goschedule/internal/http/rest"
	"goschedule/internal/storage/postgres"
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

const migrationsPath = "file://db/migrations"

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	undo := zap.ReplaceGlobals(logger)
	defer undo()

	if err := godotenv.Load(); err != nil {
		zap.S().Fatal("Error loading .env file")
		panic(err)
	}

	databaseURL := buildDatabaseURL()

	storage, err := postgres.NewStorage(databaseURL)
	if err != nil {
		panic(err)
	}
	defer storage.Close()

	server := rest.NewServer()

	// Set config
	server.Addr = os.Getenv("ADDR")
	server.Port = os.Getenv("PORT")
	server.HashKey = os.Getenv("SESSION_HASH_KEY")
	server.BlockKey = os.Getenv("SESSION_BLOCK_KEY")
	server.GitHubClientID = os.Getenv("GITHUB_CLIENT_ID")
	server.GitHubClientSecret = os.Getenv("GITHUB_CLIENT_SECRET")

	// Configure services
	server.AccountsService = accounts.NewService(storage)

	if err := server.Open(); err != nil {
		zap.S().Fatalf("%+v", err)
	}
}

func buildDatabaseURL() string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_DB"),
	)
}

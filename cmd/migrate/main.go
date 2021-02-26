package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

const migrationsPath = "file://db/migrations"

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	var isUp bool
	flag.BoolVar(&isUp, "up", false, "Should migrate up")
	var isDown bool
	flag.BoolVar(&isDown, "down", false, "Should migrate down")
	flag.Parse()

	m, err := migrate.New(migrationsPath, buildDatabaseURL())
	if err != nil {
		panic(err)
	}

	if isUp {
		migrateUp(m)
	} else if isDown {
		migrateDown(m)
	} else {
		panic(fmt.Errorf("which direction do you want"))
	}
}

func migrateUp(m *migrate.Migrate) {
	if err := m.Up(); err != nil {
		panic(err)
	}
}

func migrateDown(m *migrate.Migrate) {
	if err := m.Down(); err != nil {
		panic(err)
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

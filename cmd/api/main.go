package main

import (
	"os"

	"github.com/Krysik/go-auth/internal/server"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TODO: get port from environment variable
const PORT = "8080"

func main() {
	dbFilePath := os.Getenv("DB_FILE_PATH")

	if dbFilePath == "" {
		dbFilePath = "db.sqlite"
	}

	db, err := gorm.Open(sqlite.Open(dbFilePath), &gorm.Config{})

	if err != nil {
		panic("cannot establish database connection")
	}

	deps := server.AppDeps{
		DB: db,
	}
	server := server.NewServer(&deps)
	server.Logger.Fatal(server.Start(":" + PORT), "failed to start server")
}

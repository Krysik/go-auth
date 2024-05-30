package main

import (
	"log"
	"os"

	"github.com/Krysik/go-auth/internal/server"
	"github.com/joho/godotenv"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TODO: get port from environment variable
const PORT = "8080"

func main() {
	dbFilePath := os.Getenv("DB_FILE_PATH")
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

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
	server.Logger.Fatal(server.Start(":"+PORT), "failed to start server")
}

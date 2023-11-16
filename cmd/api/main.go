package main

import (
	"github.com/Krysik/go-auth/internal/server"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const PORT = "8080"

func main() {
	db, err := gorm.Open(sqlite.Open("db.sqlite"), &gorm.Config{})

	if err != nil {
    panic("failed to connect database")
  }

	deps := server.AppDeps{DB: db}

	server := server.NewServer(&deps)
	server.Logger.Fatal(server.Start(":" + PORT))
}
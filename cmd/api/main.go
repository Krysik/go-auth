package main

import (
	"fmt"
	"log"

	"github.com/Krysik/go-auth/internal/server"
	"github.com/Krysik/go-auth/internal/server/auth"
	"github.com/joho/godotenv"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	loadDotenvErr := godotenv.Load()

	if loadDotenvErr != nil {
		log.Fatal("Error loading .env file")
	}
	env, newEnvErr := server.NewEnv()

	if newEnvErr != nil {
		log.Fatalf("Error while parsing env variables %v", newEnvErr)
	}

	db, err := gorm.Open(sqlite.Open(env.DbFilePath), &gorm.Config{})

	if err != nil {
		panic("cannot establish database connection")
	}

	err = db.AutoMigrate(&auth.Account{}, &auth.RefreshToken{})

	if err != nil {
		panic("failed to migrate database")
	}

	deps := server.AppDeps{
		DB:  db,
		ENV: env,
	}
	server := server.NewServer(&deps)
	server.Logger.Fatal(server.Start(":"+fmt.Sprint(env.Port)), "failed to start server")
}

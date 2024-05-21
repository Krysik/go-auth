package server

import (
	"github.com/Krysik/go-auth/internal/server/auth"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

type AppDeps struct {
	DB *gorm.DB
}

func NewServer(appDeps *AppDeps) *echo.Echo {
	server := echo.New()
	server.Use(middleware.Logger())

	err := appDeps.DB.AutoMigrate(&auth.Account{})

	if err != nil {
		panic("failed to migrate database")
	}

	registerRoutes(server, appDeps)

	return server
}

package server

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type AppDeps struct {
	DB *gorm.DB
}

func NewServer(appDeps *AppDeps) *echo.Echo {
	server := echo.New()

	registerRoutes(server, appDeps.DB)

	return server
}

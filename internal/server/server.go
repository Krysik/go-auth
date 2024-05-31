package server

import (
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type AppDeps struct {
	DB  *gorm.DB
	ENV *ENV
}

var logLevels = map[string]log.Lvl{
	"debug": log.DEBUG,
	"info":  log.INFO,
	"warn":  log.WARN,
	"error": log.ERROR,
}

func NewServer(appDeps *AppDeps) *echo.Echo {
	server := echo.New()
	server.Use(middleware.RequestID())
	server.Use(middleware.Logger())

	server.Use(echoprometheus.NewMiddleware("auth"))
	server.Logger.SetLevel(logLevels[appDeps.ENV.LogLevel])

	server.GET("/metrics", echoprometheus.NewHandler())

	registerRoutes(server, &RouteDeps{
		DB:  appDeps.DB,
		ENV: appDeps.ENV,
	})

	return server
}

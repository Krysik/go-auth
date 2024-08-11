package server

import (
	"strings"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

var logLevels = map[string]log.Lvl{
	"debug": log.DEBUG,
	"info":  log.INFO,
	"warn":  log.WARN,
	"error": log.ERROR,
}

type Server struct {
	DB  *gorm.DB
	ENV *ENV
}

func (s *Server) Initialize() *echo.Echo {
	server := echo.New()
	server.Use(middleware.RequestID())
	server.Use(middleware.Logger())

	server.Use(echoprometheus.NewMiddleware("auth"))
	server.Logger.SetLevel(logLevels[strings.ToLower(s.ENV.LogLevel)])

	server.GET("/metrics", echoprometheus.NewHandler())

	api := &Api{
		DB:     s.DB,
		ENV:    s.ENV,
		Server: server,
	}
	api.RegisterRoutes()

	return server
}

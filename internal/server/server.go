package server

import (
	"net/http"
	"strings"

	"github.com/asaskevich/govalidator"
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

type (
	Server struct {
		DB  *gorm.DB
		ENV *ENV
	}

	StructValidator func(s interface{}) (bool, error)

	CustomValidator struct {
		validator StructValidator
	}
)

func (c *CustomValidator) Validate(i interface{}) error {
	_, err := govalidator.ValidateStruct(i)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, HttpErrorResponse{
			Errors: []HttpError{
				{
					Code:    "BAD_REQUEST",
					Title:   "Validation error",
					Details: err.Error(),
				},
			},
		})
	}

	return err
}

func (s *Server) Initialize() *echo.Echo {
	server := echo.New()
	server.Use(middleware.RequestID())
	server.Use(middleware.Logger())

	server.Use(echoprometheus.NewMiddleware("auth"))
	server.Validator = &CustomValidator{validator: govalidator.ValidateStruct}

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

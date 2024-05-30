package server

import (
	"github.com/Krysik/go-auth/internal/server/auth"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type AppDeps struct {
	DB *gorm.DB
}

func NewServer(appDeps *AppDeps) *echo.Echo {
	server := echo.New()
	server.Use(middleware.RequestID())
	server.Use(middleware.Logger())

	server.Use(echoprometheus.NewMiddleware("auth"))
	env, err := NewEnv()

	if err != nil {
		server.Logger.Fatal("Invalid environment variables ", err.Error())
		panic(err)
	}

	server.Logger.SetLevel(log.INFO)

	server.GET("/metrics", echoprometheus.NewHandler())
	err = appDeps.DB.AutoMigrate(&auth.Account{}, &auth.RefreshToken{})

	if err != nil {
		panic("failed to migrate database")
	}

	registerRoutes(server, &RouteDeps{
		DB:  appDeps.DB,
		ENV: env,
	})

	return server
}

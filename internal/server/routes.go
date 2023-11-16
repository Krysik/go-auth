package server

import (
	"net/http"
	"time"

	"github.com/Krysik/go-auth/internal/server/auth"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type HttpResource struct {
	data interface{}
}

type AccountResource struct {
	Id string
	Type string
	FullName string
	Email string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func registerRoutes(server *echo.Echo, db *gorm.DB) {
	server.GET("/", func (ctx echo.Context) error {
		return ctx.String(http.StatusOK, "Hello world")
	})

	server.POST("/accounts", func (ctx echo.Context) error  {
		acc, err := auth.CreateAccount(
			db,
			auth.NewAccount{FullName: "John Doe", Email: "jdoe@test", Password: "1234"},
		)

		if err != nil {
			return ctx.String(200, "internal server error")
		}
		

		return ctx.JSON(200, &HttpResource{data: acc})
	})
}
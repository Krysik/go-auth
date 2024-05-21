package server

import (
	"net/http"
	"time"

	"github.com/Krysik/go-auth/internal/server/auth"

	"github.com/labstack/echo/v4"
)

type HttpResource struct {
	Data AccountResource
}

type HttpError struct {
	Code string
	Title string
	Details string
}

type HttpErrorResponse struct {
	Errors []HttpError
}

type HttpResourceList struct {
	Data []AccountResource
}

type AccountResource struct {
	Id string
	Type string
	FullName string
	Email string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func registerRoutes(server *echo.Echo, deps *AppDeps) {
	server.GET("/", func (ctx echo.Context) error {
		return ctx.String(http.StatusOK, "Hello world")
	})

	server.POST("/accounts", func (ctx echo.Context) error  {
		/*
		TODO
		get request body
		hash password
		*/

		acc, err := auth.CreateAccount(
			deps.DB,
			auth.NewAccount{FullName: "John Doe", Email: "jdoe@test", Password: "1234"},
		)

		if err != nil {
			ctx.Logger().Error(err.Error(), " failed to create account")

			return ctx.JSON(500, HttpErrorResponse{
				Errors: []HttpError{
					{
						Code: "INTERNAL_SERVER_ERROR",
						Title: "Internal Server Error",
						Details: "Something went wrong",
					},
				},
			})
		}

		return ctx.JSON(200, HttpResource{Data: AccountResource{
			Id: acc.ID,
			Type: "account",
			FullName: acc.FullName,
			Email: acc.Email,
			CreatedAt: acc.CreatedAt,
			UpdatedAt: acc.UpdatedAt,
		}})
	})

	server.GET("/accounts", func(ctx echo.Context) error {
		accounts, err := auth.ListAccounts(deps.DB)

		if err != nil {
			ctx.Logger().Error(err.Error(), "failed to list accounts")
			return ctx.JSON(200, HttpResourceList{
				Data: []AccountResource{},
			})
		}

		var accountResources []AccountResource
		
		for _, account := range accounts {
			accountResources = append(accountResources, AccountResource{
				Id: account.ID,
				Type: "account",
				FullName: account.FullName,
				Email: account.Email,
				CreatedAt: account.CreatedAt,
				UpdatedAt: account.UpdatedAt,
			})
		}

		return ctx.JSON(200, HttpResourceList{
			Data: accountResources,
		})
	})

	server.POST("/sessions", func(ctx echo.Context) error {
		return ctx.String(200, "OK")
	})
}

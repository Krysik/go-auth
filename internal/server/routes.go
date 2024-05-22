package server

import (
	"time"

	"github.com/Krysik/go-auth/internal/server/auth"

	"github.com/labstack/echo/v4"
)

type HttpResource struct {
	Data AccountResource `json:"data"`
}

type HttpError struct {
	Code    string `json:"code"`
	Title   string `json:"title"`
	Details string `json:"details"`
}

type HttpErrorResponse struct {
	Errors []HttpError `json:"errors"`
}

type HttpResourceList struct {
	Data []AccountResource `json:"data"`
}

type AccountResource struct {
	Id        string    `json:"id"`
	Type      string    `json:"type"`
	FullName  string    `json:"fullName"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CreateAccountPayload struct {
	FullName string `json:"fullName"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func registerRoutes(server *echo.Echo, deps *AppDeps) {
	server.POST("/accounts", func(ctx echo.Context) error {
		/*
			TODO: hash password
		*/
		payload := new(CreateAccountPayload)
		if err := ctx.Bind(payload); err != nil {
			return ctx.JSON(400, HttpErrorResponse{
				Errors: []HttpError{
					{
						Code:    "BAD_REQUEST",
						Title:   "Validation error",
						Details: "Invalid body payload",
					},
				},
			})
		}

		acc, err := auth.CreateAccount(
			deps.DB,
			auth.NewAccount{
				FullName: payload.FullName,
				Email:    payload.Email,
				Password: payload.Password,
			},
		)

		if err != nil {
			ctx.Logger().Error(err.Error(), " failed to create account")

			return ctx.JSON(500, HttpErrorResponse{
				Errors: []HttpError{
					{
						Code:    "INTERNAL_SERVER_ERROR",
						Title:   "Internal Server Error",
						Details: "Something went wrong",
					},
				},
			})
		}

		return ctx.JSON(200, HttpResource{Data: AccountResource{
			Id:        acc.ID,
			Type:      "account",
			FullName:  acc.FullName,
			Email:     acc.Email,
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
				Id:        account.ID,
				Type:      "account",
				FullName:  account.FullName,
				Email:     account.Email,
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

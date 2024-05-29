package server

import (
	"github.com/Krysik/go-auth/internal/server/auth"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type createAccountRouteDeps struct {
	DB     *gorm.DB
	Server *echo.Echo
}

type newAccountPayload struct {
	FullName string `json:"fullName"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func registerCreateAccountRoute(deps *createAccountRouteDeps) {
	deps.Server.POST("/accounts", func(ctx echo.Context) error {
		payload := new(newAccountPayload)

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

		return ctx.JSON(201, HttpResource{Data: AccountResource{
			Id:        acc.ID,
			Type:      "account",
			FullName:  acc.FullName,
			Email:     acc.Email,
			CreatedAt: acc.CreatedAt,
			UpdatedAt: acc.UpdatedAt,
		}})
	})
}

package server

import (
	"net/http"

	"github.com/Krysik/go-auth/internal/server/auth"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type CreateAccountRoute struct {
	DB     *gorm.DB
	Server *echo.Echo
}

type newAccountPayload struct {
	FullName string `json:"fullName" valid:"required"`
	Email    string `json:"email" valid:"email"`
	Password string `json:"password" valid:"required"`
}

func (r *CreateAccountRoute) Mount() {
	r.Server.POST("/accounts", func(ctx echo.Context) error {
		payload := new(newAccountPayload)

		if err := ctx.Bind(payload); err != nil {
			return ctx.JSON(http.StatusBadRequest, invalidPayloadResponse)
		}

		if err := ctx.Validate(payload); err != nil {
			return err
		}

		acc, err := auth.CreateAccount(
			r.DB,
			auth.NewAccount{
				FullName: payload.FullName,
				Email:    payload.Email,
				Password: payload.Password,
			},
		)

		if err != nil {
			ctx.Logger().Error(err.Error(), " failed to create account")
			return ctx.JSON(http.StatusInternalServerError, internalServerErrorResponse)
		}

		return ctx.JSON(http.StatusCreated, HttpResource{Data: AccountResource{
			Id:        acc.ID,
			Type:      "account",
			FullName:  acc.FullName,
			Email:     acc.Email,
			CreatedAt: acc.CreatedAt,
			UpdatedAt: acc.UpdatedAt,
		}})
	})
}

package server

import (
	"github.com/Krysik/go-auth/internal/server/auth"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type ListAccountsRoute struct {
	DB     *gorm.DB
	Server *echo.Echo
	ENV    *ENV
}

func (r *ListAccountsRoute) Mount() {
	r.Server.GET("/accounts", func(ctx echo.Context) error {
		accounts, err := auth.ListAccounts(r.DB)

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
	}, func(next echo.HandlerFunc) echo.HandlerFunc {
		return newAuthMiddlewareContext(next, r.ENV.TokenIssuer, r.ENV.JwtSecret)
	})

}

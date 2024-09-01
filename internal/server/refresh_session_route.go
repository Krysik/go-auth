package server

import (
	"net/http"

	"github.com/Krysik/go-auth/internal/server/auth"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type RefreshSessionRoute struct {
	DB     *gorm.DB
	Server *echo.Echo
	ENV    *ENV
}

func (r *RefreshSessionRoute) Mount() {
	r.Server.PATCH("/sessions", func(ctx echo.Context) error {
		refreshTokenCookie, refreshTokenCookieErr := ctx.Cookie("refreshToken")

		if refreshTokenCookieErr != nil {
			return ctx.JSON(http.StatusBadRequest, HttpErrorResponse{
				Errors: []HttpError{
					{
						Code:  "MISSING_REFRESH_TOKEN",
						Title: "Missing refresh token cookie",
					},
				},
			})
		}

		ac := ctx.(*AuthContext)
		accountId := ac.AccountId
		err := r.DB.Transaction(func(tx *gorm.DB) error {
			_, err := auth.GetAccountById(tx, accountId)

			if err != nil {
				return ctx.JSON(http.StatusNotFound, HttpErrorResponse{
					Errors: []HttpError{
						{
							Code:  "ACCOUNT_NOT_FOUND",
							Title: "Account not found",
						},
					},
				})
			}
			refreshToken := refreshTokenCookie.Value
			_, err = auth.GetRefreshToken(tx, refreshToken, accountId)

			if err != nil {
				return ctx.JSON(http.StatusUnauthorized, HttpErrorResponse{
					Errors: []HttpError{
						{
							Code:  "REFRESH_TOKEN_NOT_FOUND",
							Title: "Refresh token not found",
						},
					},
				})
			}

			authTokens, err := auth.GenerateAuthTokens(auth.TokenOpts{
				Issuer:    r.ENV.TokenIssuer,
				JwtSecret: r.ENV.JwtSecret,
				Subject:   accountId,
			})
			if err != nil {
				ctx.Logger().Error(err.Error(), " failed to generate auth tokens")
				return err
			}

			if err = auth.SaveRefreshToken(tx, authTokens.RefreshToken, accountId); err != nil {
				return err
			}

			ctx.SetCookie(&http.Cookie{
				Name:     "accessToken",
				Value:    authTokens.AccessToken,
				HttpOnly: true,
				Secure:   ctx.IsTLS(),
				Path:     "/",
				Expires:  authTokens.AccessTokenTtl,
			})
			ctx.SetCookie(&http.Cookie{
				Name:     "refreshToken",
				Value:    authTokens.RefreshToken,
				HttpOnly: true,
				Secure:   ctx.IsTLS(),
				Path:     "/",
				Expires:  authTokens.RefreshTokenTtl,
			})

			return ctx.NoContent(http.StatusNoContent)
		})

		return err
	}, func(next echo.HandlerFunc) echo.HandlerFunc {
		return newAuthMiddlewareContext(next, r.ENV.TokenIssuer, r.ENV.JwtSecret)
	})
}

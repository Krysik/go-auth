package server

import (
	"net/http"

	"github.com/Krysik/go-auth/internal/server/auth"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type refreshSessionHandlerDeps struct {
	DB     *gorm.DB
	Server *echo.Echo
	ENV    *ENV
}

func registerRefreshSessionRoute(deps *refreshSessionHandlerDeps) {
	deps.Server.PATCH("/sessions", func(ctx echo.Context) error {
		// TODO: run in database transaction
		refreshTokenCookie, refreshTokenCookieErr := ctx.Cookie("refreshToken")

		if refreshTokenCookieErr != nil {
			return ctx.JSON(http.StatusBadRequest, HttpErrorResponse{
				Errors: []HttpError{
					{
						Code:  "MISSING_REFRESH_TOKEN_COOKIE",
						Title: "Missing refresh token cookie",
					},
				},
			})
		}

		ac := ctx.(*AuthContext)
		accountId := ac.AccountId
		_, err := auth.GetAccountById(deps.DB, accountId)

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
		_, err = auth.GetRefreshToken(deps.DB, refreshToken, accountId)

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

		authTokens, err := auth.GenerateAuthTokens("localhost", deps.ENV.JWT_SECRET, accountId)

		if err != nil {
			ctx.Logger().Error(err.Error(), " failed to generate auth tokens")
			return err
		}

		err = auth.SaveRefreshToken(deps.DB, authTokens.RefreshToken, accountId)

		if err != nil {
			return err
		}

		ctx.SetCookie(&http.Cookie{
			Name:     "accessToken",
			Value:    authTokens.AccessToken,
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
			Expires:  authTokens.AccessTokenTtl,
		})
		ctx.SetCookie(&http.Cookie{
			Name:     "refreshToken",
			Value:    authTokens.RefreshToken,
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
			Expires:  authTokens.RefreshTokenTtl,
		})

		return ctx.String(200, "Refreshed")
	}, func(next echo.HandlerFunc) echo.HandlerFunc {
		return newAuthMiddlewareContext(next, deps.ENV.JWT_SECRET, issuer)
	})
}

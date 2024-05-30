package server

import (
	"net/http"

	"github.com/Krysik/go-auth/internal/server/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type AuthContext struct {
	echo.Context
	AccountId string
}

func newAuthMiddlewareContext(next echo.HandlerFunc, issuer, jwtSecret string) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		accessTokenCookie, err := ctx.Request().Cookie("accessToken")

		if err != nil {
			ctx.Logger().Warn("failed to get access token cookie ", err.Error())
			return ctx.JSON(http.StatusUnauthorized, HttpErrorResponse{
				Errors: []HttpError{
					{
						Code:    "UNAUTHORIZED",
						Title:   "Unauthorized",
						Details: "Invalid credentials",
					},
				},
			})
		}

		accessToken := accessTokenCookie.Value

		token, err := jwt.ParseWithClaims(accessToken, &auth.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		}, jwt.WithIssuer(issuer), jwt.WithExpirationRequired())

		if err != nil {
			ctx.Logger().Warn("failed to parse access token ", err)
			return ctx.JSON(http.StatusUnauthorized, HttpErrorResponse{
				Errors: []HttpError{
					{
						Code:    "UNAUTHORIZED",
						Title:   "Unauthorized",
						Details: "Invalid credentials",
					},
				},
			})
		}

		claims := token.Claims.(*auth.TokenClaims)
		accountId := claims.Subject

		ac := &AuthContext{
			AccountId: accountId,
			Context:   ctx,
		}

		return next(ac)
	}
}

package server

import (
	"net/http"

	"github.com/Krysik/go-auth/internal/server/auth"
	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type SignInRoute struct {
	DB     *gorm.DB
	Server *echo.Echo
	ENV    *ENV
}

type NewSessionPayload struct {
	Email    string `json:"email" valid:"email,required"`
	Password string `json:"password" valid:"required"`
}

func (r *SignInRoute) Mount() {
	r.Server.POST("/sessions", func(ctx echo.Context) error {

		payload := new(NewSessionPayload)

		if err := ctx.Bind(payload); err != nil {
			return ctx.JSON(http.StatusBadRequest, HttpErrorResponse{
				Errors: []HttpError{
					{
						Code:    "Bad request",
						Title:   "Validation Error",
						Details: "Invalid payload format",
					},
				},
			})
		}

		_, err := govalidator.ValidateStruct(NewSessionPayload{
			Email:    payload.Email,
			Password: payload.Password,
		})

		ctx.Logger().Infof("payload values email: %s password: %s", payload.Email, payload.Password)

		if err != nil {
			ctx.Logger().Error(err.Error(), "Validation failed")
			return ctx.JSON(http.StatusBadRequest, HttpErrorResponse{
				Errors: []HttpError{
					{
						Code:    "BAD_REQUEST",
						Title:   "Validation error",
						Details: err.Error(),
					},
				},
			})
		}

		account, err := auth.ValidateCredentials(r.DB, payload.Email, payload.Password)

		if err != nil {
			ctx.Logger().Error(err.Error(), " failed to validate credentials")

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
		authTokens, err := auth.GenerateAuthTokens(auth.TokenOpts{
			Issuer:    r.ENV.TokenIssuer,
			JwtSecret: r.ENV.JwtSecret,
			Subject:   account.ID,
		})

		if err != nil {
			ctx.Logger().Error(err.Error(), " failed to generate auth tokens")
			return err
		}

		err = auth.SaveRefreshToken(r.DB, authTokens.RefreshToken, account.ID)

		if err != nil {
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

		return ctx.JSON(200, HttpResource{Data: AccountResource{
			Id:        account.ID,
			Type:      "account",
			FullName:  account.FullName,
			Email:     account.Email,
			CreatedAt: account.CreatedAt,
			UpdatedAt: account.UpdatedAt,
		}})
	})
}

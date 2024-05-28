package server

import (
	"net/http"
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

type NewAccountPayload struct {
	FullName string `json:"fullName"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type NewSessionPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

const issuer = "localhost"

func registerRoutes(server *echo.Echo, deps *AppDeps) {
	server.POST("/accounts", func(ctx echo.Context) error {
		payload := new(NewAccountPayload)

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
	}, func(next echo.HandlerFunc) echo.HandlerFunc {
		return newAuthMiddlewareContext(next, issuer)
	})

	server.POST("/sessions", func(ctx echo.Context) error {
		payload := new(NewSessionPayload)

		if err := ctx.Bind(payload); err != nil {
			ctx.Logger().Error(err.Error(), " failed to bind payload")
			return ctx.JSON(http.StatusBadRequest, HttpErrorResponse{
				Errors: []HttpError{
					{
						Code:    "BAD_REQUEST",
						Title:   "Validation error",
						Details: "Invalid body payload",
					},
				},
			})
		}

		account, err := auth.ValidateCredentials(deps.DB, payload.Email, payload.Password)

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

		authTokens, err := auth.GenerateAuthTokens(issuer, account.ID)

		if err != nil {
			ctx.Logger().Error(err.Error(), " failed to generate auth tokens")
			return err
		}

		err = auth.SaveRefreshToken(deps.DB, authTokens.RefreshToken, account.ID)

		if err != nil {
			return err
		}

		ctx.SetCookie(&http.Cookie{
			Name:     "accessToken",
			Value:    authTokens.AccessToken,
			HttpOnly: true,
			Path:     "/",
			Expires:  authTokens.AccessTokenTtl,
		})
		ctx.SetCookie(&http.Cookie{
			Name:     "refreshToken",
			Value:    authTokens.RefreshToken,
			HttpOnly: true,
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

	server.PATCH("/sessions", func(ctx echo.Context) error {
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

		accountId := ctx.(*AuthContext).AccountId
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

		authTokens, err := auth.GenerateAuthTokens("localhost", accountId)

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
	})
}

package server

import (
	"time"

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

type NewSessionPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

const issuer = "localhost"

func registerRoutes(server *echo.Echo, deps *AppDeps) {
	registerListAccountsRoute(&listAccountsDeps{
		DB:     deps.DB,
		Server: server,
	})

	registerCreateAccountRoute(&createAccountRouteDeps{
		DB:     deps.DB,
		Server: server,
	})

	registerSignInRoute(&signInRouteDeps{
		DB:     deps.DB,
		Server: server,
	})

	registerRefreshSessionRoute(&refreshSessionHandlerDeps{
		DB:     deps.DB,
		Server: server,
	})
}

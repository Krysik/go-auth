package server

import (
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
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

type Api struct {
	Server *echo.Echo
	DB     *gorm.DB
	ENV    *ENV
}

func (api *Api) RegisterRoutes() {
	(&ListAccountsRoute{
		DB:     api.DB,
		Server: api.Server,
		ENV:    api.ENV,
	}).Mount()

	(&CreateAccountRoute{
		DB:     api.DB,
		Server: api.Server,
	}).Mount()

	(&SignInRoute{
		DB:     api.DB,
		Server: api.Server,
		ENV:    api.ENV,
	}).Mount()

	(&RefreshSessionRoute{
		DB:     api.DB,
		Server: api.Server,
		ENV:    api.ENV,
	}).Mount()
}

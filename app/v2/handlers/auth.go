package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/memnix/memnix-rest/app/v2/views/components"
	"github.com/memnix/memnix-rest/domain"
	"github.com/memnix/memnix-rest/services/auth"
)

type AuthController struct {
	useCase auth.IUseCase
}

const (
	// SessionTokenCookieKey is the key for the session token cookie.
	SessionTokenCookieKey = "session_token"
	ExpiresDuration       = 24 * time.Hour
)

func NewAuthController(auth auth.IUseCase) AuthController {
	return AuthController{useCase: auth}
}

func (a *AuthController) PostLogin(c echo.Context) error {
	// Get the username and password from the request
	email := c.FormValue("email")
	password := c.FormValue("password")

	slog.Info("Auth: ", slog.String("email", email), slog.String("password", password))

	// Call the use case to authenticate the user
	jwtToken, err := a.useCase.Login(context.Background(), password, email)
	if err != nil {
		loginError := components.LoginError("Invalid email or password")
		slog.Info("Auth: ", slog.String("error", err.Error()))
		return Render(c, http.StatusForbidden, loginError)
	}

	cookie := &http.Cookie{
		Name:     SessionTokenCookieKey,
		Value:    jwtToken,
		Path:     "/",
		Expires:  time.Now().Add(ExpiresDuration),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(cookie)

	return c.Redirect(http.StatusFound, "/")
}

func (a *AuthController) PostLogout(c echo.Context) error {
	cookie := &http.Cookie{
		Name:     SessionTokenCookieKey,
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(cookie)
	return c.Redirect(http.StatusFound, "/")
}

func (a *AuthController) PostRegister(c echo.Context) error {
	// Get the username and password from the request
	email := c.FormValue("email")
	password := c.FormValue("password")
	username := c.FormValue("username")

	slog.Info("Auth: ", slog.String("email", email), slog.String("username", username))

	registerStruct := domain.Register{
		Email:    email,
		Password: password,
		Username: username,
	}

	// Call the use case to authenticate the user
	_, err := a.useCase.Register(c.Request().Context(), registerStruct)
	if err != nil {
		loginError := components.RegisterError("Invalid email or password")
		slog.Info("Auth: ", slog.String("error", err.Error()))
		return Render(c, http.StatusForbidden, loginError)
	}

	return c.Redirect(http.StatusFound, "/")
}

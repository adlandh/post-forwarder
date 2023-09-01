package driver

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/adlandh/post-forwarder/internal/post-forwarder/config"
	"github.com/adlandh/post-forwarder/internal/post-forwarder/domain"

	"github.com/labstack/echo/v4"
)

type HttpServer struct {
	token string
	app   domain.ApplicationInterface
}

var _ ServerInterface = (*HttpServer)(nil)

func NewHttpServer(cfg *config.Config, app domain.ApplicationInterface) *HttpServer {
	return &HttpServer{
		token: cfg.AuthToken,
		app:   app,
	}
}

func (h HttpServer) HealthCheck(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "Ok")
}

func (h HttpServer) PostWebhook(ctx echo.Context, token string, service string) error {
	return h.webhook(ctx, token, service)
}

func (h HttpServer) GetWebhook(ctx echo.Context, token string, service string) error {
	return h.webhook(ctx, token, service)
}

func (h HttpServer) webhook(ctx echo.Context, token string, service string) error {
	// checking if the token is valid
	if token != h.token {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid auth token")
	}

	var msg string

	// checking query parameters first
	for key, values := range ctx.QueryParams() {
		msg += fmt.Sprintf("%s=%s\n", key, strings.Join(values, ","))
	}

	// checking body parameters
	body, err := io.ReadAll(ctx.Request().Body)

	if err != nil {
		// if parameters were empty, just throw error
		if msg == "" {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		} else { // add error to msg
			msg += fmt.Sprintf("error reading body: %s", err.Error())
		}
	} else { // if no error add body to msg
		msg += string(body)
	}

	err = h.app.ProcessRequest(ctx.Request().Context(), service, msg)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.NoContent(http.StatusOK)
}

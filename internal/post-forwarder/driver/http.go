package driver

import (
	"io"
	"net/http"

	"github.com/adlandh/post-forwarder/internal/post-forwarder/config"
	"github.com/adlandh/post-forwarder/internal/post-forwarder/domain"

	"github.com/labstack/echo/v4"
)

type HttpServer struct {
	token string
	app   domain.ApplicationInterface
}

func (h *HttpServer) Webhook(ctx echo.Context, token string, service string) error {
	if token != h.token {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid auth token")
	}

	body, err := io.ReadAll(ctx.Request().Body)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = h.app.ProcessRequest(ctx.Request().Context(), service, string(body))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.NoContent(http.StatusOK)
}

var _ ServerInterface = (*HttpServer)(nil)

func NewHttpServer(cfg *config.Config, app domain.ApplicationInterface) *HttpServer {
	return &HttpServer{
		token: cfg.AuthToken,
		app:   app,
	}
}

func (h *HttpServer) HealthCheck(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "Ok")
}

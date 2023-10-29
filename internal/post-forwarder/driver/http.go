package driver

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/adlandh/post-forwarder/internal/post-forwarder/config"
	"github.com/adlandh/post-forwarder/internal/post-forwarder/domain"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type HTTPServer struct {
	app   domain.ApplicationInterface
	token string
}

var _ ServerInterface = (*HTTPServer)(nil)

const ErrorSendingResponseMessage = "error sending response: %w"

func NewHTTPServer(cfg *config.Config, app domain.ApplicationInterface) *HTTPServer {
	return &HTTPServer{
		token: cfg.AuthToken,
		app:   app,
	}
}

func (h HTTPServer) HealthCheck(ctx echo.Context) error {
	if err := ctx.String(http.StatusOK, "Ok"); err != nil {
		return fmt.Errorf(ErrorSendingResponseMessage, err)
	}

	return nil
}

func (h HTTPServer) PostWebhook(ctx echo.Context, token string, service string) error {
	return h.webhook(ctx, token, service)
}

func (h HTTPServer) GetWebhook(ctx echo.Context, token string, service string) error {
	return h.webhook(ctx, token, service)
}

func (h HTTPServer) webhook(ctx echo.Context, token string, service string) error {
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
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		} else { // add error to msg
			msg += fmt.Sprintf("error reading body: %s", err.Error())
		}
	} else { // if no error add body to msg
		msg += string(body)
	}

	err = h.app.ProcessRequest(ctx.Request().Context(), getURLFromRequest(ctx.Request()), service, msg)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if err := ctx.NoContent(http.StatusOK); err != nil {
		return fmt.Errorf(ErrorSendingResponseMessage, err)
	}

	return nil
}

func getURLFromRequest(request *http.Request) string {
	scheme := "http"
	if request.TLS != nil {
		scheme += "s"
	}

	return fmt.Sprintf("%s://%s", scheme, request.Host)
}

func (h HTTPServer) ShowMessage(ctx echo.Context, id string) error {
	newUUID, err := uuid.Parse(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	msg, createdAt, err := h.app.GetMessage(ctx.Request().Context(), id)
	if err != nil {
		if errors.Is(domain.ErrorNotFound, err) {
			return echo.NewHTTPError(http.StatusNotFound, domain.ErrorNotFound.Error())
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = ctx.JSON(http.StatusOK, Message{
		CreatedAt: createdAt,
		Id:        newUUID,
		Message:   msg,
	})

	if err != nil {
		return fmt.Errorf(ErrorSendingResponseMessage, err)
	}

	return nil
}

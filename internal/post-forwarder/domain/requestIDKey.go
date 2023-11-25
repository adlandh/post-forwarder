package domain

import (
	"context"

	"github.com/labstack/echo/v4"
)

type requestIDKey string

func (r requestIDKey) String() string {
	return string(r)
}

func (r requestIDKey) Saver(e echo.Context, id string) {
	ctx := context.WithValue(e.Request().Context(), r, id)
	e.SetRequest(e.Request().WithContext(ctx))
}

var RequestID = requestIDKey("request_id")

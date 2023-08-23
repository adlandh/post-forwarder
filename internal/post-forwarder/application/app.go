package application

import (
	"context"

	"github.com/adlandh/post-forwarder/internal/post-forwarder/domain"
)

var _ domain.ApplicationInterface = (*Application)(nil)

type Application struct {
	dst domain.MessageDestination
}

func NewApplication(dst domain.MessageDestination) *Application {
	return &Application{
		dst: dst,
	}
}

func (a Application) ProcessRequest(ctx context.Context, service string, msg string) (err error) {
	return a.dst.Send(ctx, service, msg)
}

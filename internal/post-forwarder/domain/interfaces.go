package domain

import "context"

type ApplicationInterface interface {
	ProcessRequest(ctx context.Context, service string, msg string) (err error)
}

type Notifier interface {
	Send(ctx context.Context, service, msg string) (err error)
}

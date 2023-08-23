package domain

import "context"

type ApplicationInterface interface {
	ProcessRequest(ctx context.Context, service string, msg string) (err error)
}

type MessageDestination interface {
	Send(ctx context.Context, service, msg string) (err error)
}

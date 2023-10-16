// Package domain contains application domain layer
package domain

import "context"

// ApplicationInterface in an interface for application.Application
type ApplicationInterface interface {
	ProcessRequest(ctx context.Context, service string, msg string) (err error)
}

// Notifier in an interface for Notifiers
type Notifier interface {
	Send(ctx context.Context, service, msg string) (err error)
}

// Package domain contains application domain layer
package domain

import (
	"context"
	"fmt"
	"time"
)

// ApplicationInterface in an interface for application.Application
type ApplicationInterface interface {
	ProcessRequest(ctx context.Context, url string, service string, msg string) (err error)
	ShowMessage(ctx context.Context, id string) (msg string, createdAt time.Time, err error)
}

// Notifier in an interface for Notifiers
type Notifier interface {
	Send(ctx context.Context, service, msg string) (err error)
}

var ErrorNotFound = fmt.Errorf("not found")

type MessageStorage interface {
	Store(ctx context.Context, msg string) (id string, err error)
	Read(ctx context.Context, id string) (msg string, createdAt time.Time, err error)
}

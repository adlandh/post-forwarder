// Package domain contains application domain layer
package domain

import (
	"context"
	"fmt"
	"time"
)

type ApplicationInterface interface {
	ProcessRequest(ctx context.Context, url string, service string, msg string) (err error)
	GetMessage(ctx context.Context, id string) (msg string, createdAt time.Time, err error)
}

type Notifier interface {
	Send(ctx context.Context, service, msg string) (err error)
}

var ErrorNotFound = fmt.Errorf("not found")

type MessageStorage interface {
	Store(ctx context.Context, msg string) (id string, err error)
	Read(ctx context.Context, id string) (msg string, createdAt time.Time, err error)
}

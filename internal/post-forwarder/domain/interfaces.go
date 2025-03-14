// Package domain contains application domain layer
//
//go:generate mockery
package domain

import (
	"context"
	"fmt"
	"time"
)

//go:generate gowrap gen -i ApplicationInterface -t https://raw.githubusercontent.com/adlandh/gowrap-templates/main/sentry.gotmpl -o ./wrappers/ApplicationInterfaceWithSentry.go -l "" -g -v InstanceName=application
type ApplicationInterface interface {
	ProcessRequest(ctx context.Context, url string, service string, msg string) (err error)
	GetMessage(ctx context.Context, id string) (msg string, createdAt time.Time, err error)
}

//go:generate gowrap gen -i MessageStorage -t https://raw.githubusercontent.com/adlandh/gowrap-templates/main/sentry.gotmpl -o ./wrappers/MessageStorageWithSentry.go -l "" -g -v InstanceName=notifier
type Notifier interface {
	Send(ctx context.Context, service, msg string) (err error)
}

var ErrorNotFound = fmt.Errorf("not found")

//go:generate gowrap gen -i Notifier -t https://raw.githubusercontent.com/adlandh/gowrap-templates/main/sentry.gotmpl -o ./wrappers/NotifierWithSentry.go -l "" -g -v InstanceName=redis
type MessageStorage interface {
	Store(ctx context.Context, msg string) (id string, err error)
	Read(ctx context.Context, id string) (msg string, createdAt time.Time, err error)
}

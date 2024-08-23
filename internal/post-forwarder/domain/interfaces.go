// Package domain contains application domain layer
package domain

import (
	"context"
	"fmt"
	"time"
)

//go:generate mockery --name ApplicationInterface --with-expecter
//go:generate gowrap gen -p github.com/adlandh/post-forwarder/internal/post-forwarder/domain -i ApplicationInterface -t https://raw.githubusercontent.com/adlandh/gowrap-templates/main/sentry.gotmpl -o ./wrappers/ApplicationInterfaceWithSentry.go -l "" -g
type ApplicationInterface interface {
	ProcessRequest(ctx context.Context, url string, service string, msg string) (err error)
	GetMessage(ctx context.Context, id string) (msg string, createdAt time.Time, err error)
}

//go:generate mockery --name Notifier --with-expecter
//go:generate gowrap gen -p github.com/adlandh/post-forwarder/internal/post-forwarder/domain -i MessageStorage -t https://raw.githubusercontent.com/adlandh/gowrap-templates/main/sentry.gotmpl -o ./wrappers/MessageStorageWithSentry.go -l "" -g
type Notifier interface {
	Send(ctx context.Context, service, msg string) (err error)
}

var ErrorNotFound = fmt.Errorf("not found")

//go:generate mockery --name MessageStorage --with-expecter
//go:generate gowrap gen -p github.com/adlandh/post-forwarder/internal/post-forwarder/domain -i Notifier -t https://raw.githubusercontent.com/adlandh/gowrap-templates/main/sentry.gotmpl -o ./wrappers/NotifierWithSentry.go -l "" -g
type MessageStorage interface {
	Store(ctx context.Context, msg string) (id string, err error)
	Read(ctx context.Context, id string) (msg string, createdAt time.Time, err error)
}

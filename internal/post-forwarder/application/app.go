// Package application contains application layer
package application

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/adlandh/post-forwarder/internal/post-forwarder/domain"
	"go.uber.org/zap"
)

const (
	MaxMessageLength    = 4_096
	ErrorSendingMessage = "error sending message"
)

var r = regexp.MustCompile(`<.*?>`)

var _ domain.ApplicationInterface = (*Application)(nil)

type Application struct {
	notifier domain.Notifier
	logger   *zap.Logger
}

func NewApplication(notifier domain.Notifier, logger *zap.Logger) *Application {
	return &Application{
		notifier: notifier,
		logger:   logger,
	}
}

func (a Application) ProcessRequest(ctx context.Context, service string, msg string) error {
	subject := genSubject(service)

	msg = r.ReplaceAllString(msg, "")

	msg = limitMessageLength(subject, msg)

	err := a.notifier.Send(ctx, subject, msg)
	if err != nil {
		a.logger.Error(ErrorSendingMessage, zap.String("subject", subject), zap.String("msg", msg), zap.Error(err))
		return errors.New(ErrorSendingMessage)
	}

	return nil
}

func genSubject(service string) string {
	return fmt.Sprintf("Message from <b>%s</b>:", service)
}

func limitMessageLength(subject, msg string) string {
	if len([]rune(subject+"\n"+msg)) <= MaxMessageLength {
		return msg
	}

	return msg[:MaxMessageLength-len(subject)-1]
}

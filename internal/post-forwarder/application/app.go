// Package application contains application layer
package application

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	contextlogger "github.com/adlandh/context-logger"
	"github.com/adlandh/post-forwarder/internal/post-forwarder/domain"

	"go.uber.org/zap"
)

const (
	MaxMessageLength    = 100
	ErrorSendingMessage = "error sending message"
)

var r = regexp.MustCompile(`<.*?>`)

var _ domain.ApplicationInterface = (*Application)(nil)

type Application struct {
	notifier domain.Notifier
	logger   *contextlogger.ContextLogger
	storage  domain.MessageStorage
}

func NewApplication(notifier domain.Notifier, logger *zap.Logger, storage domain.MessageStorage) *Application {
	return &Application{
		notifier: notifier,
		logger:   contextlogger.WithContext(logger, contextlogger.WithValueExtractor("request_id")),
		storage:  storage,
	}
}

func (a Application) ProcessRequest(ctx context.Context, url string, service string, msg string) error {
	subject := genSubject(service)

	if isMessageLong(subject, msg) {
		id, err := a.storage.Store(ctx, msg)
		if err != nil {
			a.logger.Ctx(ctx).Error(ErrorSendingMessage, zap.String("subject", subject), zap.String("msg", msg), zap.Error(err))
			return errors.New(ErrorSendingMessage)
		}

		msg = genURL(url, id)
	} else {
		msg = r.ReplaceAllString(msg, "")
	}

	err := a.notifier.Send(ctx, subject, msg)
	if err != nil {
		a.logger.Ctx(ctx).Error(ErrorSendingMessage, zap.String("subject", subject), zap.String("msg", msg), zap.Error(err))
		return errors.New(ErrorSendingMessage)
	}

	return nil
}

func (a Application) GetMessage(ctx context.Context, id string) (msg string, createdAt time.Time, err error) {
	msg, createdAt, err = a.storage.Read(ctx, id)
	if err == nil || errors.Is(err, domain.ErrorNotFound) {
		return
	}

	a.logger.Ctx(ctx).Error("error getting message", zap.String("id", id), zap.Error(err))

	return msg, createdAt, fmt.Errorf("error getting message: %s", id)
}

func genSubject(service string) string {
	return fmt.Sprintf("Message from <b>%s</b>", service)
}

func isMessageLong(subject, msg string) bool {
	return len([]rune(subject+"\n"+msg)) > MaxMessageLength
}

func genURL(url string, id string) string {
	url = url + "/api/message/" + id
	return "Full message is here: " + fmt.Sprintf("<a href=\"%s\">%s</a>", url, url)
}

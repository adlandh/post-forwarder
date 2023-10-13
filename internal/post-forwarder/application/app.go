package application

import (
	"context"
	"fmt"
	"regexp"

	"github.com/adlandh/post-forwarder/internal/post-forwarder/domain"
)

const MaxMessageLength = 4_096

var r = regexp.MustCompile(`<.*?>`)

var _ domain.ApplicationInterface = (*Application)(nil)

type Application struct {
	notifier domain.Notifier
}

func NewApplication(notifier domain.Notifier) *Application {
	return &Application{
		notifier: notifier,
	}
}

func (a Application) ProcessRequest(ctx context.Context, service string, msg string) error {
	subject := genSubject(service)

	msg = r.ReplaceAllString(msg, "")

	msg = limitMessageLength(subject, msg)

	return a.notifier.Send(ctx, subject, msg)
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

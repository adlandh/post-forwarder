package application

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/adlandh/post-forwarder/internal/post-forwarder/domain/mocks"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
)

func TestProcessRequest(t *testing.T) {
	notifiers := new(mocks.Notifier)
	app := NewApplication(notifiers)
	service := gofakeit.Word()
	msg := gofakeit.SentenceSimple()
	ctx := context.Background()

	subject := genSubject(service)

	t.Run("happy case", func(t *testing.T) {
		notifiers.On("Send", ctx, subject, msg).Return(nil).Once()
		err := app.ProcessRequest(ctx, service, msg)
		require.NoError(t, err)
	})

	t.Run("happy case with long string", func(t *testing.T) {
		longMsg := strings.Repeat("A", MaxMessageLength+1)
		require.Greater(t, len(longMsg), MaxMessageLength)
		shortenMsg := limitMessageLength(subject, longMsg)
		require.LessOrEqual(t, len(subject+"\n"+shortenMsg), MaxMessageLength)
		notifiers.On("Send", ctx, subject, shortenMsg).Return(nil).Once()
		err := app.ProcessRequest(ctx, service, longMsg)
		require.NoError(t, err)
	})

	t.Run("error case", func(t *testing.T) {
		fakeErr := errors.New(gofakeit.SentenceSimple())
		notifiers.On("Send", ctx, subject, msg).Return(fakeErr).Once()
		err := app.ProcessRequest(ctx, service, msg)
		require.Error(t, err)
		require.True(t, errors.Is(err, fakeErr))
	})

	notifiers.AssertExpectations(t)

}

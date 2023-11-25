package application

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/adlandh/post-forwarder/internal/post-forwarder/domain"
	"github.com/adlandh/post-forwarder/internal/post-forwarder/domain/mocks"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestProcessRequest(t *testing.T) {
	notifiers := new(mocks.Notifier)
	storage := new(mocks.MessageStorage)
	logger := zaptest.NewLogger(t)
	app := NewApplication(notifiers, logger, storage)
	service := gofakeit.Word()
	msg := gofakeit.SentenceSimple()
	ctx := context.WithValue(context.Background(), domain.RequestID, gofakeit.UUID())
	url := gofakeit.URL()

	subject := genSubject(service)

	t.Run("happy case", func(t *testing.T) {
		if isMessageLong(subject, msg) {
			id := gofakeit.UUID()
			storage.On("Store", ctx, msg).Return(id, nil).Once()
		}

		notifiers.On("Send", ctx, subject, msg).Return(nil).Once()
		err := app.ProcessRequest(ctx, url, service, msg)
		require.NoError(t, err)
	})

	t.Run("happy case with long string", func(t *testing.T) {
		longMsg := strings.Repeat("A", MaxMessageLength+1)
		require.Greater(t, len(longMsg), MaxMessageLength)
		id := gofakeit.UUID()
		storage.On("Store", ctx, longMsg).Return(id, nil).Once()
		notifiers.On("Send", ctx, subject, genURL(url, id)).Return(nil).Once()
		err := app.ProcessRequest(ctx, url, service, longMsg)
		require.NoError(t, err)
	})

	t.Run("error case", func(t *testing.T) {
		fakeErr := errors.New(gofakeit.SentenceSimple())
		notifiers.On("Send", ctx, subject, msg).Return(fakeErr).Once()
		err := app.ProcessRequest(ctx, url, service, msg)
		require.Error(t, err)
		require.Equal(t, err.Error(), ErrorSendingMessage)
	})

	storage.AssertExpectations(t)
	notifiers.AssertExpectations(t)
}

func TestReplaceHtml(t *testing.T) {
	testString := "<b>Test</b>"
	resultString := r.ReplaceAllString(testString, "")
	require.Equal(t, "Test", resultString)

	testString = "<a href=\"https://example.com\">Test</a>"
	resultString = r.ReplaceAllString(testString, "")
	require.Equal(t, "Test", resultString)
}

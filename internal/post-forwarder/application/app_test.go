package application

import (
	"context"
	"strings"
	"testing"

	"github.com/adlandh/post-forwarder/internal/post-forwarder/domain"
	"github.com/adlandh/post-forwarder/internal/post-forwarder/domain/mocks"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestProcessRequest(t *testing.T) {
	notifiers := new(mocks.Notifier)
	storage := new(mocks.MessageStorage)
	logger := zaptest.NewLogger(t)
	app := NewApplication(notifiers, logger, storage)
	service := gofakeit.Word()
	shortMessage := gofakeit.Sentence(3)
	ctx := context.WithValue(context.Background(), domain.RequestID, gofakeit.UUID())
	url := gofakeit.URL()

	subject := genSubject(service)

	t.Run("happy case with short string", func(t *testing.T) {
		notifiers.On("Send", ctx, subject, shortMessage).Return(nil).Once()
		err := app.ProcessRequest(ctx, url, service, shortMessage)
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
		fakeErr := gofakeit.Error()
		notifiers.On("Send", ctx, subject, shortMessage).Return(fakeErr).Once()
		err := app.ProcessRequest(ctx, url, service, shortMessage)
		require.Error(t, err)
		require.ErrorIs(t, err, ErrorSendingMessageError)
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

package application

import (
	"context"
	"strings"
	"testing"

	contextlogger "github.com/adlandh/context-logger"
	"github.com/adlandh/post-forwarder/internal/post-forwarder/domain"
	"github.com/adlandh/post-forwarder/internal/post-forwarder/domain/mocks"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestProcessRequest(t *testing.T) {
	notifiers := mocks.NewNotifier(t)
	storage := mocks.NewMessageStorage(t)
	logger := contextlogger.WithContext(zaptest.NewLogger(t))
	app := NewApplication(notifiers, logger, storage)
	service := gofakeit.Word()
	shortMessage := gofakeit.Sentence()
	ctx := context.WithValue(context.Background(), domain.RequestID, gofakeit.UUID())
	url := gofakeit.URL()

	subject := genSubject(service)

	t.Run("happy case with short string", func(t *testing.T) {
		notifiers.EXPECT().Send(ctx, subject, shortMessage).Return(nil).Once()
		err := app.ProcessRequest(ctx, url, service, shortMessage)
		require.NoError(t, err)
	})

	t.Run("happy case with long string", func(t *testing.T) {
		longMsg := strings.Repeat("A", MaxMessageLength+1)
		require.Greater(t, len(longMsg), MaxMessageLength)
		id := gofakeit.UUID()
		storage.EXPECT().Store(ctx, longMsg).Return(id, nil).Once()
		notifiers.EXPECT().Send(ctx, subject, genURL(url, id)).Return(nil).Once()
		err := app.ProcessRequest(ctx, url, service, longMsg)
		require.NoError(t, err)
	})

	t.Run("error case", func(t *testing.T) {
		fakeErr := gofakeit.Error()
		notifiers.EXPECT().Send(ctx, subject, shortMessage).Return(fakeErr).Once()
		err := app.ProcessRequest(ctx, url, service, shortMessage)
		require.Error(t, err)
		require.ErrorIs(t, err, ErrorSendingMessageError)
	})
}

func TestReplaceHtml(t *testing.T) {
	testString := "<b>Test</b>"
	resultString := r.ReplaceAllString(testString, "")
	require.Equal(t, "Test", resultString)

	testString = "<a href=\"https://example.com\">Test</a>"
	resultString = r.ReplaceAllString(testString, "")
	require.Equal(t, "Test", resultString)
}

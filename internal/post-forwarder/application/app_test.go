package application

import (
	"context"
	"errors"
	"testing"

	"github.com/adlandh/post-forwarder/internal/post-forwarder/domain/mocks"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
)

func TestProcessRequest(t *testing.T) {
	messageDestination := new(mocks.MessageDestination)
	app := NewApplication(messageDestination)
	service := gofakeit.Word()
	msg := gofakeit.SentenceSimple()
	ctx := context.Background()

	t.Run("happy case", func(t *testing.T) {
		messageDestination.On("Send", ctx, service, msg).Return(nil).Once()
		err := app.ProcessRequest(ctx, service, msg)
		require.NoError(t, err)
	})

	t.Run("error case", func(t *testing.T) {
		fakeErr := errors.New(gofakeit.SentenceSimple())
		messageDestination.On("Send", ctx, service, msg).Return(fakeErr).Once()
		err := app.ProcessRequest(ctx, service, msg)
		require.Error(t, err)
		require.True(t, errors.Is(err, fakeErr))
	})

	messageDestination.AssertExpectations(t)

}

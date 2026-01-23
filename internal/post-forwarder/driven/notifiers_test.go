package driven

import (
	"reflect"
	"testing"

	"github.com/adlandh/post-forwarder/internal/post-forwarder/config"
	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/telegram"
	"github.com/stretchr/testify/require"
)

func TestNewNotifiersNilConfig(t *testing.T) {
	_, err := NewNotifiers(nil)
	require.Error(t, err)
}

func TestNewNotifiersEmptyConfig(t *testing.T) {
	cfg := &config.Config{}
	_, err := NewNotifiers(cfg)
	require.Error(t, err)
}

func TestNewNotifiersUnknownService(t *testing.T) {
	cfg := &config.Config{
		Notifiers: []string{"unknown"},
	}

	_, err := NewNotifiers(cfg)
	require.Error(t, err)
}

func TestNewNotifiersSuccess(t *testing.T) {
	origNewTelegramService := newTelegramService
	t.Cleanup(func() {
		newTelegramService = origNewTelegramService
	})
	newTelegramService = func(string) (*telegram.Telegram, error) {
		return &telegram.Telegram{}, nil
	}

	cfg := &config.Config{
		Notifiers: []string{" TELEGRAM ", "sLaCk", "pushover"},
		Telegram: config.TelegramConfig{
			Token:   "telegram-token",
			ChatIDs: []int64{123, 456},
		},
		Slack: config.SlackConfig{
			Token:      "slack-token",
			ChannelIDs: []string{"C1", "C2"},
		},
		Pushover: config.PushoverConfig{
			Token: "pushover-token",
			User:  "pushover-user",
		},
	}

	notifier, err := NewNotifiers(cfg)
	require.NoError(t, err)

	require.Equal(t, 3, notifierServiceCount(t, notifier))
}

func notifierServiceCount(t *testing.T, notifier *notify.Notify) int {
	t.Helper()

	require.NotNil(t, notifier)

	value := reflect.ValueOf(notifier)
	if value.Kind() == reflect.Pointer {
		value = value.Elem()
	}

	field := value.FieldByName("notifiers")
	require.True(t, field.IsValid(), "notifier notifiers field not found")
	require.Equal(t, reflect.Slice, field.Kind(), "notifier notifiers field kind is %s", field.Kind())

	return field.Len()
}

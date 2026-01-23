// Package driven contains driven layer
package driven

import (
	"errors"
	"fmt"
	"strings"

	"github.com/adlandh/post-forwarder/internal/post-forwarder/config"
	"github.com/adlandh/post-forwarder/internal/post-forwarder/domain"
	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/slack"
	"github.com/nikoksr/notify/service/telegram"
)

var _ domain.Notifier = (*notify.Notify)(nil)

var (
	newTelegramService = telegram.New
	newSlackService    = slack.New
	newPushoverService = NewPushover
)

func NewNotifiers(cfg *config.Config) (*notify.Notify, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is nil")
	}

	if len(cfg.Notifiers) == 0 {
		return nil, fmt.Errorf("no notifiers configured")
	}

	notifier := notify.New()

	var (
		used int
		err  error
	)

	for _, rawService := range cfg.Notifiers {
		service := strings.ToLower(strings.TrimSpace(rawService))

		if service == "" {
			continue
		}

		notifyService, serviceErr := buildNotifierService(cfg, service)
		if serviceErr != nil {
			err = errors.Join(err, serviceErr)
			continue
		}

		notifier.UseServices(notifyService)

		used++
	}

	if err != nil {
		return nil, fmt.Errorf("error creating notifiers: %w", err)
	}

	if used == 0 {
		return nil, fmt.Errorf("no notifiers configured")
	}

	return notifier, nil
}

func buildNotifierService(cfg *config.Config, service string) (notify.Notifier, error) {
	switch service {
	case config.TelegramService:
		return buildTelegramService(cfg)
	case config.SlackService:
		return buildSlackService(cfg), nil
	case config.PushoverService:
		return buildPushoverService(cfg), nil
	default:
		return nil, fmt.Errorf("unknown notifier service: %s", service)
	}
}

func buildTelegramService(cfg *config.Config) (*telegram.Telegram, error) {
	telegramService, err := newTelegramService(cfg.Telegram.Token)
	if err != nil {
		return nil, fmt.Errorf("error creating telegram service: %w", err)
	}

	for _, chatID := range cfg.Telegram.ChatIDs {
		telegramService.AddReceivers(chatID)
	}

	return telegramService, nil
}

func buildSlackService(cfg *config.Config) *slack.Slack {
	slackService := newSlackService(cfg.Slack.Token)

	for _, channelID := range cfg.Slack.ChannelIDs {
		slackService.AddReceivers(channelID)
	}

	return slackService
}

func buildPushoverService(cfg *config.Config) *Pushover {
	pushoverService := newPushoverService(cfg.Pushover.Token)
	pushoverService.AddReceivers(cfg.Pushover.User)

	return pushoverService
}

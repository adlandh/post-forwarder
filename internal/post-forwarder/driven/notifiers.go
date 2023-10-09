package driven

import (
	"fmt"

	"github.com/adlandh/post-forwarder/internal/post-forwarder/config"
	"github.com/adlandh/post-forwarder/internal/post-forwarder/domain"
	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/slack"
	"github.com/nikoksr/notify/service/telegram"
)

var _ domain.Notifier = (*notify.Notify)(nil)

func NewNotifiers(cfg *config.Config) (*notify.Notify, error) {
	notifier := notify.Default()
	for _, service := range cfg.Notifiers {
		switch service {
		case config.TelegramService:
			telegramService, err := telegram.New(cfg.Telegram.Token)
			if err != nil {
				return nil, fmt.Errorf("error creating telegram service: %w", err)
			}

			telegramService.SetParseMode(telegram.ModeMarkdown)

			for _, chatID := range cfg.Telegram.ChatIDs {
				telegramService.AddReceivers(chatID)
			}

			notify.UseServices(telegramService)
		case config.SlackService:
			slackService := slack.New(cfg.Slack.Token)

			for _, channelID := range cfg.Slack.ChannelIDs {
				slackService.AddReceivers(channelID)
			}

			notify.UseServices(slackService)
		}
	}
	return notifier, nil
}
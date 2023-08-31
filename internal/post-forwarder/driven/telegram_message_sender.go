package driven

import (
	"context"
	"fmt"

	"github.com/adlandh/post-forwarder/internal/post-forwarder/config"
	"github.com/adlandh/post-forwarder/internal/post-forwarder/domain"
	"github.com/go-telegram/bot"
)

const MaxMessageLength = 4_096

var _ domain.MessageDestination = (*TelegramMessageSender)(nil)

type TelegramMessageSender struct {
	bot    *bot.Bot
	chatId int64
	token  string
}

func NewTelegramMessageSender(cfg *config.Config) (*TelegramMessageSender, error) {
	t := &TelegramMessageSender{
		chatId: cfg.BotChatID,
		token:  cfg.BotToken,
	}

	b, err := bot.New(cfg.BotToken)
	if err != nil {
		return nil, err
	}

	t.bot = b

	return t, nil
}

func (t TelegramMessageSender) Send(ctx context.Context, service, msg string) error {
	fullMessage := fmt.Sprintf("Message from %s: %s", service, msg)
	if len(fullMessage) > MaxMessageLength {
		fullMessage = fullMessage[:MaxMessageLength]
	}

	_, err := t.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: t.chatId,
		Text:   fullMessage,
	})

	return err
}

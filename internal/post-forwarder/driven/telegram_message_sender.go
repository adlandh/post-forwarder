package driven

import (
	"context"
	"fmt"

	"github.com/adlandh/post-forwarder/internal/post-forwarder/config"
	"github.com/adlandh/post-forwarder/internal/post-forwarder/domain"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
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
	fullMessage := fmt.Sprintf("*%s*:\n%s", bot.EscapeMarkdown(service), bot.EscapeMarkdown(msg))
	if len([]rune(fullMessage)) > MaxMessageLength {
		fullMessage = fullMessage[:MaxMessageLength]
	}

	_, err := t.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    t.chatId,
		Text:      fullMessage,
		ParseMode: models.ParseModeMarkdown,
	})

	return err
}

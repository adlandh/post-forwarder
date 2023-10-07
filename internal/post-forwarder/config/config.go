package config

import (
	"fmt"
	"strings"

	"github.com/caarlos0/env/v9"
)

const TelegramService = "telegram"
const SlackService = "slack"

type SentryConfig struct {
	DSN                string  `env:"DSN"`
	Environment        string  `env:"ENVIRONMENT"`
	TracesSampleRate   float64 `env:"TRACES_SAMPLE_RATE" envDefault:"1.0"`
	ProfilesSampleRate float64 `env:"PROFILES_SAMPLE_RATE" envDefault:"1.0"`
}

type TelegramConfig struct {
	Token   string  `env:"TOKEN"`
	ChatIDs []int64 `env:"CHAT_IDS" envSeparator:","`
}

type SlackConfig struct {
	Token      string   `env:"TOKEN"`
	ChannelIDs []string `env:"CHANNEL_IDS" envSeparator:","`
}

type Config struct {
	Port      string         `env:"PORT" envDefault:"8080"`
	AuthToken string         `env:"AUTH_TOKEN,notEmpty"`
	Telegram  TelegramConfig `envPrefix:"TELEGRAM_"`
	Slack     SlackConfig    `envPrefix:"SLACK_"`
	Notifiers []string       `env:"NOTIFIERS" envSeparator:"," envDefault:"TELEGRAM"`
	Sentry    SentryConfig   `envPrefix:"SENTRY_"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	err := env.ParseWithOptions(&cfg, env.Options{
		RequiredIfNoDef: false,
	})
	if err != nil {
		return nil, err
	}

	err = checkNotifiers(cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func checkNotifiers(cfg Config) error {
	for i := range cfg.Notifiers {
		cfg.Notifiers[i] = strings.ToLower(cfg.Notifiers[i])
		switch cfg.Notifiers[i] {
		case TelegramService:
			if cfg.Telegram.Token == "" {
				return fmt.Errorf("empty telegram token")
			}
			if len(cfg.Telegram.ChatIDs) == 0 {
				return fmt.Errorf("no chat id provided")
			}
		case SlackService:
			if cfg.Slack.Token == "" {
				return fmt.Errorf("empty slack api token")
			}
			if len(cfg.Slack.ChannelIDs) == 0 {
				return fmt.Errorf("no channel id provided")
			}
		}
	}
	return nil
}

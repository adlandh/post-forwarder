// Package config contains application configuration
package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/caarlos0/env/v9"
)

const TelegramService = "telegram"
const SlackService = "slack"
const PushoverService = "pushover"

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

type PushoverConfig struct {
	Token string `env:"TOKEN"`
	User  string `env:"USER"`
}

type Config struct {
	Port      string         `env:"PORT" envDefault:"8080"`
	AuthToken string         `env:"AUTH_TOKEN,notEmpty"`
	Telegram  TelegramConfig `envPrefix:"TELEGRAM_"`
	Slack     SlackConfig    `envPrefix:"SLACK_"`
	Pushover  PushoverConfig `envPrefix:"PUSHOVER_"`
	Notifiers []string       `env:"NOTIFIERS" envSeparator:"," envDefault:"TELEGRAM"`
	Sentry    SentryConfig   `envPrefix:"SENTRY_"`
}

func NewConfig() (*Config, error) {
	var cfg Config

	if err := env.ParseWithOptions(&cfg, env.Options{
		RequiredIfNoDef: false,
	}); err != nil {
		return nil, fmt.Errorf("error parsing config: %w", err)
	}

	if err := checkNotifiers(cfg); err != nil {
		return nil, fmt.Errorf("error checking notifiers: %w", err)
	}

	return &cfg, nil
}

func checkNotifiers(cfg Config) (err error) {
	for i := range cfg.Notifiers {
		cfg.Notifiers[i] = strings.ToLower(cfg.Notifiers[i])
		switch cfg.Notifiers[i] {
		case TelegramService:
			err2 := checkTelegramConfig(cfg.Telegram)
			if err2 != nil {
				err = errors.Join(err, err2)
			}
		case SlackService:
			err2 := checkSlackConfig(cfg.Slack)
			if err2 != nil {
				err = errors.Join(err, err2)
			}
		case PushoverService:
			err2 := checkPushoverConfig(cfg.Pushover)
			if err2 != nil {
				err = errors.Join(err, err2)
			}
		}
	}

	return
}

func checkPushoverConfig(cfg PushoverConfig) (err error) {
	if cfg.Token == "" {
		err = fmt.Errorf("empty pushover application api token")
	}

	if cfg.User == "" {
		err = errors.Join(err, fmt.Errorf("empty pushover user api token"))
	}

	return
}

func checkSlackConfig(cfg SlackConfig) (err error) {
	if cfg.Token == "" {
		err = fmt.Errorf("empty slack api token")
	}

	if len(cfg.ChannelIDs) == 0 {
		err = errors.Join(err, fmt.Errorf("no channel id provided"))
	}

	return
}

func checkTelegramConfig(cfg TelegramConfig) (err error) {
	if cfg.Token == "" {
		err = fmt.Errorf("empty telegram token")
	}

	if len(cfg.ChatIDs) == 0 {
		err = errors.Join(err, fmt.Errorf("no chat id provided"))
	}

	return
}

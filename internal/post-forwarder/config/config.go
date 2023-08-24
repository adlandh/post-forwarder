package config

import "github.com/caarlos0/env/v9"

type Config struct {
	Port      string `env:"PORT" envDefault:"8080"`
	AuthToken string `env:"AUTH_TOKEN,notEmpty"`
	BotToken  string `env:"BOT_TOKEN,notEmpty"`
	BotChatID int64  `env:"BOT_CHAT_ID,notEmpty"`
	SentryDSN string `env:"SENTRY_DSN"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	err := env.ParseWithOptions(&cfg, env.Options{
		RequiredIfNoDef: false,
	})

	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

package config

import "github.com/caarlos0/env/v9"

type Config struct {
	Port                     string  `env:"PORT" envDefault:"8080"`
	AuthToken                string  `env:"AUTH_TOKEN,notEmpty"`
	BotToken                 string  `env:"BOT_TOKEN,notEmpty"`
	BotChatID                string  `env:"BOT_CHAT_ID,notEmpty"`
	SentryDSN                string  `env:"SENTRY_DSN"`
	SentryTracesSampleRate   float64 `env:"SENTRY_TRACES_SAMPLE_RATE" envDefault:"1.0"`
	SentryProfilesSampleRate float64 `env:"SENTRY_PROFILES_SAMPLE_RATE" envDefault:"1.0"`
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

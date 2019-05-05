package main

import (
	"time"

	"github.com/caarlos0/env"
)

type configuration struct {
	TelegramToken   string        `env:"TELEGRAM_TOKEN,required"`
	TelegramChannel string        `env:"TELEGRAM_CHANNEL,required"`
	RSSFeed         string        `env:"RSS_FEED,required"`
	RSSFeedPeriod   time.Duration `env:"RSS_FEED_PERIOD" envDefault:"5m"`
	DBPath          string        `env:"DB_PATH,required"`
}

func readConfig() (*configuration, error) {
	cfg := configuration{}
	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

package habrabot

import (
	"log"
	"sort"
	"time"

	"github.com/caarlos0/env"

	delivery2 "github.com/kapitanov/habrabot/internal/delivery"
	source2 "github.com/kapitanov/habrabot/internal/source"
	storage2 "github.com/kapitanov/habrabot/internal/storage"
)

type configuration struct {
	TelegramToken   string        `env:"TELEGRAM_TOKEN,required"`
	TelegramChannel string        `env:"TELEGRAM_CHANNEL,required"`
	RSSFeed         string        `env:"RSS_FEED,required"`
	RSSFeedPeriod   time.Duration `env:"RSS_FEED_PERIOD" envDefault:"5m"`
	RSSDBPath       string        `env:"RSS_DB_PATH,required"`
}

func readConfig() (*configuration, error) {
	cfg := configuration{}
	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func sync(feed source2.Feed, channel delivery2.Channel, db storage2.Driver) error {
	articles, err := feed.Read()
	if err != nil {
		return err
	}

	// TODO filter by tags (optional)

	sort.Slice(articles, func(i, j int) bool {
		return articles[i].Time.Before(articles[j].Time)
	})

	newArticleCount := 0
	for _, article := range articles {
		status, err := db.Store(article)
		if err != nil {
			return err
		}

		if status == storage2.New {
			log.Printf("New article from feed: %s", article.ID)
			newArticleCount++

			err := channel.Publish(article)
			if err != nil {
				return err
			}
		}
	}

	if newArticleCount > 0 {
		log.Printf("Sync completed, %d new article(s) were found", newArticleCount)
		log.Println("------------------")
	}

	return nil
}

// Main is an entrypoint for application.
func Main() {
	config, err := readConfig()
	if err != nil {
		log.Fatal(err)
	}

	channel, err := delivery2.NewTelegramChannel(config.TelegramToken, config.TelegramChannel)
	if err != nil {
		panic(err)
	}

	db, err := storage2.NewBoltDBDriver(config.RSSDBPath)
	if err != nil {
		log.Fatal(err)
	}

	feed := source2.NewRSSFeed(config.RSSFeed)

	for {
		err = sync(feed, channel, db)
		if err != nil {
			log.Fatal(err)
		}

		time.Sleep(config.RSSFeedPeriod)
	}
}

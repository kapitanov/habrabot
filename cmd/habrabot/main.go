package habrabot

import (
	"errors"
	"flag"
	"golang.org/x/net/context"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"

	"github.com/kapitanov/habrabot/internal/data"
	"github.com/kapitanov/habrabot/internal/opengraph"
	"github.com/kapitanov/habrabot/internal/rss"
	"github.com/kapitanov/habrabot/internal/storage"
	"github.com/kapitanov/habrabot/internal/telegram"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

var (
	envFilePath *string
)

func init() {
	envFilePath = flag.String("env", "", "path to .env file to load")
}

type configuration struct {
	TelegramToken   string        `env:"TELEGRAM_TOKEN,required"`
	TelegramChannel string        `env:"TELEGRAM_CHANNEL,required"`
	RSSFeedURL      string        `env:"RSS_FEED,required"`
	RSSFeedPeriod   time.Duration `env:"RSS_FEED_PERIOD" envDefault:"5m"`
	BoltDBPath      string        `env:"BOLTDB_PATH,required"`
}

func readConfig() (configuration, error) {
	if *envFilePath != "" {
		err := godotenv.Load(*envFilePath)
		if err != nil {
			return configuration{}, err
		}
	}

	cfg := configuration{}
	err := env.Parse(&cfg)
	if err != nil {
		return configuration{}, err
	}

	return cfg, nil
}

func (c configuration) CreateFeed() (data.Feed, error) {
	// RSS feed is a root feed.
	feed, err := rss.New(c.RSSFeedURL)
	if err != nil {
		return nil, err
	}

	// Then it should be wrapped into opengraph enricher.
	feed = opengraph.Enrich(feed)

	// Then it should be filtered by BoltDB database.
	feed = storage.UseBoltDB(feed, c.BoltDBPath)

	return feed, nil
}

func (c configuration) CreateConsumer() data.Consumer {
	return telegram.New(c.TelegramToken, c.TelegramChannel)
}

func runOnce(ctx context.Context, feed data.Feed, consumer data.Consumer) error {
	newArticleCount := 0
	feed = data.Transform(feed, data.TransformationFunc(func(_ context.Context, article *data.Article) error {
		log.Printf("new article from feed: %s", article.ID)
		newArticleCount++

		return nil
	}))

	err := feed.Read(ctx, consumer)
	if err != nil {
		return err
	}

	if newArticleCount > 0 {
		log.Info().Int("new", newArticleCount).Msg("sync completed")
	}

	return nil
}

// Main is an entrypoint for application.
func Main() {
	flag.Parse()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out: os.Stderr,
	})
	log.Logger = log.Logger.With().Timestamp().Logger()

	config, err := readConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("unable to load configuration")
	}

	feed, err := config.CreateFeed()
	if err != nil {
		log.Fatal().Err(err).Msg("unable to create feed")
	}

	consumer := config.CreateConsumer()

	run(feed, consumer, config)
}

func run(feed data.Feed, consumer data.Consumer, config configuration) {
	syncTrigger := make(chan struct{}, 1)
	syncTrigger <- struct{}{}

	ctx, cancel := context.WithCancel(context.Background())

	timer := time.NewTicker(config.RSSFeedPeriod)
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		for range timer.C {
			syncTrigger <- struct{}{}
		}
	}()

	go func() {
		defer wg.Done()

		for range syncTrigger {
			err := runOnce(ctx, feed, consumer)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					return
				}

				log.Fatal().Err(err).Msg("unable to run sync routine")
			}
		}
	}()

	log.Info().Msg("shutting down")
	close(syncTrigger)
	cancel()
	timer.Stop()
	wg.Wait()

	log.Info().Msg("goodbye")
}

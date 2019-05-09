package main

import (
	"log"
	"sort"
	"time"

	"github.com/kapitanov/habrabot/delivery"
	"github.com/kapitanov/habrabot/source"
	"github.com/kapitanov/habrabot/storage"
)

func sync(feed source.Feed, channel delivery.Channel, db storage.Driver) error {
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

		if status == storage.StatusNew {
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

func main() {
	config, err := readConfig()
	if err != nil {
		log.Fatal(err)
	}

	channel, err := delivery.NewTelegramChannel(config.TelegramToken, config.TelegramChannel)
	if err != nil {
		panic(err)
	}

	db, err := storage.NewBoltDBDriver(config.RSSDBPath)
	if err != nil {
		log.Fatal(err)
	}

	feed := source.NewRSSFeed(config.RSSFeed)

	for {
		err = sync(feed, channel, db)
		if err != nil {
			log.Fatal(err)
		}

		time.Sleep(config.RSSFeedPeriod)
	}
}

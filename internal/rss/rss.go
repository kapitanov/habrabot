package rss

import (
	"log"
	"net/url"
	"sort"
	"strings"

	"github.com/mmcdole/gofeed"

	"github.com/kapitanov/habrabot/internal/data"
)

type feed struct {
	URL string
}

// New creates new RSS feed reader.
func New(url string) data.Feed {
	log.Printf("will read news from rss feed \"%s\"", url)
	return &feed{url}
}

// Read method reads feed items and streams them into the consumer.
func (r *feed) Read(consumer data.Consumer) error {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(r.URL)
	if err != nil {
		return err
	}

	var articles []data.Article
	for _, item := range feed.Items {
		article, err := parseArticleFromRss(item)
		if err != nil {
			return err
		}

		articles = append(articles, article)
	}

	// Sort articles by time in ascending order.
	sort.Slice(articles, func(i, j int) bool {
		return articles[i].Time.Before(articles[j].Time)
	})

	for _, article := range articles {
		err := consumer.On(article)
		if err != nil {
			return err
		}
	}

	return nil
}

func parseArticleFromRss(item *gofeed.Item) (data.Article, error) {
	description, err := normalizeHTML(item.Description)
	if err != nil {
		return data.Article{}, err
	}

	imageURL, err := extractImageURL(item.Description)
	if err != nil {
		return data.Article{}, err
	}

	u, err := url.Parse(item.Link)
	if err != nil {
		return data.Article{}, err
	}

	u.RawQuery = ""

	article := data.Article{
		ID:          item.GUID,
		Time:        *item.PublishedParsed,
		Title:       item.Title,
		LinkURL:     u.String(),
		Description: description,
		Author:      "",
		ImageURL:    imageURL,
	}

	if item.Author != nil {
		article.Author = item.Author.Name
	}

	tags := make(map[string]string)
	for _, cat := range item.Categories {
		cat = strings.ToLower(cat)
		tags[cat] = cat
	}

	article.Tags = make([]string, len(tags))
	i := 0
	for tag := range tags {
		article.Tags[i] = tag
		i++
	}

	return article, nil
}

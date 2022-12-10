package rss

import (
	"context"
	"net/url"
	"sort"
	"strings"

	"github.com/hashicorp/go-retryablehttp"

	"github.com/kapitanov/habrabot/internal/httpclient"

	"github.com/mmcdole/gofeed"
	"github.com/rs/zerolog/log"

	"github.com/kapitanov/habrabot/internal/data"
)

type feed struct {
	URL        string
	HTTPClient *retryablehttp.Client
}

// New creates new RSS feed reader.
func New(url string) (data.Feed, error) {
	log.Info().Str("url", url).Msg("using rss feed")

	httpClient, err := httpclient.New(httpclient.RSSPolicy)
	if err != nil {
		return nil, err
	}

	return &feed{
		URL:        url,
		HTTPClient: httpClient,
	}, nil
}

// Read method reads feed items and streams them into the consumer.
func (r *feed) Read(ctx context.Context, consumer data.Consumer) error {
	fp := gofeed.NewParser()
	fp.Client = r.HTTPClient.HTTPClient

	feed, err := fp.ParseURL(r.URL)
	if err != nil {
		log.Error().Err(err).Str("url", r.URL).Msg("unable to parse rss url")
		return err
	}

	var articles []data.Article
	for _, item := range feed.Items {
		if isCanceled(ctx) {
			return context.Canceled
		}

		article, err := parseArticleFromRSS(item)
		if err != nil {
			log.Error().Err(err).Str("url", r.URL).Msg("unable to item from rss feed")
			return err
		}

		articles = append(articles, article)
	}

	// Sort articles by time in ascending order.
	sort.Slice(articles, func(i, j int) bool {
		return articles[i].Time.Before(articles[j].Time)
	})

	for _, article := range articles {
		err := consumer.On(ctx, article)
		if err != nil {
			return err
		}
	}

	return nil
}

func parseArticleFromRSS(item *gofeed.Item) (data.Article, error) {
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

func isCanceled(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

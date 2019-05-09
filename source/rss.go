package source

import (
	"log"
	"net/url"
	"strings"

	"github.com/mmcdole/gofeed"
)

type rssFeed struct {
	URL string
}

// NewRSSFeed creates new RSS feed reader
func NewRSSFeed(url string) Feed {
	log.Printf("Will read news from RSS feed \"%s\"", url)
	return &rssFeed{url}
}

// Reads a list of feed items from RSS
func (r *rssFeed) Read() ([]*Article, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(r.URL)
	if err != nil {
		return nil, err
	}

	articles := make([]*Article, len(feed.Items))
	for i, item := range feed.Items {
		article, err := parseArticleFromRss(item)
		if err != nil {
			return nil, err
		}

		articles[i] = article
	}

	return articles, nil
}

func parseArticleFromRss(item *gofeed.Item) (*Article, error) {
	description, err := normalizeHTML(item.Description)
	if err != nil {
		return nil, err
	}

	imageURL, err := extractImageURL(item.Description)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(item.Link)
	if err != nil {
		return nil, err
	}

	u.RawQuery = ""

	article := &Article{
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

package source

import (
	"log"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"

	"github.com/mmcdole/gofeed"
)

type rssFeed struct {
	URL string
}

// NewRSSFeed creates new RSS feed reader.
func NewRSSFeed(url string) Feed {
	log.Printf("Will read news from RSS feed \"%s\"", url)
	return &rssFeed{url}
}

// Read method reads a list of feed items from RSS.
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

	// Load web page and try parse OpenGraph tags
	err = loadOpengraphTags(article)
	if err != nil {
		return nil, err
	}

	return article, nil
}

func loadOpengraphTags(article *Article) error {
	resp, err := http.Get(article.LinkURL)
	if err != nil {
		return err
	}

	if resp.StatusCode >= http.StatusMultipleChoices {
		log.Printf("Unable to download \"%s\": %d", article.LinkURL, resp.StatusCode)
		return nil
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	root, err := html.Parse(resp.Body)
	if err != nil {
		return err
	}

	searchForOpengraphTags(article, root, false, false)

	return nil
}

func searchForOpengraphTags(article *Article, node *html.Node, hasTitleTag, hasImageTag bool) {
	isMetaTag, property, content := tryParseMetaTag(node)

	if isMetaTag {
		switch property {
		case "og:title":
			if !hasTitleTag {
				article.Title = content
			}
		case "og:image":
			if !hasImageTag {
				article.ImageURL = content
			}
		}
		return
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		searchForOpengraphTags(article, c, hasTitleTag, hasImageTag)
	}
}

func tryParseMetaTag(node *html.Node) (bool, string, string) {
	if node.Type == html.ElementNode && node.Data == "meta" {
		property := ""
		content := ""

		for _, attr := range node.Attr {
			switch attr.Key {
			case "property":
				property = attr.Val

			case "content":
				content = attr.Val
			}
		}

		if property != "" && content != "" {
			return true, property, content
		}
	}

	return false, "", ""
}

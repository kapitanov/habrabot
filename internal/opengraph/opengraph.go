package opengraph

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
	"golang.org/x/net/html"

	"github.com/kapitanov/habrabot/internal/data"
)

// Enrich adds opengraph data into the stream of articles.
func Enrich(feed data.Feed) data.Feed {
	return data.Transform(feed, data.TransformationFunc(func(article *data.Article) error {
		// Load web page and try parse OpenGraph tags
		t, err := loadTags(article.LinkURL)
		if err == nil {
			// Errors are ignored here
			t.Enrich(article)
		}
		return nil
	}))
}

type tags struct {
	Title    *string
	ImageURL *string
}

func (t tags) Enrich(article *data.Article) {
	if t.Title != nil {
		article.Title = *t.Title
	}

	if t.ImageURL != nil {
		article.ImageURL = t.ImageURL
	}
}

func loadTags(sourceURL string) (tags, error) {
	//nolint:gosec // Suppress "G107: Potential HTTP request made with variable url"
	resp, err := http.Get(sourceURL)
	if err != nil {
		log.Warn().Err(err).Str("url", sourceURL).Msg("unable to download web page")
		return tags{}, err
	}

	if resp.StatusCode >= http.StatusMultipleChoices {
		log.Warn().Err(err).
			Str("url", sourceURL).
			Int("status", resp.StatusCode).
			Msg("unable to download web page")
		return tags{}, fmt.Errorf("unable to download \"%s\": %v", sourceURL, resp.Status)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	root, err := html.Parse(resp.Body)
	if err != nil {
		log.Warn().Err(err).Str("url", sourceURL).Msg("unable to parse web page")
		return tags{}, err
	}

	t := parseTags(root)
	return t, nil
}

func parseTags(root *html.Node) tags {
	t := tags{}

	// Find <html> node
	htmlNode := findNode(root, "html")
	if htmlNode != nil {
		// Find <head> node
		headNode := findNode(htmlNode, "head")

		// Scan its children looking for <meta>
		if headNode != nil {
			for node := headNode.FirstChild; node != nil; node = node.NextSibling {
				if node.Type != html.ElementNode || node.Data != "meta" {
					continue
				}

				processMetaTag(node, &t)
			}
		}
	}

	return t
}

func findNode(root *html.Node, name string) *html.Node {
	queue := []*html.Node{root}

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		if node.Type == html.ElementNode && node.Data == name {
			return node
		}

		for c := node.FirstChild; c != nil; c = c.NextSibling {
			queue = append(queue, c)
		}
	}

	return nil
}

func processMetaTag(node *html.Node, t *tags) {
	key, value := decodeMetaTag(node)

	switch key {
	case "og:title":
		t.Title = &value
	case "og:image":
		t.ImageURL = &value
	case "og:description":

		_ = value
	}
}

func decodeMetaTag(node *html.Node) (key, value string) {
	for _, attr := range node.Attr {
		switch attr.Key {
		case "property":
			key = attr.Val
		case "content":
			value = attr.Val
		}
	}

	return
}

package opengraph

import (
	"log"
	"net/http"

	"golang.org/x/net/html"

	"github.com/kapitanov/habrabot/internal/data"
)

// Enrich adds opengraph data into the stream of articles.
func Enrich(feed data.Feed) data.Feed {
	return data.Transform(feed, data.TransformationFunc(func(article *data.Article) error {
		// Load web page and try parse OpenGraph tags
		loadOpengraphTags(article)
		return nil
	}))
}

func loadOpengraphTags(article *data.Article) {
	resp, err := http.Get(article.LinkURL)
	if err != nil {
		log.Printf("unable to download \"%s\": %v", article.LinkURL, err)
		return
	}

	if resp.StatusCode >= http.StatusMultipleChoices {
		log.Printf("unable to download \"%s\": %v", article.LinkURL, resp.Status)
		return
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	root, err := html.Parse(resp.Body)
	if err != nil {
		log.Printf("unable to process \"%s\": %v", article.LinkURL, err)
		return
	}

	searchForOpengraphTags(article, root, false, false)
}

func searchForOpengraphTags(article *data.Article, node *html.Node, hasTitleTag, hasImageTag bool) {
	isMetaTag, property, content := tryParseMetaTag(node)

	if isMetaTag {
		switch property {
		case "og:title":
			if !hasTitleTag {
				article.Title = content
			}
		case "og:image":
			if !hasImageTag {
				article.ImageURL = &content
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

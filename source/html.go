package source

import (
	"net/url"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

func normalizeHTML(input string) (string, error) {
	if input == "" {
		return "", nil
	}

	node, err := html.Parse(strings.NewReader(input))
	if err != nil {
		return "", err
	}

	text := extractHTMLNodeText(node, "")
	text = strings.Trim(text, " \n")
	text = replaceRegexp(text, "[\\n\\s]+\\n\\s*", "\n")
	text = replaceRegexp(text, "  ", " ")
	return text, nil
}

func replaceRegexp(text, regex, replace string) string {
	r := regexp.MustCompile(regex)
	text = r.ReplaceAllString(text, replace)
	return text
}

func extractHTMLNodeText(node *html.Node, text string) string {
	if node.Type == html.TextNode {
		return text + strings.Trim(node.Data, "\n")
	}

	if node.Type == html.ElementNode {
		switch node.Data {
		case "br":
			if len(text) > 0 {
				text = text + "\n"
			}
			break

		case "a":
			if !shouldExtractAnchorText(node) {
				return text
			}
			break
		}
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		text = extractHTMLNodeText(c, text)
	}

	return text
}

func shouldExtractAnchorText(node *html.Node) bool {
	// Skip hyperlinks wuth utm_source query parameter
	for _, a := range node.Attr {
		if a.Key == "href" {
			u, err := url.Parse(a.Val)
			if err != nil {
				return false
			}

			_, exists := u.Query()["utm_source"]
			if exists {
				return false
			}
		}
	}

	return true
}

func extractImageURL(input string) (string, error) {
	if input == "" {
		return "", nil
	}

	nodes, err := html.ParseFragment(strings.NewReader(input), nil)
	if err != nil {
		return "", err
	}

	for _, node := range nodes {
		url := extractHTMLNodeImageURL(node)
		if url != "" {
			return url, nil
		}
	}

	return "", nil
}

func extractHTMLNodeImageURL(node *html.Node) string {
	if node.Type == html.ElementNode && node.Data == "img" {
		for _, a := range node.Attr {
			if a.Key == "src" && a.Val != "" {
				return a.Val
			}
		}
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		url := extractHTMLNodeImageURL(c)
		if url != "" {
			return url
		}
	}

	return ""
}

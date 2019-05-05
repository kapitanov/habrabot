package source

import (
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

	return text, nil
}

func extractHTMLNodeText(node *html.Node, text string) string {
	if node.Type == html.TextNode {
		return text + strings.Trim(node.Data, "\n")
	}

	if node.Type == html.ElementNode {
		if node.Data == "br" {
			if len(text) > 0 {
				text = text + "\n"
			}

			return text
		}

		if node.Data == "a" {
			return text
		}
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		text = extractHTMLNodeText(c, text)
	}

	return text
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

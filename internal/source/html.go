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

	text := extractHTMLNodeText(node)
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

//nolint:cyclop // Will refactor later
func extractHTMLNodeText(node *html.Node) string {
	// Text nodes
	if node.Type == html.TextNode {
		return strings.Trim(node.Data, "\n")
	}

	// Special handling for <br>
	if node.Type == html.ElementNode && node.Data == "br" {
		return "\n"
	}

	// Wrapping tag detection
	wrappingTag := ""
	textPrefix := ""
	addFinalNewline := false
	if node.Type == html.ElementNode {
		switch node.Data {
		case "a":
			// UTM links should be removed
			if isUtmHyperlink(node) {
				return ""
			}

		case "b", "strong", "i", "em", "code", "s", "strike", "del", "u":
			// Supported tags
			wrappingTag = node.Data

		case "pre":
			// Special handling for <pre>
			return extractPreHTMLNodeText(node)

		case "p":
			addFinalNewline = true

		case "li":
			textPrefix = "- "
			addFinalNewline = true
		}
	}

	// Opening tag
	text := ""
	if wrappingTag != "" {
		text += "<" + wrappingTag + ">"
	}
	if textPrefix != "" {
		text += textPrefix
	}

	// Extract child content
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		text += extractHTMLNodeText(c)
	}

	// Closing tag
	if wrappingTag != "" {
		text += "</" + wrappingTag + ">"
	}

	if addFinalNewline {
		text += "\n"
	}

	return text
}

func extractPreHTMLNodeText(node *html.Node) string {
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "code" {
			return extractPreWithNestedCodeNodeText(c)
		}
	}

	return extractPreWithoutNestedCodeNodeText(node)
}

func extractPreWithoutNestedCodeNodeText(node *html.Node) string {
	lang := ""
	for _, attr := range node.Attr {
		if attr.Key == "class" {
			lang = attr.Val
		}
	}

	var text string
	if lang != "" {
		text = "<pre language=\"" + html.EscapeString(lang) + "\">"
	} else {
		text = "<pre>"
	}

	// Extract child content
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			text += strings.Trim(c.Data, "\n")
			text += "\n"
		}
	}

	text += "</pre>\n"
	return text
}

func extractPreWithNestedCodeNodeText(node *html.Node) string {
	lang := ""
	for _, attr := range node.Attr {
		if attr.Key == "class" {
			lang = attr.Val
		}
	}

	var text string
	if lang != "" {
		text = "<pre language=\"" + html.EscapeString(lang) + "\">"
	} else {
		text = "<pre>"
	}

	// Extract child content
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			text += strings.Trim(c.Data, "\n")
			text += "\n"
		}
	}

	text += "</pre>\n"
	return text
}

func isUtmHyperlink(node *html.Node) bool {
	// Skip hyperlinks wuth utm_source query parameter
	for _, a := range node.Attr {
		if a.Key == "href" {
			u, err := url.Parse(a.Val)
			if err != nil {
				return true
			}

			_, exists := u.Query()["utm_source"]
			if exists {
				return true
			}
		}
	}

	return false
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
		u := extractHTMLNodeImageURL(node)
		if u != "" {
			return u, nil
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
		u := extractHTMLNodeImageURL(c)
		if u != "" {
			return u
		}
	}

	return ""
}

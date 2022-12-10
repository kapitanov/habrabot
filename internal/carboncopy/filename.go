package carboncopy

import (
	"github.com/kapitanov/habrabot/internal/data"
	"regexp"
	"strings"
	"unicode"
)

var articleIDRegex = regexp.MustCompile("([0-9]{3,})")

func extractFileName(article data.Article) string {
	var filename []rune

	articleID := articleIDRegex.FindString(article.LinkURL)
	for _, ch := range articleID {
		filename = append(filename, ch)
	}

	hasPrefix := false

	for _, ch := range strings.TrimSpace(article.Title) {
		if isAllowedRune(ch) {
			if !hasPrefix {
				hasPrefix = true

				if len(filename) > 0 {
					filename = append(filename, ' ', '-', ' ')
				}
			}

			if !unicode.IsSpace(ch) {
				filename = append(filename, ch)
			} else {
				filename = append(filename, ' ')
			}
		}
	}

	str := string(filename)
	str = strings.TrimRight(str, ". ")
	str += ".html"
	return str
}

func isAllowedRune(ch rune) bool {
	if unicode.IsLetter(ch) || unicode.IsDigit(ch) || unicode.IsSpace(ch) {
		return true
	}

	switch ch {
	case '-', ',', '.', '_':
		return true
	default:
		return false
	}
}

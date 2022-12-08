package source

import (
	"time"
)

// Article is a single item of a feed.
type Article struct {
	ID          string    // ID of article.
	Title       string    // Title of article.
	Time        time.Time // An article's publication time.
	Description string    // An article's description text.
	LinkURL     string    // A hyperlink to article's web page.
	ImageURL    string    // A hyperlink to article's title image if available.
	Author      string    // Article's author name.
	Tags        []string  // List of article's tags.
}

// Feed reads article list from remote source.
type Feed interface {
	// Read methods reads a list of feed items.
	Read() ([]*Article, error)
}

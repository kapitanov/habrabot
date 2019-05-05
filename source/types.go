package source

import (
	"time"
)

// Article is a single item of a feed
type Article struct {
	ID          string
	Title       string
	Time        time.Time
	Description string
	LinkURL     string
	ImageURL    string
	Author      string
	Tags        []string
}

// Feed reads article list from remote source
type Feed interface {
	// Reads a list of feed items
	Read() ([]*Article, error)
}

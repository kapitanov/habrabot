package data

import (
	"context"
	"time"
)

// Article is a single item of a feed.
type Article struct {
	ID          string    // ID of article.
	Title       string    // Title of article.
	Time        time.Time // An article's publication time.
	Description string    // An article's description text.
	LinkURL     string    // A hyperlink to article's web page.
	ImageURL    *string   // A hyperlink to article's title image if available.
	Author      string    // Article's author name.
	Tags        []string  // List of article's tags.
}

// Feed reads article list from remote source.
type Feed interface {
	// Read method reads feed items and streams them into the consumer.
	Read(ctx context.Context, consumer Consumer) error
}

// Consumer consumes feed items.
type Consumer interface {
	// On method is invoked when an article is received from the feed.
	On(ctx context.Context, article Article) error
}

// ConsumerFunc is a function-based implementation of Consumer.
type ConsumerFunc func(ctx context.Context, article Article) error

// On method is invoked when an article is received from the feed.
func (c ConsumerFunc) On(ctx context.Context, article Article) error {
	return c(ctx, article)
}

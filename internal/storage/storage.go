package storage

import "github.com/kapitanov/habrabot/internal/data"

// Storage acts as a filter for feed stream.
type Storage interface {
	// Filter creates new filtered feed.
	Filter(feed data.Feed) data.Feed
}

package storage

import (
	"github.com/kapitanov/habrabot/internal/source"
)

// Status is an article storage status.
type Status int16

const (
	// New means an article has not been stored yet.
	New Status = iota

	// Old means an article has been stored already.
	Old
)

// Driver defines methods to access the local database.
type Driver interface {
	// Store tries to write an article and returns storage status.
	Store(article *source.Article) (Status, error)
}

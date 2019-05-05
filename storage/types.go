package storage

import "github.com/kapitanov/habrabot/source"

// StorageStatus is an article storage status
type StorageStatus int16

const (
	// StatusNew means an article has not been stored yet
	StatusNew StorageStatus = iota
	// StatusNew means an article has been stored already
	StatusOld
)

// Driver defines methods to access the local database
type Driver interface {
	// Tries to write an article and returns storage status
	Store(article *source.Article) (StorageStatus, error)
}

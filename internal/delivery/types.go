package delivery

import (
	"github.com/kapitanov/habrabot/internal/source"
)

// Channel provides an entrypoint for message delivery.
type Channel interface {
	// Publish publishes a message.
	Publish(article *source.Article) error
}

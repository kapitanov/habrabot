package delivery

import "github.com/kapitanov/habrabot/source"

type Channel interface {
	Publish(article *source.Article) error
}

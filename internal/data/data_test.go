package data

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func NewArticle(id string) Article {
	return Article{
		ID: id,
	}
}

func NewArticles(ids ...string) []Article {
	var articles []Article
	for _, id := range ids {
		articles = append(articles, NewArticle(id))
	}
	return articles
}

func NewInMemoryFeed(articles []Article) Feed {
	return InMemoryFeed(articles)
}

type InMemoryFeed []Article

func (f InMemoryFeed) Read(consumer Consumer) error {
	for _, article := range f {
		err := consumer.On(article)
		if err != nil {
			return err
		}
	}

	return nil
}

type Pipeline func(feed Feed) Feed

func RunFeed(
	t *testing.T,
	input []Article,
	pipeline Pipeline,
	consumer Consumer,
) error {
	feed := NewInMemoryFeed(input)
	feed = pipeline(feed)
	require.NotNil(t, feed)

	return feed.Read(consumer)
}

func RunFeedInMemory(
	t *testing.T,
	input []Article,
	pipeline Pipeline,
) ([]Article, error) {
	var output []Article
	consumer := ConsumerFunc(func(article Article) error {
		output = append(output, article)
		return nil
	})

	err := RunFeed(t, input, pipeline, consumer)
	return output, err
}

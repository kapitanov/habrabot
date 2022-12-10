package data

import (
	"context"
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

func (f InMemoryFeed) Read(ctx context.Context, consumer Consumer) error {
	for _, article := range f {
		err := consumer.On(ctx, article)
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
	if pipeline != nil {
		feed = pipeline(feed)
	}
	require.NotNil(t, feed)

	return feed.Read(context.Background(), consumer)
}

func RunFeedInMemory(
	t *testing.T,
	input []Article,
	pipeline Pipeline,
) ([]Article, error) {
	var output []Article
	consumer := ConsumerFunc(func(_ context.Context, article Article) error {
		output = append(output, article)
		return nil
	})

	err := RunFeed(t, input, pipeline, consumer)
	return output, err
}

type InMemoryConsumer struct {
	articles []Article
}

func NewInMemoryConsumer() *InMemoryConsumer {
	return &InMemoryConsumer{}
}

func (c *InMemoryConsumer) Items() []Article {
	return c.articles
}

func (c *InMemoryConsumer) On(_ context.Context, article Article) error {
	c.articles = append(c.articles, article)
	return nil
}

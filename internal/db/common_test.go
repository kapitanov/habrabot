package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kapitanov/habrabot/internal/data"
)

func NewArticle(id string) data.Article {
	return data.Article{
		ID: id,
	}
}

func NewArticles(ids ...string) []data.Article {
	var articles []data.Article
	for _, id := range ids {
		articles = append(articles, NewArticle(id))
	}
	return articles
}

func NewInMemoryFeed(articles []data.Article) data.Feed {
	return InMemoryFeed(articles)
}

type InMemoryFeed []data.Article

func (f InMemoryFeed) Read(ctx context.Context, consumer data.Consumer) error {
	for _, article := range f {
		err := consumer.On(ctx, article)
		if err != nil {
			return err
		}
	}

	return nil
}

func Execute(t *testing.T, feed data.Feed) []data.Article {
	var articles []data.Article
	var consumer data.ConsumerFunc = func(_ context.Context, article data.Article) error {
		articles = append(articles, article)
		return nil
	}

	err := feed.Read(context.Background(), consumer)
	require.NoError(t, err)
	return articles
}

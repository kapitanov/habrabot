package data

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrap_NoOp(t *testing.T) {
	var pipeline Pipeline = func(feed Feed) Feed {
		return Wrap(feed, MiddlewareFunc(func(_ context.Context, _ Article, _ NextFunc) error {
			return nil
		}))
	}
	input := NewArticles("1", "2", "3")
	output, err := RunFeedInMemory(t, input, pipeline)

	assert.NoError(t, err)
	assert.Empty(t, output)
}

func TestWrap_MiddlewareError(t *testing.T) {
	expectedError := errors.New("expected error")

	var pipeline Pipeline = func(feed Feed) Feed {
		return Wrap(feed, MiddlewareFunc(func(_ context.Context, _ Article, _ NextFunc) error {
			return expectedError
		}))
	}
	input := NewArticles("1", "2", "3")
	output, err := RunFeedInMemory(t, input, pipeline)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Empty(t, output)
}

func TestWrap_ConsumerError(t *testing.T) {
	expectedError := errors.New("expected error")

	var pipeline Pipeline = func(feed Feed) Feed {
		return Wrap(feed, MiddlewareFunc(func(_ context.Context, article Article, next NextFunc) error {
			return next(article)
		}))
	}
	var consumer ConsumerFunc = func(_ context.Context, _ Article) error {
		return expectedError
	}
	input := NewArticles("1", "2", "3")
	err := RunFeed(t, input, pipeline, consumer)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
}

func TestWrap_Success(t *testing.T) {
	beforeCount, afterCount := 0, 0

	var pipeline Pipeline = func(feed Feed) Feed {
		return Wrap(feed, MiddlewareFunc(func(_ context.Context, article Article, next NextFunc) error {
			beforeCount++

			err := next(article)
			if err != nil {
				return err
			}

			afterCount++
			return nil
		}))
	}
	input := NewArticles("1", "2", "3")
	output, err := RunFeedInMemory(t, input, pipeline)

	assert.NoError(t, err)
	assert.Len(t, output, len(input))
	assert.Equal(t, beforeCount, len(input))
	assert.Equal(t, afterCount, len(input))
	for i := range input {
		assert.Equal(t, input[i], output[i])
	}
}

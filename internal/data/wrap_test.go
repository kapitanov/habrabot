package data

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrap_NoOp(t *testing.T) {
	var pipeline Pipeline = func(feed Feed) Feed {
		return Wrap(feed, MiddlewareFunc(func(_ Article, _ func(article Article) error) error {
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
		return Wrap(feed, MiddlewareFunc(func(_ Article, _ func(article Article) error) error {
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
		return Wrap(feed, MiddlewareFunc(func(article Article, next func(article Article) error) error {
			return next(article)
		}))
	}
	var consumer ConsumerFunc = func(_ Article) error {
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
		return Wrap(feed, MiddlewareFunc(func(article Article, next func(article Article) error) error {
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

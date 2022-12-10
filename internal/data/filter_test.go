package data

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilter_Block(t *testing.T) {
	var pipeline Pipeline = func(feed Feed) Feed {
		return Filter(feed, PredicateFunc(func(_ context.Context, _ Article) (bool, error) {
			return false, nil
		}))
	}
	input := NewArticles("1", "2", "3")
	output, err := RunFeedInMemory(t, input, pipeline)

	assert.NoError(t, err)
	assert.Empty(t, output)
}

func TestFilter_Pass(t *testing.T) {
	var pipeline Pipeline = func(feed Feed) Feed {
		return Filter(feed, PredicateFunc(func(_ context.Context, _ Article) (bool, error) {
			return true, nil
		}))
	}
	input := NewArticles("1", "2", "3")
	output, err := RunFeedInMemory(t, input, pipeline)

	assert.NoError(t, err)
	assert.Len(t, output, len(input))
	for i := range input {
		assert.Equal(t, input[i], output[i])
	}
}

func TestFilter_PassAndBlock(t *testing.T) {
	input := NewArticles("1", "2", "3", "4")

	var pipeline Pipeline = func(feed Feed) Feed {
		return Filter(feed, PredicateFunc(func(_ context.Context, article Article) (bool, error) {
			if article.ID == input[0].ID || article.ID == input[2].ID {
				return true, nil
			}

			return false, nil
		}))
	}
	output, err := RunFeedInMemory(t, input, pipeline)

	assert.NoError(t, err)
	if assert.Len(t, output, 2) {
		assert.Equal(t, input[0], output[0])
		assert.Equal(t, input[2], output[1])
	}
}

func TestFilter_Error(t *testing.T) {
	expectedError := errors.New("expected error")

	var pipeline Pipeline = func(feed Feed) Feed {
		return Filter(feed, PredicateFunc(func(_ context.Context, _ Article) (bool, error) {
			return false, expectedError
		}))
	}
	input := NewArticles("1", "2", "3")
	output, err := RunFeedInMemory(t, input, pipeline)

	assert.Error(t, err)
	assert.Equal(t, err, expectedError)
	assert.Empty(t, output, len(input))
}

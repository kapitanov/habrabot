package data

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTransform_NoOp(t *testing.T) {
	var pipeline Pipeline = func(feed Feed) Feed {
		return Transform(feed, TransformationFunc(func(article *Article) error {
			return nil
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

func TestTransform_Error(t *testing.T) {
	expectedError := errors.New("expected error")

	var pipeline Pipeline = func(feed Feed) Feed {
		return Transform(feed, TransformationFunc(func(article *Article) error {
			return expectedError
		}))
	}
	input := NewArticles("1", "2", "3")
	output, err := RunFeedInMemory(t, input, pipeline)

	assert.Error(t, err)
	assert.Equal(t, err, expectedError)
	assert.Empty(t, output)
}

func TestTransform_Modify(t *testing.T) {
	var pipeline Pipeline = func(feed Feed) Feed {
		return Transform(feed, TransformationFunc(func(article *Article) error {
			article.ID = fmt.Sprintf("NEW:%v", article.ID)
			return nil
		}))
	}

	input := NewArticles("1", "2", "3")
	output, err := RunFeedInMemory(t, input, pipeline)

	assert.NoError(t, err)
	assert.Len(t, output, len(input))
	for i := range input {
		assert.Contains(t, output[i].ID, input[i].ID)
		assert.Contains(t, output[i].ID, "NEW:")
	}
}

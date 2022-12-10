package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTee(t *testing.T) {
	input := NewArticles("1", "2", "3")

	c1 := NewInMemoryConsumer()
	c2 := NewInMemoryConsumer()
	c3 := NewInMemoryConsumer()

	c := Tee(c1, c2, c3)

	err := RunFeed(t, input, nil, c)

	assert.NoError(t, err)

	assert.Len(t, c1.Items(), len(input), "c1")
	assert.Len(t, c2.Items(), len(input), "c2")
	assert.Len(t, c3.Items(), len(input), "c3")

	for i := range input {
		assert.Equal(t, input[i], c1.Items()[i], "c1", i)
		assert.Equal(t, input[i], c2.Items()[i], "c2", i)
		assert.Equal(t, input[i], c3.Items()[i], "c3", i)
	}
}

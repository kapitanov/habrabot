package db

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUseBoltDB_Pass(t *testing.T) {
	f, err := os.CreateTemp(os.TempDir(), "*")
	require.NoError(t, err)
	dbPath := f.Name()
	t.Logf("data file: %v", dbPath)
	err = f.Close()
	require.NoError(t, err)

	defer func() {
		_ = os.Remove(dbPath)
	}()

	input := NewArticles("1", "2", "3")
	feed := NewInMemoryFeed(input)
	feed = Use(feed, dbPath)

	output := Execute(t, feed)

	assert.Len(t, output, len(input))
	for i := range output {
		assert.Equal(t, input[i], output[i])
	}
}

func TestUseBoltDB_Block(t *testing.T) {
	f, err := os.CreateTemp(os.TempDir(), "*")
	require.NoError(t, err)
	dbPath := f.Name()
	t.Logf("data file: %v", dbPath)
	err = f.Close()
	require.NoError(t, err)

	defer func() {
		_ = os.Remove(dbPath)
	}()

	input := NewArticles("1", "2", "3")

	feed := NewInMemoryFeed(append(input, input...))
	feed = Use(feed, dbPath)

	output := Execute(t, feed)

	assert.Len(t, output, len(input))
	for i := range output {
		assert.Equal(t, input[i], output[i])
	}
}

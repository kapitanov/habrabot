package carboncopy

import (
	"github.com/kapitanov/habrabot/internal/data"
	"github.com/stretchr/testify/assert"
	"testing"
)

type ExtractFileNameTestCase struct {
	Title    string
	URL      string
	Expected string
}

func TestExtractFileName(t *testing.T) {
	testCases := []ExtractFileNameTestCase{
		{
			Title:    "Программный рендер в стиле игры Doom",
			URL:      "https://habr.com/ru/post/704622/",
			Expected: "704622 - Программный рендер в стиле игры Doom.html",
		},
		{
			Title:    "  Программный рендер в стиле игры Doom  ",
			URL:      "https://habr.com/post/704622",
			Expected: "704622 - Программный рендер в стиле игры Doom.html",
		},
		{
			Title:    "  Программный рендер - в стиле, игры. Doom. /?+*/\\  ",
			URL:      "https://habr.com/post/704622",
			Expected: "704622 - Программный рендер - в стиле, игры. Doom.html",
		},
		{
			Title:    "Программный рендер в стиле игры Doom",
			URL:      "https://habr.com/post/no-article-id",
			Expected: "Программный рендер в стиле игры Doom.html",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.Expected, func(t *testing.T) {
			actual := extractFileName(data.Article{
				Title:   tc.Title,
				LinkURL: tc.URL,
			})

			assert.Equal(t, tc.Expected, actual)
		})
	}
}

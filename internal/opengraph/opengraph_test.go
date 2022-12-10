//nolint:goconst // it's OK for tests
package opengraph

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestParseTags_NoHeadTag(t *testing.T) {
	input := `
<html>
<body>
	<h1>Page title text</h1>
	<p>
		Page content
	</p>
</body>
</html
`

	output := parseTagsTestHelper(t, input)

	assert.Nil(t, output.Title, "Title")
	assert.Nil(t, output.ImageURL, "ImageURL")
}

func TestParseTags_NoTags(t *testing.T) {
	input := `
<html>
<head>
	<title>Page title</title>
	<meta property="og:title:alt" content="Alt. page title" />
</head>
<body>
	<h1>Page title text</h1>
	<p>
		Page content
	</p>
</body>
</html
`

	output := parseTagsTestHelper(t, input)

	assert.Nil(t, output.Title, "Title")
	assert.Nil(t, output.ImageURL, "ImageURL")
}

func TestParseTags_TitleTagOnly(t *testing.T) {
	input := `
<html>
<head>
	<title>Page title</title>
	<meta property="og:title" content="Alt. page title" />
	<meta property="og:title:alt" content="Alt. page title 2" />
</head>
<body>
	<h1>Page title text</h1>
	<p>
		Page content
	</p>
</body>
</html
`

	output := parseTagsTestHelper(t, input)

	if assert.NotNil(t, output.Title, "Title") {
		assert.Equal(t, "Alt. page title", *output.Title, "Title")
	}

	assert.Nil(t, output.ImageURL, "ImageURL")
}

func TestParseTags_ImageTagOnly(t *testing.T) {
	input := `
<html>
<head>
	<title>Page title</title>
	<meta property="og:title:alt" content="Alt. page title" />
	<meta property="og:image" content="https://example.com/image.jpg" />
	<meta property="og:image:url" content="https://example.com/insecure/image.jpg" />
	<meta property="og:image:secure_url" content="https://example.com/secure/image.jpg" />
</head>
<body>
	<h1>Page title text</h1>
	<p>
		Page content
	</p>
</body>
</html
`

	output := parseTagsTestHelper(t, input)

	assert.Nil(t, output.Title, "Title")

	if assert.NotNil(t, output.ImageURL, "ImageURL") {
		assert.Equal(t, "https://example.com/image.jpg", *output.ImageURL, "ImageURL")
	}
}

func TestParseTags_BothTitleAndImageTags(t *testing.T) {
	input := `
<html>
<head>
	<title>Page title</title>
	<meta property="og:title" content="Alt. page title" />
	<meta property="og:title:alt" content="Alt. page title 2" />
	<meta property="og:image" content="https://example.com/image.jpg" />
	<meta property="og:image:url" content="https://example.com/insecure/image.jpg" />
	<meta property="og:image:secure_url" content="https://example.com/secure/image.jpg" />
</head>
<body>
	<h1>Page title text</h1>
	<p>
		Page content
	</p>
</body>
</html
`

	output := parseTagsTestHelper(t, input)

	if assert.NotNil(t, output.Title, "Title") {
		assert.Equal(t, "Alt. page title", *output.Title, "Title")
	}

	if assert.NotNil(t, output.ImageURL, "ImageURL") {
		assert.Equal(t, "https://example.com/image.jpg", *output.ImageURL, "ImageURL")
	}
}

func parseTagsTestHelper(t *testing.T, input string) tags {
	root, err := html.Parse(strings.NewReader(input))
	require.NoError(t, err)

	return parseTags(root)
}

func TestLoadTags_OK(t *testing.T) {
	body := `
<html>
<head>
	<title>Page title</title>
	<meta property="og:title" content="Alt. page title" />
	<meta property="og:title:alt" content="Alt. page title 2" />
	<meta property="og:image" content="https://example.com/image.jpg" />
	<meta property="og:image:url" content="https://example.com/insecure/image.jpg" />
	<meta property="og:image:secure_url" content="https://example.com/secure/image.jpg" />
</head>
<body>
	<h1>Page title text</h1>
	<p>
		Page content
	</p>
</body>
</html
`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/html")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(body))
	}))
	defer server.Close()

	sourceURL := server.URL
	output, err := loadTags(sourceURL)

	if assert.NoError(t, err) {
		if assert.NotNil(t, output.Title, "Title") {
			assert.Equal(t, "Alt. page title", *output.Title, "Title")
		}

		if assert.NotNil(t, output.ImageURL, "ImageURL") {
			assert.Equal(t, "https://example.com/image.jpg", *output.ImageURL, "ImageURL")
		}
	}
}

func TestLoadTags_NonSuccessfulResponse(t *testing.T) {
	statusCodes := []int{
		http.StatusBadRequest,
		http.StatusUnauthorized,
		http.StatusForbidden,
		http.StatusNotFound,
		http.StatusInternalServerError,
		http.StatusNotImplemented,
		http.StatusBadGateway,
		http.StatusGatewayTimeout,
	}

	for _, statusCode := range statusCodes {
		statusCode := statusCode
		t.Run(fmt.Sprintf("HTTP_%v", statusCode), func(t *testing.T) {
			testLoadTagsNonSuccessfulResponseImpl(t, statusCode)
		})
	}
}

func testLoadTagsNonSuccessfulResponseImpl(t *testing.T, status int) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
	}))
	defer server.Close()

	sourceURL := server.URL
	_, err := loadTags(sourceURL)

	assert.Error(t, err)
}

func TestLoadTags_Redirect(t *testing.T) {
	body := `
<html>
<head>
	<title>Page title</title>
	<meta property="og:title" content="Alt. page title" />
	<meta property="og:title:alt" content="Alt. page title 2" />
	<meta property="og:image" content="https://example.com/image.jpg" />
	<meta property="og:image:url" content="https://example.com/insecure/image.jpg" />
	<meta property="og:image:secure_url" content="https://example.com/secure/image.jpg" />
</head>
<body>
	<h1>Page title text</h1>
	<p>
		Page content
	</p>
</body>
</html
`

	actualServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/html")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(body))
	}))
	defer actualServer.Close()

	redirectServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("location", actualServer.URL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}))
	defer redirectServer.Close()

	sourceURL := redirectServer.URL
	output, err := loadTags(sourceURL)

	if assert.NoError(t, err) {
		if assert.NotNil(t, output.Title, "Title") {
			assert.Equal(t, "Alt. page title", *output.Title, "Title")
		}

		if assert.NotNil(t, output.ImageURL, "ImageURL") {
			assert.Equal(t, "https://example.com/image.jpg", *output.ImageURL, "ImageURL")
		}
	}
}

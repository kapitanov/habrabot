//nolint:goconst // it's OK for tests
package telegram

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/kapitanov/habrabot/internal/data"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/utf8string"
)

func TestFormatMessageText(t *testing.T) {
	const href = "https://google.com"

	// <a href="https://google.com"><strong>TITLE</strong></a>NNTEXT
	// 1234567890123456789012345678901234567890123456789012345678901234567890
	//           111111111222222222233333333334444444444555555555566666666
	// 61

	testCases := []struct {
		Name      string
		Title     string
		Text      string
		MaxLength int
		Expected  string
	}{
		{
			Name:      "NoTrim",
			Title:     "TITLE OF THE MESSAGE",
			Text:      "TEXT OF THE MESSAGE",
			MaxLength: 1000,
			Expected:  "<a href=\"https://google.com\"><strong>TITLE OF THE MESSAGE</strong></a>\n\nTEXT OF THE MESSAGE",
		},
		{
			Name:      "TrimText",
			Title:     "TITLE OF THE MESSAGE",
			Text:      "TEXT OF THE MESSAGE",
			MaxLength: 85,
			Expected:  "<a href=\"https://google.com\"><strong>TITLE OF THE MESSAGE</strong></a>\n\nTEXT OF THE \u2026",
		},
		{
			Name:      "TitleOnly",
			Title:     "TITLE OF THE MESSAGE",
			Text:      "TEXT OF THE MESSAGE",
			MaxLength: 72,
			Expected:  "<a href=\"https://google.com\"><strong>TITLE OF THE MESSAGE</strong></a>",
		},
		{
			Name:      "TrimTitle",
			Title:     "TITLE OF THE MESSAGE",
			Text:      "TEXT OF THE MESSAGE",
			MaxLength: 60,
			Expected:  "<a href=\"https://google.com\"><strong>TITLE OF \u2026</strong></a>",
		},
		{
			Name:      "HrefOnly",
			Title:     "TITLE OF THE MESSAGE",
			Text:      "TEXT OF THE MESSAGE",
			MaxLength: 30,
			Expected:  "https://google.com",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.Name, func(t *testing.T) {
			actual := formatMessageText(tc.Title, tc.Text, href, tc.MaxLength)
			t.Logf("Output: %q", actual)
			t.Logf("%d runes", unicodeLength(actual))
			assert.Equal(t, tc.Expected, actual)
		})
	}
}

func TestTrimLongText(t *testing.T) {
	const maxLength = 20

	strLengths := []int{
		maxLength / 2,
		maxLength - 1,
		maxLength,
		maxLength + 1,
		maxLength * 2,
	}

	for _, strLength := range strLengths {
		strLength := strLength

		t.Run(fmt.Sprintf("LEN=%v", strLength), func(t *testing.T) {
			str := ""
			for unicodeLen(str) < strLength {
				str += "lorem ⌘ фыщъ "
			}

			input := unicodeSlice(str, 0, strLength)
			output := trimLongText(input, maxLength)

			t.Logf("input:  %d chars, %q", unicodeLen(input), input)
			t.Logf("output: %d chars, %q", unicodeLen(output), output)

			assert.True(t, unicodeLen(output) <= maxLength)

			if unicodeLen(input) <= maxLength {
				assert.Equal(t, input, output)
			} else {
				assert.True(t, strings.HasSuffix(output, ellipsis))
			}
		})
	}
}

func TestCreateTextMessage_NoTrim(t *testing.T) {
	article := data.Article{
		Title:       "TITLE",
		Description: "",
		LinkURL:     "https://google.com",
	}

	for unicodeLength(article.Description) < maxTextLength {
		article.Description += "lorem "
	}

	article.Description = unicodeSlice(article.Description, 0, maxTextLength-10)
	chatID := int64(1024)
	chattable := createTextMessage(article, chatID)

	if assert.NotNil(t, chattable) {
		if assert.IsType(t, tgbotapi.MessageConfig{}, chattable) {
			msg := chattable.(tgbotapi.MessageConfig)

			assert.Truef(
				t,
				unicodeLength(msg.Text) <= maxTextLength,
				"expected len(msg.Text) = %v <- %v",
				unicodeLength(msg.Text),
				maxTextLength,
			)

			assert.Equal(t, tgbotapi.ModeHTML, msg.ParseMode)
			assert.True(t, msg.DisableWebPagePreview)
			assert.Equal(t, chatID, msg.ChatID)
			assert.Equal(t, "", msg.ChannelUsername)
			assert.Equal(t, 0, msg.ReplyToMessageID)
			assert.Nil(t, msg.ReplyMarkup)
			assert.False(t, msg.DisableNotification)
		}
	}
}

func TestCreateTextMessage_Trim(t *testing.T) {
	strLengths := []int{
		maxTextLength,
		maxTextLength + 1,
		maxTextLength * 2,
	}

	for _, strLength := range strLengths {
		strLength := strLength

		t.Run(fmt.Sprintf("LEN=%v", strLength), func(t *testing.T) {
			article := data.Article{
				Title:       "TITLE",
				Description: "",
				LinkURL:     "https://google.com",
			}

			for unicodeLength(article.Description) < strLength {
				article.Description += "lorem "
			}

			article.Description = unicodeSlice(article.Description, 0, strLength)
			chatID := int64(1024)
			chattable := createTextMessage(article, chatID)

			if assert.NotNil(t, chattable) {
				if assert.IsType(t, tgbotapi.MessageConfig{}, chattable) {
					msg := chattable.(tgbotapi.MessageConfig)

					assert.Truef(
						t,
						unicodeLen(msg.Text) <= maxTextLength,
						"expected len(msg.Text) = %v <= %v",
						unicodeLength(msg.Text),
						maxTextLength,
					)

					assert.Equal(t, tgbotapi.ModeHTML, msg.ParseMode)
					assert.True(t, msg.DisableWebPagePreview)
					assert.Equal(t, chatID, msg.ChatID)
					assert.Equal(t, "", msg.ChannelUsername)
					assert.Equal(t, 0, msg.ReplyToMessageID)
					assert.Nil(t, msg.ReplyMarkup)
					assert.False(t, msg.DisableNotification)
				}
			}
		})
	}
}

func TestCreateTextAndImageMessage_NoTrim(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("image"))
	}))
	defer server.Close()

	article := data.Article{
		Title:       "TITLE",
		Description: "",
		LinkURL:     "https://google.com",
		ImageURL:    &server.URL,
	}
	for len(article.Description) < maxMediaCaptionLength {
		article.Description += "lorem "
	}

	article.Description = unicodeSlice(article.Description, 0, maxMediaCaptionLength-10)
	chatID := int64(1024)
	chattable, err := createTextAndImageMessage(context.Background(), article, chatID, http.DefaultClient)
	assert.NoError(t, err)

	if assert.NotNil(t, chattable) {
		if assert.IsType(t, tgbotapi.PhotoConfig{}, chattable) {
			msg := chattable.(tgbotapi.PhotoConfig)

			assert.Truef(
				t,
				unicodeLen(msg.Caption) <= maxMediaCaptionLength,
				"expected len(msg.Text) = %v <= %v",
				unicodeLength(msg.Caption),
				maxMediaCaptionLength,
			)

			assert.Equal(t, tgbotapi.ModeHTML, msg.ParseMode)
			assert.Equal(t, chatID, msg.ChatID)
			assert.Equal(t, "", msg.ChannelUsername)
			assert.Equal(t, 0, msg.ReplyToMessageID)
			assert.Nil(t, msg.ReplyMarkup)
			assert.False(t, msg.DisableNotification)
		}
	}
}

func TestCreateTextAndImageMessage_Trim(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("image"))
	}))
	defer server.Close()

	strLengths := []int{
		maxMediaCaptionLength,
		maxMediaCaptionLength + 1,
		maxMediaCaptionLength * 2,
	}

	for _, strLength := range strLengths {
		strLength := strLength

		t.Run(fmt.Sprintf("LEN=%v", strLength), func(t *testing.T) {
			article := data.Article{
				Title:       "TITLE",
				Description: "",
				LinkURL:     "https://google.com",
				ImageURL:    &server.URL,
			}
			for unicodeLength(article.Description) < strLength {
				article.Description += "lorem "
			}

			article.Description = unicodeSlice(article.Description, 0, strLength)
			chatID := int64(1024)
			chattable, err := createTextAndImageMessage(context.Background(), article, chatID, http.DefaultClient)
			assert.NoError(t, err)

			if assert.NotNil(t, chattable) {
				if assert.IsType(t, tgbotapi.PhotoConfig{}, chattable) {
					msg := chattable.(tgbotapi.PhotoConfig)

					assert.Truef(
						t,
						unicodeLen(msg.Caption) <= maxMediaCaptionLength,
						"expected len(msg.Text) = %v <= %v",
						unicodeLength(msg.Caption),
						maxMediaCaptionLength,
					)

					assert.Equal(t, tgbotapi.ModeHTML, msg.ParseMode)
					assert.Equal(t, chatID, msg.ChatID)
					assert.Equal(t, "", msg.ChannelUsername)
					assert.Equal(t, 0, msg.ReplyToMessageID)
					assert.Nil(t, msg.ReplyMarkup)
					assert.False(t, msg.DisableNotification)
				}
			}
		})
	}
}

func TestSanitizeText_ReplaceNBSPs(t *testing.T) {
	input := "foo\u00A0bar"
	expected := "foo bar"

	actual := sanitizeText(input)

	assert.Equal(t, expected, actual)
}

func unicodeSlice(s string, i, j int) string {
	return utf8string.NewString(s).Slice(i, j)
}

func unicodeLen(s string) int {
	return utf8string.NewString(s).RuneCount()
}

//nolint:goconst // it's OK for tests
package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTrimLongText(t *testing.T) {
	const maxLength = 1024

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
			for len(str) < strLength {
				str += "lorem "
			}

			input := str[0:strLength]
			output := trimLongText(input, maxLength)
			assert.Truef(
				t,
				len(output) <= maxLength,
				"expected len(output) = %v <- %v",
				len(output),
				maxLength,
			)
		})
	}
}

func TestCreateTextMessage_NoTrim(t *testing.T) {
	str := ""
	for len(str) < maxTextLength {
		str += "lorem "
	}

	str = str[0:(maxTextLength - 10)]
	chatID := int64(1024)
	chattable := createTextMessage(str, chatID)

	if assert.NotNil(t, chattable) {
		if assert.IsType(t, tgbotapi.MessageConfig{}, chattable) {
			msg := chattable.(tgbotapi.MessageConfig)

			assert.Truef(
				t,
				len(msg.Text) <= maxTextLength,
				"expected len(msg.Text) = %v <- %v",
				len(msg.Text),
				maxTextLength,
			)

			assert.Equal(t, str, msg.Text)
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
			str := ""
			for len(str) < strLength {
				str += "lorem "
			}

			str = str[0:strLength]
			chatID := int64(1024)
			chattable := createTextMessage(str, chatID)

			if assert.NotNil(t, chattable) {
				if assert.IsType(t, tgbotapi.MessageConfig{}, chattable) {
					msg := chattable.(tgbotapi.MessageConfig)

					assert.Truef(
						t,
						len(msg.Text) <= maxTextLength,
						"expected len(msg.Text) = %v <= %v",
						len(msg.Text),
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

	str := ""
	for len(str) < maxMediaCaptionLength {
		str += "lorem "
	}

	str = str[0:(maxMediaCaptionLength - 10)]
	chatID := int64(1024)
	chattable, err := createTextAndImageMessage(str, server.URL, chatID)
	assert.NoError(t, err)

	if assert.NotNil(t, chattable) {
		if assert.IsType(t, tgbotapi.PhotoConfig{}, chattable) {
			msg := chattable.(tgbotapi.PhotoConfig)

			assert.Truef(
				t,
				len(msg.Caption) <= maxMediaCaptionLength,
				"expected len(msg.Text) = %v <= %v",
				len(msg.Caption),
				maxMediaCaptionLength,
			)

			assert.Equal(t, str, msg.Caption)
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
			str := ""
			for len(str) < strLength {
				str += "lorem "
			}

			str = str[0:strLength]
			chatID := int64(1024)
			chattable, err := createTextAndImageMessage(str, server.URL, chatID)
			assert.NoError(t, err)

			if assert.NotNil(t, chattable) {
				if assert.IsType(t, tgbotapi.PhotoConfig{}, chattable) {
					msg := chattable.(tgbotapi.PhotoConfig)

					assert.Truef(
						t,
						len(msg.Caption) <= maxMediaCaptionLength,
						"expected len(msg.Text) = %v <= %v",
						len(msg.Caption),
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

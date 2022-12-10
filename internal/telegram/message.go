package telegram

import (
	"fmt"
	"html"
	"io"
	"net/http"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"golang.org/x/exp/utf8string"

	"github.com/kapitanov/habrabot/internal/data"
)

const (
	maxTextLength         = 4096 - 4
	maxMediaCaptionLength = 1024 - 32

	ellipsis = "\u2026"
)

var ellipsisUTF8 = utf8string.NewString(ellipsis)

func prepareMessage(article data.Article, chatID int64) (tgbotapi.Chattable, error) {
	text := fmt.Sprintf(
		"<a href=\"%s\"><strong>%s</strong></a>\n\n%s",
		html.EscapeString(article.LinkURL),
		article.Title,
		article.Description,
	)

	text = sanitizeText(text)

	if article.ImageURL == nil {
		return createTextMessage(text, chatID), nil
	}

	return createTextAndImageMessage(text, *article.ImageURL, chatID)
}

func createTextMessage(text string, chatID int64) tgbotapi.Chattable {
	text = trimLongText(text, maxTextLength)

	msg := tgbotapi.NewMessageToChannel("", text)
	msg.ChatID = chatID
	msg.ParseMode = tgbotapi.ModeHTML
	msg.DisableWebPagePreview = true

	return msg
}

func createTextAndImageMessage(text, imageURL string, chatID int64) (tgbotapi.Chattable, error) {
	bytes, err := downloadImage(imageURL)
	if err != nil {
		return nil, err
	}

	text = trimLongText(text, maxMediaCaptionLength)

	photo := tgbotapi.NewPhotoUpload(chatID, tgbotapi.FileBytes{Bytes: bytes})
	photo.Caption = text
	photo.ParseMode = tgbotapi.ModeHTML

	return photo, nil
}

func downloadImage(url string) ([]byte, error) {
	r, err := http.DefaultClient.Get(url)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = r.Body.Close()
	}()

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func trimLongText(text string, max int) string {
	uText := utf8string.NewString(text)
	if uText.RuneCount() > max {
		rem := max - ellipsisUTF8.RuneCount()
		text = uText.Slice(0, rem)
		text += ellipsis
	}

	return text
}

func sanitizeText(text string) string {
	text = strings.ReplaceAll(text, "\u00A0", " ")
	return text
}

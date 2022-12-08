package telegram

import (
	"fmt"
	"html"
	"io"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/kapitanov/habrabot/internal/data"
)

const (
	maxTextLength         = 4096 - 4
	maxMediaCaptionLength = 1024 - 32
)

func prepareMessage(article data.Article, chatID int64) (tgbotapi.Chattable, error) {
	text := fmt.Sprintf(
		"<a href=\"%s\"><strong>%s</strong></a>\n\n%s",
		html.EscapeString(article.LinkURL),
		article.Title,
		article.Description,
	)

	if article.ImageURL == nil {
		text = trimLongText(text, maxMediaCaptionLength)
		return createTextMessage(text, chatID), nil
	}

	text = trimLongText(text, maxTextLength)
	return createTextAndImageMessage(text, *article.ImageURL, chatID)
}

func createTextMessage(text string, chatID int64) tgbotapi.Chattable {
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
	if len(text) > max {
		text = text[0:max-3] + "..."
	}
	return text
}

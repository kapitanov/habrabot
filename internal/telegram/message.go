package telegram

import (
	"context"
	"fmt"
	"html"
	"io"
	"net/http"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/kapitanov/habrabot/internal/data"
	"golang.org/x/exp/utf8string"
)

const (
	maxTextLength         = 4096 - 4
	maxMediaCaptionLength = 1024 - 32

	ellipsis = "\u2026"
)

var ellipsisUTF8 = utf8string.NewString(ellipsis)

func prepareMessage(
	ctx context.Context,
	article data.Article,
	chatID int64,
	httpClient *http.Client,
) (tgbotapi.Chattable, error) {
	if article.ImageURL == nil {
		return createTextMessage(article, chatID), nil
	}

	return createTextAndImageMessage(ctx, article, chatID, httpClient)
}

func createTextMessage(article data.Article, chatID int64) tgbotapi.Chattable {
	text := formatMessageText(article.Title, article.Description, article.LinkURL, maxTextLength)

	msg := tgbotapi.NewMessageToChannel("", text)
	msg.ChatID = chatID
	msg.ParseMode = tgbotapi.ModeHTML
	msg.DisableWebPagePreview = true

	return msg
}

func createTextAndImageMessage(
	ctx context.Context,
	article data.Article,
	chatID int64,
	httpClient *http.Client,
) (tgbotapi.Chattable, error) {
	bytes, err := downloadImage(ctx, *article.ImageURL, httpClient)
	if err != nil {
		return nil, err
	}

	text := formatMessageText(article.Title, article.Description, article.LinkURL, maxMediaCaptionLength)

	photo := tgbotapi.NewPhotoUpload(chatID, tgbotapi.FileBytes{Bytes: bytes})
	photo.Caption = text
	photo.ParseMode = tgbotapi.ModeHTML

	return photo, nil
}

func downloadImage(
	ctx context.Context,
	url string,
	httpClient *http.Client,
) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func formatMessageText(title, text, href string, maxLength int) string {
	const titleTextSeparator = "\n\n"

	title = sanitizeText(title)

	formattedTitle := fmt.Sprintf("<a href=\"%s\"><strong>%s</strong></a>", html.EscapeString(href), title)

	if unicodeLength(formattedTitle) > maxLength {
		maxTitleLength := maxLength - unicodeLength(
			fmt.Sprintf("<a href=\"%s\"><strong></strong></a>", html.EscapeString(href)),
		)

		if maxTitleLength <= 0 {
			return href
		}

		title = trimLongText(title, maxTitleLength)
		formattedTitle = fmt.Sprintf("<a href=\"%s\"><strong>%s</strong></a>", html.EscapeString(href), title)
		return formattedTitle
	}

	remMaxTextLength := maxLength - unicodeLength(formattedTitle) - unicodeLength(titleTextSeparator)

	if remMaxTextLength <= 0 {
		return formattedTitle
	}

	text = sanitizeText(text)
	text = trimLongText(text, remMaxTextLength)

	formattedText := fmt.Sprintf("%s%s%s", formattedTitle, titleTextSeparator, text)
	return formattedText
}

func unicodeLength(text string) int {
	return utf8string.NewString(text).RuneCount()
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
	text = strings.ReplaceAll(text, "\u00a0", " ")
	return text
}

package delivery

import (
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/kapitanov/habrabot/source"
)

type telegramChannel struct {
	bot  *tgbotapi.BotAPI
	chat tgbotapi.Chat
}

// NewTelegramChannel creates new delivery channel that publishes messages into Telegram channel
func NewTelegramChannel(token, channelNameOrID string) (Channel, error) {
	bot, err := connectToTelegram(token)
	if err != nil {
		return nil, err
	}

	me, err := bot.GetMe()
	if err != nil {
		return nil, err
	}
	log.Printf("Connected to telegram as @%s", me.UserName)

	chat, err := bot.GetChat(tgbotapi.ChatConfig{SuperGroupUsername: channelNameOrID})
	if err != nil {
		return nil, err
	}
	log.Printf("Will post messages to telegram channel @%s (%d)", chat.UserName, chat.ID)

	return &telegramChannel{bot, chat}, nil
}

func connectToTelegram(token string) (*tgbotapi.BotAPI, error) {
	proxyURLStr := os.Getenv("HTTP_PROXY")
	if proxyURLStr == "" {
		return tgbotapi.NewBotAPI(token)
	}

	proxyURL, err := url.Parse(proxyURLStr)
	if err != nil {
		return nil, err
	}

	log.Printf("Will use proxy server \"%s://%s\" to connect to telegram", proxyURL.Scheme, proxyURL.Host)

	httpClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}}
	bot, err := tgbotapi.NewBotAPIWithClient(token, httpClient)
	if err != nil {
		return nil, err
	}

	return bot, nil
}

func (c *telegramChannel) Publish(article *source.Article) error {
	text := fmt.Sprintf("<a href=\"%s\"><strong>%s</strong></a>\n\n%s", html.EscapeString(article.LinkURL), article.Title, article.Description)

	if article.ImageURL != "" {
		text = trimLongText(text, 4096);
		return c.publishTextAndImage(article, text, article.ImageURL)
	} else {
		text = trimLongText(text, 1024);
		return c.publishText(article, text)
	}
}

func (c *telegramChannel) publishText(article *source.Article, text string) error {
	msg := tgbotapi.NewMessageToChannel("", text)
	msg.ChatID = c.chat.ID
	msg.ParseMode = tgbotapi.ModeHTML
	msg.DisableWebPagePreview = true

	result, err := c.bot.Send(msg)
	if err != nil {
		return err
	}

	log.Printf("Posted new text message #%d to telegram channel @%s", result.MessageID, c.chat.UserName)
	return nil
}

func (c *telegramChannel) publishTextAndImage(article *source.Article, text, imageURL string) error {
	bytes, err := downloadImage(imageURL)
	if err != nil {
		return err
	}

	photo := tgbotapi.NewPhotoUpload(c.chat.ID, tgbotapi.FileBytes{Bytes: bytes})
	photo.Caption = text
	photo.ParseMode = tgbotapi.ModeHTML

	result, err := c.bot.Send(photo)
	if err != nil {
		return err
	}

	log.Printf("Posted new image message #%d to telegram channel @%s", result.MessageID, c.chat.UserName)
	return nil
}

func downloadImage(url string) ([]byte, error) {
	r, err := http.DefaultClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func trimLongText(text string, max int) string {
	if len(text) > max {
		text = text[0:max - 3] + "..."
	}
	return text
}

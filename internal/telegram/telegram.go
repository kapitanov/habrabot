package telegram

import (
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/kapitanov/habrabot/internal/data"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// New creates new consumed that publishes messages into Telegram channel.
func New(token, channelNameOrID string) data.Consumer {
	t := &transmitter{
		token:           token,
		channelNameOrID: channelNameOrID,
		bot:             nil,
		chat:            nil,
	}

	return t
}

type transmitter struct {
	token           string
	channelNameOrID string
	bot             *tgbotapi.BotAPI
	chat            *tgbotapi.Chat
}

// On method is invoked when an article is received from the feed.
func (t *transmitter) On(article data.Article) error {
	if t.bot == nil {
		bot, err := connectToTelegram(t.token)
		if err != nil {
			return err
		}

		t.bot = bot
	}

	if t.chat == nil {
		chat, err := selectChat(t.bot, t.channelNameOrID)
		if err != nil {
			return err
		}

		t.chat = chat
	}

	msg, err := prepareMessage(article, t.chat.ID)
	if err != nil {
		return err
	}

	result, err := t.bot.Send(msg)
	if err != nil {
		return err
	}

	log.Printf("posted new message #%d to telegram channel @%s", result.MessageID, t.chat.UserName)
	return nil
}

func connectToTelegram(token string) (*tgbotapi.BotAPI, error) {
	httpTransport := &http.Transport{}

	proxyURLStr := os.Getenv("HTTP_PROXY")
	if proxyURLStr == "" {
		return tgbotapi.NewBotAPI(token)
	} else {
		proxyURL, err := url.Parse(proxyURLStr)
		if err != nil {
			return nil, err
		}

		httpTransport.Proxy = http.ProxyURL(proxyURL)

		log.Printf("will use proxy server \"%s://%s\" to connect to telegram", proxyURL.Scheme, proxyURL.Host)
	}

	httpClient := &http.Client{
		Transport: httpTransport,
	}
	bot, err := tgbotapi.NewBotAPIWithClient(token, httpClient)
	if err != nil {
		return nil, err
	}

	me, err := bot.GetMe()
	if err != nil {
		return nil, err
	}
	log.Printf("connected to telegram as @%s", me.UserName)
	return bot, nil
}

func selectChat(bot *tgbotapi.BotAPI, channelNameOrID string) (*tgbotapi.Chat, error) {
	chat, err := bot.GetChat(tgbotapi.ChatConfig{SuperGroupUsername: channelNameOrID})
	if err != nil {
		return nil, err
	}

	log.Printf("will post messages to telegram channel @%s (%d)", chat.UserName, chat.ID)
	return &chat, nil
}

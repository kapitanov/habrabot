package telegram

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/kapitanov/habrabot/internal/data"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rs/zerolog/log"
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
			log.Error().Err(err).Msg("unable to connect to telegram")
			return err
		}

		t.bot = bot
	}

	if t.chat == nil {
		chat, err := selectChat(t.bot, t.channelNameOrID)
		if err != nil {
			log.Error().Err(err).Str("chat", t.channelNameOrID).Msg("unable to select chat")
			return err
		}

		t.chat = chat
	}

	msg, err := prepareMessage(article, t.chat.ID)
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare telegram message")
		return err
	}

	result, err := t.bot.Send(msg)
	if err != nil {
		log.Error().
			Err(err).
			Str("title", article.Title).
			Str("id", article.ID).
			Msg("unable to send to telegram")
		return err
	}

	log.Info().
		Int("msg", result.MessageID).
		Str("channel", fmt.Sprintf("@%v", t.chat.UserName)).
		Str("title", article.Title).
		Str("id", article.ID).
		Msg("posted a telegram message")
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

		log.Info().
			Str("proxy", fmt.Sprintf("%s://%s", proxyURL.Scheme, proxyURL.Host)).
			Msg("will use proxy server for telegram")
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

	log.Info().
		Str("me", fmt.Sprintf("@%v", me.UserName)).
		Msg("connected to telegram")
	return bot, nil
}

func selectChat(bot *tgbotapi.BotAPI, channelNameOrID string) (*tgbotapi.Chat, error) {
	chat, err := bot.GetChat(tgbotapi.ChatConfig{SuperGroupUsername: channelNameOrID})
	if err != nil {
		return nil, err
	}

	log.Info().
		Str("channel", fmt.Sprintf("@%v", chat.UserName)).
		Int64("id", chat.ID).
		Msg("will post messages to telegram channel")
	return &chat, nil
}

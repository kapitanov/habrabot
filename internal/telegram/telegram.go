package telegram

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/rs/zerolog/log"

	"github.com/kapitanov/habrabot/internal/data"
	"github.com/kapitanov/habrabot/internal/httpclient"
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
	httpClient      *retryablehttp.Client
	bot             *tgbotapi.BotAPI
	chat            *tgbotapi.Chat
}

// On method is invoked when an article is received from the feed.
func (t *transmitter) On(article data.Article) error {
	err := t.connectToTelegram()
	if err != nil {
		log.Error().Err(err).Msg("unable to connect to telegram")
		return err
	}

	err = t.selectChat()
	if err != nil {
		log.Error().Err(err).Str("chat", t.channelNameOrID).Msg("unable to select chat")
		return err
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

func (t *transmitter) createHTTPClient() error {
	if t.httpClient != nil {
		return nil
	}

	httpClient, err := createHTTPClient()
	if err != nil {
		return err
	}

	t.httpClient = httpClient
	return nil
}

func createHTTPClient() (*retryablehttp.Client, error) {
	httpClient, err := httpclient.New(httpclient.TelegramPolicy)
	if err != nil {
		return nil, err
	}

	return httpClient, nil
}

func (t *transmitter) connectToTelegram() error {
	if t.bot != nil {
		return nil
	}

	err := t.createHTTPClient()
	if err != nil {
		return err
	}

	bot, err := connectToTelegram(t.httpClient, t.token)
	if err != nil {
		return err
	}

	t.bot = bot
	return nil
}

func connectToTelegram(httpClient *retryablehttp.Client, token string) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPIWithClient(token, httpClient.StandardClient())
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

func (t *transmitter) selectChat() error {
	if t.chat == nil {
		chat, err := selectChat(t.bot, t.channelNameOrID)
		if err != nil {
			log.Error().Err(err).Str("chat", t.channelNameOrID).Msg("unable to select chat")
			return err
		}

		t.chat = chat
	}
	return nil
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

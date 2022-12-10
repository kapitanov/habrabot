package httpclient

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type policy interface {
	ConfigureHTTP(client *retryablehttp.Client)
	CreateLogger() zerolog.Logger
}

var (
	TelegramPolicy  policy = telegramPolicy{}
	RSSPolicy       policy = rssPolicy{}
	OpengraphPolicy policy = opengraphPolicy{}
)

// New creates new HTTP client with proper resilience policy.
func New(p policy) (*retryablehttp.Client, error) {
	logger := p.CreateLogger()

	innerHTTPClient, err := createInnerHTTPClient(logger)
	if err != nil {
		return nil, err
	}

	httpClient := retryablehttp.NewClient()
	p.ConfigureHTTP(httpClient)

	httpClient.Logger = loggerAdapter{logger}
	httpClient.HTTPClient = innerHTTPClient

	return httpClient, nil
}

func createInnerHTTPClient(logger zerolog.Logger) (*http.Client, error) {
	httpTransport := &http.Transport{}

	proxyURLStr := os.Getenv("HTTP_PROXY")
	if proxyURLStr != "" {
		proxyURL, err := url.Parse(proxyURLStr)
		if err != nil {
			return nil, err
		}

		httpTransport.Proxy = http.ProxyURL(proxyURL)

		logger.Info().
			Str("proxy", fmt.Sprintf("%s://%s", proxyURL.Scheme, proxyURL.Host)).
			Msg("will use proxy server")
	}

	return &http.Client{
		Transport: httpTransport,
	}, nil
}

type loggerAdapter struct {
	log zerolog.Logger
}

func (a loggerAdapter) Error(format string, args ...interface{}) {
	a.log.Error().Msgf(format, args...)
}

func (a loggerAdapter) Warn(format string, args ...interface{}) {
	a.log.Warn().Msgf(format, args...)
}

func (a loggerAdapter) Info(format string, args ...interface{}) {
	a.log.Info().Msgf(format, args...)
}

func (a loggerAdapter) Debug(format string, args ...interface{}) {
	a.log.Debug().Msgf(format, args...)
}

type telegramPolicy struct{}

func (_ telegramPolicy) ConfigureHTTP(client *retryablehttp.Client) {
	client.Backoff = retryablehttp.LinearJitterBackoff
	client.RetryMax = 10
	client.RetryWaitMin = time.Second
	client.RetryWaitMin = 30 * time.Second
}

func (_ telegramPolicy) CreateLogger() zerolog.Logger {
	return log.Logger.With().Str("component", "telegram").Logger()
}

type rssPolicy struct{}

func (_ rssPolicy) ConfigureHTTP(client *retryablehttp.Client) {
	client.Backoff = retryablehttp.LinearJitterBackoff
	client.RetryMax = 10
	client.RetryWaitMin = time.Second
	client.RetryWaitMin = 30 * time.Second
	client.CheckRetry = func(ctx context.Context, resp *http.Response, err error) (bool, error) {
		if ctx.Err() != nil {
			return false, ctx.Err()
		}

		if resp.StatusCode == http.StatusNotFound {
			return false, nil
		}

		return true, nil
	}
}

func (_ rssPolicy) CreateLogger() zerolog.Logger {
	return log.Logger.With().Str("component", "rss").Logger()
}

type opengraphPolicy struct{}

func (_ opengraphPolicy) ConfigureHTTP(client *retryablehttp.Client) {
	client.Backoff = retryablehttp.LinearJitterBackoff
	client.RetryMax = 10
	client.RetryWaitMin = time.Second
	client.RetryWaitMin = 30 * time.Second
}

func (_ opengraphPolicy) CreateLogger() zerolog.Logger {
	return log.Logger.With().Str("component", "opengraph").Logger()
}

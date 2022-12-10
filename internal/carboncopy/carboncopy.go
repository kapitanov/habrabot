package carboncopy

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/rs/zerolog/log"

	"github.com/kapitanov/habrabot/internal/data"
	"github.com/kapitanov/habrabot/internal/httpclient"
)

func Use(dirPath string) data.Consumer {
	log.Info().Str("dir", dirPath).Msg("will store local copy of feed items")
	return &consumer{
		dirPath: dirPath,
	}
}

type consumer struct {
	dirPath    string
	httpClient *retryablehttp.Client
}

// On method is invoked when an article is received from the feed.
func (c *consumer) On(ctx context.Context, article data.Article) error {
	f, err := c.open(article)
	if err != nil {
		log.Error().Err(err).Str("id", article.ID).Msg("unable to open cc file")
		return err
	}
	defer func() {
		_ = f.Close()
	}()

	resp, err := c.download(ctx, article)
	if err != nil {
		log.Error().Err(err).Str("id", article.ID).Msg("unable to download cc file")
		return err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		log.Error().Err(err).Str("id", article.ID).Msg("unable to download cc file")
		return err
	}

	log.Info().Str("id", article.ID).Str("file", f.Name()).Msg("feed item has been stored locally")
	return nil
}

func (c *consumer) open(article data.Article) (*os.File, error) {
	fullpath := filepath.Join(c.dirPath, extractFileName(article))

	dirName := filepath.Dir(fullpath)
	err := os.MkdirAll(dirName, 0755)
	if err != nil {
		log.Error().Err(err).Str("path", fullpath).Msg("mkdir error")
		return nil, err
	}

	f, err := os.OpenFile(fullpath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
	if err != nil {
		log.Error().Err(err).Str("path", fullpath).Msg("fopen error")
		return nil, err
	}

	return f, nil
}

func (c *consumer) download(ctx context.Context, article data.Article) (*http.Response, error) {
	if c.httpClient == nil {
		httpClient, err := httpclient.New(httpclient.CCPolicy)
		if err != nil {
			return nil, err
		}

		c.httpClient = httpClient
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, article.LinkURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.StandardClient().Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= http.StatusMultipleChoices {
		log.Warn().Err(err).
			Str("url", article.LinkURL).
			Int("status", resp.StatusCode).
			Msg("unable to download web page")
		return nil, fmt.Errorf("unable to download \"%s\": %v", article.LinkURL, resp.Status)
	}

	return resp, nil
}

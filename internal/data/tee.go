package data

import "context"

// Tee splits feed stream between consumers, implementing a fanout pattern.
func Tee(consumers ...Consumer) Consumer {
	return tee{
		consumers: consumers,
	}
}

type tee struct {
	consumers []Consumer
}

// On method is invoked when an article is received from the feed.
func (t tee) On(ctx context.Context, article Article) error {
	for _, consumer := range t.consumers {
		err := consumer.On(ctx, article)
		if err != nil {
			return err
		}
	}

	return nil
}

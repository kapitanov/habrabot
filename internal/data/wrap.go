package data

// Middleware allows to hook into stream processing pipeline.
type Middleware interface {
	// Do method executes an action over a stream item.
	Do(article Article, next func(article Article) error) error
}

// MiddlewareFunc is a function that implements Middleware interface.
type MiddlewareFunc func(article Article, next func(article Article) error) error

// Do method executes an action over a stream item.
func (mw MiddlewareFunc) Do(article Article, next func(Article) error) error {
	return mw(article, next)
}

// Wrap applies wraps stream into a middleware.
func Wrap(feed Feed, middleware Middleware) Feed {
	return &wrapper{
		feed:       feed,
		middleware: middleware,
	}
}

type wrapper struct {
	feed       Feed
	middleware Middleware
}

// Read method reads feed items and streams them into the consumer.
func (w *wrapper) Read(consumer Consumer) error {
	return w.feed.Read(ConsumerFunc(func(article Article) error {
		return w.middleware.Do(article, consumer.On)
	}))
}

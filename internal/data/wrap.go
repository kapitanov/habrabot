package data

import "context"

// NextFunc is a "next" function for middlewares.
type NextFunc func(article Article) error

// Middleware allows to hook into stream processing pipeline.
type Middleware interface {
	// Do method executes an action over a stream item.
	Do(ctx context.Context, article Article, next NextFunc) error
}

// MiddlewareFunc is a function that implements Middleware interface.
type MiddlewareFunc func(ctx context.Context, article Article, next NextFunc) error

// Do method executes an action over a stream item.
func (mw MiddlewareFunc) Do(ctx context.Context, article Article, next NextFunc) error {
	return mw(ctx, article, next)
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
func (w *wrapper) Read(ctx context.Context, consumer Consumer) error {
	return w.feed.Read(ctx, ConsumerFunc(func(ctx context.Context, article Article) error {
		var next NextFunc = func(a Article) error {
			return consumer.On(ctx, a)
		}

		return w.middleware.Do(ctx, article, next)
	}))
}

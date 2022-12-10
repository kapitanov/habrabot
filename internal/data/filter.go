package data

import "context"

// Predicate is a predicate interface for feed stream's filtering.
type Predicate interface {
	// Filter returns true if an article passes through the filter, and false otherwise.
	Filter(ctx context.Context, article Article) (bool, error)
}

// PredicateFunc is a function that implements Predicate interface.
type PredicateFunc func(ctx context.Context, article Article) (bool, error)

// Filter returns true if an article passes through the filter, and false otherwise.
func (f PredicateFunc) Filter(ctx context.Context, article Article) (bool, error) {
	return f(ctx, article)
}

// Filter applies a predicate to the feed's stream.
func Filter(feed Feed, predicate Predicate) Feed {
	return Wrap(feed, MiddlewareFunc(func(ctx context.Context, article Article, next NextFunc) error {
		ok, err := predicate.Filter(ctx, article)

		if err != nil {
			return err
		}

		if !ok {
			return nil
		}

		return next(article)
	}))
}

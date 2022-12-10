package data

import "context"

// Transformation defines a feed item's transformation.
type Transformation interface {
	// Apply updates an article according to transformation logic.
	Apply(ctx context.Context, article *Article) error
}

// TransformationFunc is a function that implements Transformation interface.
type TransformationFunc func(ctx context.Context, article *Article) error

// Apply updates an article according to transformation logic.
func (f TransformationFunc) Apply(ctx context.Context, article *Article) error {
	return f(ctx, article)
}

// Transform applies a function to the feed's stream.
func Transform(feed Feed, transformation Transformation) Feed {
	return Wrap(feed, MiddlewareFunc(func(ctx context.Context, article Article, next NextFunc) error {
		err := transformation.Apply(ctx, &article)
		if err != nil {
			return err
		}

		return next(article)
	}))
}

package data

// Transformation defines a feed item's transformation.
type Transformation interface {
	// Apply updates an article according to transformation logic.
	Apply(article *Article) error
}

// TransformationFunc is a function that implements Transformation interface.
type TransformationFunc func(article *Article) error

// Apply updates an article according to transformation logic.
func (f TransformationFunc) Apply(article *Article) error {
	return f(article)
}

// Transform applies a function to the feed's stream.
func Transform(feed Feed, transformation Transformation) Feed {
	return Wrap(feed, MiddlewareFunc(func(article Article, next func(article Article) error) error {
		err := transformation.Apply(&article)
		if err != nil {
			return err
		}

		return next(article)
	}))
}

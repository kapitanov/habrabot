package data

// Predicate is a predicate interface for feed stream's filtering.
type Predicate interface {
	// Filter returns true if an article passes through the filter, and false otherwise.
	Filter(article Article) (bool, error)
}

// PredicateFunc is a function that implements Predicate interface.
type PredicateFunc func(article Article) (bool, error)

// Filter returns true if an article passes through the filter, and false otherwise.
func (f PredicateFunc) Filter(article Article) (bool, error) {
	return f(article)
}

// Filter applies a predicate to the feed's stream.
func Filter(feed Feed, predicate Predicate) Feed {
	return Wrap(feed, MiddlewareFunc(func(article Article, next func(article Article) error) error {
		ok, err := predicate.Filter(article)

		if err != nil {
			return err
		}

		if !ok {
			return nil
		}

		return next(article)
	}))
}

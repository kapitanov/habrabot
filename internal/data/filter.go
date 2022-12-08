package data

// Filter applies a predicate to the feed's stream.
func Filter(feed Feed, predicate func(article Article) (bool, error)) Feed {
	return &filter{
		feed:      feed,
		predicate: predicate,
	}
}

type filter struct {
	feed      Feed
	predicate func(article Article) (bool, error)
}

// Read method reads feed items and streams them into the consumer.
func (f *filter) Read(consumer Consumer) error {
	return f.feed.Read(ConsumerFunc(func(article Article) error {
		ok, err := f.predicate(article)

		if err != nil {
			return err
		}

		if !ok {
			return nil
		}

		return consumer.On(article)
	}))
}

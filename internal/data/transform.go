package data

// Transform applies a function to the feed's stream.
func Transform(feed Feed, transform func(article *Article) error) Feed {
	return &transformer{
		feed:      feed,
		transform: transform,
	}
}

type transformer struct {
	feed      Feed
	transform func(article *Article) error
}

// Read method reads feed items and streams them into the consumer.
func (t *transformer) Read(consumer Consumer) error {
	return t.feed.Read(ConsumerFunc(func(article Article) error {
		err := t.transform(&article)
		if err != nil {
			return err
		}

		return consumer.On(article)
	}))
}

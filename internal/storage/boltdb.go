package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	bolt "go.etcd.io/bbolt"

	"github.com/kapitanov/habrabot/internal/data"
)

var bucketName = []byte("articles")

func UseBoltDB(feed data.Feed, dbPath string) data.Feed {
	log.Info().Str("path", dbPath).Msg("using boltdb storage")

	storage := &boltDBStorage{
		dbPath: dbPath,
	}

	// TODO turn into middleware

	return data.Filter(feed, storage)
}

type boltDBStorage struct {
	dbPath string
}

// Filter returns true if an article passes through the filter, and false otherwise.
func (s *boltDBStorage) Filter(article data.Article) (bool, error) {
	var isVisible bool
	err := executeTX(s.dbPath, func(tx *bolt.Tx) error {
		bucket, e := ensureBucket(tx)
		if e != nil {
			return e
		}

		isVisible, e = doFilter(bucket, article)
		if e != nil {
			log.Error().Err(e).Msg("unable to apply filter")
			return e
		}

		return nil
	})
	return isVisible, err
}

func executeTX(dbPath string, fn func(tx *bolt.Tx) error) error {
	db, err := openDB(dbPath)
	if err != nil {
		log.Error().Err(err).Str("path", dbPath).Msg("unable to open db file")
		return err
	}

	defer func() {
		_ = db.Close()
	}()

	err = db.Update(func(tx *bolt.Tx) error {
		return fn(tx)
	})
	if err != nil {
		return err
	}

	return nil
}

func openDB(dbPath string) (*bolt.DB, error) {
	dbPath, err := filepath.Abs(dbPath)
	if err != nil {
		return nil, err
	}

	err = os.MkdirAll(filepath.Dir(dbPath), os.ModePerm)
	if err != nil {
		return nil, err
	}

	db, err := bolt.Open(dbPath, os.ModePerm, nil)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func ensureBucket(tx *bolt.Tx) (*bolt.Bucket, error) {
	bucket := tx.Bucket(bucketName)
	if bucket == nil {
		var err error
		bucket, err = tx.CreateBucket(bucketName)
		if err != nil {
			log.Error().Err(err).Str("bucket", string(bucketName)).Msg("unable to create bucket")
			return nil, err
		}
	}

	return bucket, nil
}

func doFilter(bucket *bolt.Bucket, article data.Article) (bool, error) {
	key := []byte(strings.ToLower(article.ID))
	if bucket.Get(key) != nil {
		return false, nil
	}

	value, err := json.Marshal(article)
	if err != nil {
		return false, err
	}

	err = bucket.Put(key, value)
	if err != nil {
		return false, err
	}

	return true, nil
}

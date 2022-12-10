package storage

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"

	bolt "go.etcd.io/bbolt"

	"github.com/kapitanov/habrabot/internal/data"
)

var bucketName = []byte("articles")

func UseBoltDB(feed data.Feed, dbPath string) data.Feed {
	log.Printf("initialized boltdb storage file \"%s\"", dbPath)

	storage := &boltDBStorage{
		dbPath: dbPath,
	}

	return data.Filter(feed, storage)
}

type boltDBStorage struct {
	dbPath string
}

// Filter returns true if an article passes through the filter, and false otherwise.
func (s *boltDBStorage) Filter(article data.Article) (bool, error) {
	var isVisible bool
	err := openDB(s.dbPath, func(tx *bolt.Tx) error {
		bucket, e := ensureBucket(tx)
		if e != nil {
			return e
		}

		isVisible, e = doFilter(bucket, article)
		if e != nil {
			return e
		}

		return nil
	})
	return isVisible, err
}

func openDB(dbPath string, fn func(tx *bolt.Tx) error) error {
	dbPath, err := filepath.Abs(dbPath)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(dbPath), os.ModePerm)
	if err != nil {
		return err
	}

	db, err := bolt.Open(dbPath, os.ModePerm, nil)
	if err != nil {
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

func ensureBucket(tx *bolt.Tx) (*bolt.Bucket, error) {
	bucket := tx.Bucket(bucketName)
	if bucket == nil {
		var err error
		bucket, err = tx.CreateBucket(bucketName)
		if err != nil {
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

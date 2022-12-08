package storage

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/kapitanov/habrabot/internal/source"

	bolt "go.etcd.io/bbolt"
)

type boltdbDriver struct {
	dbPath string
}

var (
	bucketName = []byte("articles")
)

// NewBoltDBDriver creates new storage driver tha uses BoltDB as a storage engine.
func NewBoltDBDriver(dbPath string) (Driver, error) {
	driver := &boltdbDriver{dbPath}
	err := driver.update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		if bucket == nil {
			_, err := tx.CreateBucket(bucketName)
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	log.Printf("Initialized BoltDB storage file \"%s\"", dbPath)

	return driver, nil
}

// Store tries to write an article and returns storage status.
func (d *boltdbDriver) Store(article *source.Article) (Status, error) {
	status := New
	err := d.update(func(tx *bolt.Tx) error {
		key := []byte(strings.ToLower(article.ID))

		bucket := tx.Bucket(bucketName)
		if bucket.Get(key) != nil {
			status = Old
			return nil
		}

		value, err := json.Marshal(article)
		if err != nil {
			return err
		}

		err = bucket.Put(key, value)
		if err != nil {
			return err
		}

		return nil
	})

	return status, err
}

func (d *boltdbDriver) update(fn func(*bolt.Tx) error) error {
	db, err := bolt.Open(d.dbPath, 0600, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Update(fn)
	return err
}

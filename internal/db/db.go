package db

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gofrs/flock"
	"github.com/spf13/viper"
	bolt "go.etcd.io/bbolt"
)

const (
	KV_ENV_VAR     = "KV_DB"
	KV_DEFAULT_DB  = ".kv.db"
	KV_BUCKET_KV   = "kv"
	KV_BUCKET_KEYS = "k"
)

type KV struct {
	db   *bolt.DB
	lock *flock.Flock
}

func hashKey(key string) string {
	hash := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hash[:])
}

func NewKV() *KV {
	dbFile := viper.GetString(KV_ENV_VAR)
	if dbFile == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		dbFile = filepath.Join(homeDir, KV_DEFAULT_DB)
	}

	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	lockFile := dbFile + ".lock"
	lock := flock.New(lockFile)

	return &KV{
		db:   db,
		lock: lock,
	}
}

func (kv *KV) Close() {
	if kv.db != nil {
		kv.db.Close()
	}
}

func (kv *KV) acquireLock() {
	locked, err := kv.lock.TryLock()
	if err != nil {
		log.Fatalf("Unable to acquire lock: %v", err)
	}
	if !locked {
		log.Fatalf("Unable to acquire lock: already locked")
	}
}

func (kv *KV) releaseLock() {
	err := kv.lock.Unlock()
	if err != nil {
		log.Fatalf("Unable to release lock: %v", err)
	}
}

func (kv *KV) Set(key, value string) {
	// multiple execs can happen at the same time, so let's make sure no
	// concurrent writes are performed
	kv.acquireLock()
	defer kv.releaseLock()

	err := kv.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(KV_BUCKET_KV))
		// keeping keys in a different bucket to allow listing
		k, err := tx.CreateBucketIfNotExists([]byte(KV_BUCKET_KEYS))
		if err != nil {
			return err
		}
		// hashing keys so I use consistent values even if keys are large freetext
		hashedKey := hashKey(key)
		err = k.Put([]byte(key), []byte(hashedKey))
		if err != nil {
			return err
		}
		return bucket.Put([]byte(hashedKey), []byte(value))
	})
	if err != nil {
		log.Fatal(err)
	}
}

// Get returns a value if current key exists
func (kv *KV) Get(key string) string {
	var value string
	err := kv.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(KV_BUCKET_KV))
		if bucket == nil {
			return fmt.Errorf("bucket not found")
		}
		hashedKey := hashKey(key)
		val := bucket.Get([]byte(hashedKey))
		if val == nil {
			return fmt.Errorf("key not found")
		}
		value = string(val)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return value
}

// List all keys in the bucket key
func (kv *KV) List() {
	err := kv.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(KV_BUCKET_KEYS))
		if bucket == nil {
			return fmt.Errorf("bucket not found")
		}
		return bucket.ForEach(func(k, v []byte) error {
			fmt.Println(string(k))
			return nil
		})
	})
	if err != nil {
		log.Fatal(err)
	}
}

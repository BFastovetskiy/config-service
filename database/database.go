package database

import (
	"config-service/about"
	"config-service/core"
	"config-service/utils"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/boltdb/bolt"

	_ "net/http/pprof"
)

func InitDatabase(log core.ILogger, dbPath, dbName string) *Database {
	if !utils.ExistFileOrDir(dbPath) {
		os.Mkdir(dbPath, os.ModeDir)
	}
	d := &Database{
		log:  log,
		path: dbPath,
		Name: dbName,
	}
	db, _ := d.openDatabase()
	if db == nil {
		return nil
	}
	d.db = db
	d.setupDB()
	return d
}

func (d Database) openDatabase() (*bolt.DB, error) {
	var err error
	var db *bolt.DB
	if db, err = bolt.Open(filepath.Join(d.path, d.Name), 0600, nil); err != nil {
		d.log.Errorf("Could not open db, %v", err)
		return nil, err
	}
	return db, nil
}

func (d Database) setupDB() {
	err := d.db.Update(func(tx *bolt.Tx) error {
		root, err := tx.CreateBucketIfNotExists(toByte(about.Database_Namespace_Root))
		if err != nil {
			d.log.Errorf("Could not create root bucket: %v", err)
		}
		d.root = root
		return nil
	})
	if err != nil {
		d.log.Errorf("Could not set up buckets, %v", err)
	}
	d.log.Info("Database setup done")
}

// GetValue implements IDatabase.
func (d Database) GetValue(namespace string, key string) (value []byte, err error) {
	return d.GetW(toByte(namespace), toByte(key))
}

func (d Database) GetW(namespace, key []byte) (value []byte, err error) {
	err = d.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(toByte(about.Database_Namespace_Root)).Bucket(namespace)
		tmpValue := bucket.Get(key)
		value = make([]byte, len(tmpValue))
		copy(value, tmpValue)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return value, nil
}

// SetValue implements IDatabase.
func (d Database) SetValue(namespace string, key string, value []byte) error {
	return d.SetW(toByte(namespace), toByte(key), value)
}

func (d Database) SetW(namespace, key, value []byte) error {
	var mu sync.Mutex
	mu.Lock()
	err := d.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(toByte(about.Database_Namespace_Root)).Bucket(namespace)
		err := bucket.Put(key, value)
		return err
	})
	mu.Unlock()
	if err != nil {
		d.log.Debugf("Failed to set key %s for namespace %s: %v", key, namespace, err)
		return err
	}

	return nil
}

// DelValue implements IDatabase.
func (d Database) DelValue(namespace, key string) error {
	return d.DelW(toByte(namespace), toByte(key))
}

func (d Database) DelW(namespace, key []byte) error {
	var mu sync.Mutex
	mu.Lock()
	err := d.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(toByte(about.Database_Namespace_Root)).Bucket(namespace)
		err := bucket.Delete(key)
		return err
	})
	mu.Unlock()
	if err != nil {
		d.log.Debugf("Failed del key %s for namespace %s: %v", key, namespace, err)
		return err
	}

	return nil
}

// HasValue implements IDatabase.
func (d Database) HasValue(namespace, key string) (ok bool, err error) {
	_, err = d.GetValue(namespace, key)
	switch err {
	case nil:
		return true, nil
	}

	return false, nil
}

func toByte(str string) []byte {
	return []byte(str)
}

// CreateApplication implements IDatabase.
func (d Database) CreateApplication(application string) error {
	var mu sync.Mutex
	mu.Lock()
	err := d.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.Bucket(toByte(about.Database_Namespace_Root)).CreateBucketIfNotExists(toByte(application))
		return err
	})
	mu.Unlock()
	if err != nil {
		d.log.Errorf("Could not create bucket for application: %v", err)
	}
	return err
}

// DeleteApplication implements IDatabase.
func (d Database) DeleteApplication(application string) error {
	var mu sync.Mutex
	mu.Lock()
	err := d.db.Update(func(tx *bolt.Tx) error {
		err := tx.Bucket(toByte(about.Database_Namespace_Root)).DeleteBucket(toByte(application))
		if err != nil {
			d.log.Errorf("Could not delete bucket for application: %v", err)
		}
		return err
	})
	mu.Unlock()
	return err
}

// SetApplicationData implements IDatabase.
func (d Database) SetApplicationData(application string, item core.DiscoveryItem) error {
	d.CreateApplication(application)
	var mu sync.Mutex
	mu.Lock()
	err := d.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(toByte(about.Database_Namespace_Root)).Bucket(toByte(application))
		data, err := json.Marshal(item)
		if err != nil {
			d.log.Errorf("Could not set application data: %v", err)
			return err
		}
		err = bucket.Put(toByte("Value"), data)
		return err
	})
	mu.Unlock()
	return err
}

// Backup implements IDatabase.
func (d Database) Backup(backup io.Writer) error {
	return d.db.View(func(tx *bolt.Tx) error {
		_, err := tx.WriteTo(backup)
		return err
	})
}

// Close implements IDatabase.
func (d Database) Close() error {
	return d.db.Close()
}

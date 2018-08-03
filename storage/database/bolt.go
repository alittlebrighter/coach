// Copyright (c) 2018, Adam Bright <brightam1@gmail.com>
// See LICENSE for licensing information

package database

import (
	"errors"
	"time"

	"github.com/buger/jsonparser"
	bolt "github.com/coreos/bbolt"
	"github.com/json-iterator/go" // for full (de)serialization
	"github.com/rs/xid"

	"github.com/alittlebrighter/coach/gen/models"
)

var (
	HistoryBucket   = []byte("history")
	SavedCmdsBucket = []byte("commands")
	IgnoreBucket    = []byte("ignore")
	buckets         = [][]byte{
		HistoryBucket,
		SavedCmdsBucket,
		IgnoreBucket,
	}

	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

type BoltDB struct {
	db *bolt.DB
}

func NewBoltDB(path string, readonly bool) (db *BoltDB, err error) {
	b := new(BoltDB)
	b.db, err = bolt.Open(path, 0600, &bolt.Options{Timeout: 2 * time.Second, ReadOnly: readonly})
	return b, b.initDB()
}

// Close closes the bolt db file.
func (d *BoltDB) Close() {
	d.db.Close()
}

func (b *BoltDB) initDB() error {
	if b.db.IsReadOnly() {
		return nil
	}
	return b.db.Update(func(tx *bolt.Tx) error {
		for _, bucket := range buckets {
			if _, err := tx.CreateBucketIfNotExists(bucket); err != nil {
				return err
			}
		}
		return nil
	})
}

func (b *BoltDB) CheckDupeCmds(command string, count int) (countReached bool) {
	b.db.View(func(tx *bolt.Tx) error {
		if shouldIgnore(tx, command) {
			countReached = false
			return nil
		}

		// Assume bucket exists and has keys
		c := tx.Bucket(HistoryBucket).Cursor()

		for k, v := c.Last(); count > 0 && k != nil; k, v = c.Prev() {
			if fullCmd, err := jsonparser.GetUnsafeString(v, "fullCommand"); err == nil && fullCmd == command {
				count--
			}
		}

		countReached = count <= 0
		return nil
	})
	return
}

// GetRecent retrieves the last count (arg) lines of history from specified tty (arg).
func (b *BoltDB) GetRecent(tty string, count int) ([]models.HistoryRecord, error) {
	records := []models.HistoryRecord{}
	b.db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		c := tx.Bucket(HistoryBucket).Cursor()

		for k, v := c.Last(); len(records) < count && k != nil; k, v = c.Prev() {
			if lineTty, err := jsonparser.GetUnsafeString(v, "tty"); err != nil || lineTty != tty {
				continue
			}

			var line models.HistoryRecord
			var err error
			err = json.Unmarshal(v, &line)
			if err != nil {
				continue
			}

			records = append([]models.HistoryRecord{line}, records...)
		}

		return nil
	})

	return records, nil
}

func (b *BoltDB) PruneHistory(max int) error {
	if max < 1 {
		return errors.New("invalid max value")
	}
	return b.db.Update(func(tx *bolt.Tx) error {
		diff := tx.Bucket(HistoryBucket).Stats().KeyN - max

		c := tx.Bucket(HistoryBucket).Cursor()

		for k, _ := c.First(); diff > 0; k, _ = c.Next() {
			if err := tx.Bucket(HistoryBucket).Delete(k); err != nil {
				return err
			}
			diff--
		}
		return nil
	})
}

func (b *BoltDB) QueryScripts(tags ...string) ([]models.DocumentedScript, error) {
	cmds := []models.DocumentedScript{}
	if len(tags) == 0 {
		return cmds, nil
	}
	err := b.db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		c := tx.Bucket(SavedCmdsBucket).Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			// this is an ugly abuse of scoping rules
			var savedCmd models.DocumentedScript
			skip := false
			jsonparser.ArrayEach(v, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
				if skip {
					return
				}
				// simple for now, just add command to the list if any tags match
				for _, tag := range tags {
					if tag == string(value) {
						err = json.Unmarshal(v, &savedCmd)
						if err != nil {
							skip = true
							break
						}
					}
				}
			}, "tags")
			cmds = append(cmds, savedCmd)
		}

		return nil
	})
	return cmds, err
}

func (b *BoltDB) GetScript(alias []byte) (command *models.DocumentedScript) {
	b.db.View(func(tx *bolt.Tx) error {
		cmdData := tx.Bucket(SavedCmdsBucket).Get(alias)
		if cmdData == nil || len(cmdData) == 0 {
			return errors.New("not found")
		}

		return json.Unmarshal(cmdData, &command)
	})
	return
}

func (b *BoltDB) DeleteScript(alias []byte) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(SavedCmdsBucket).Delete(alias)
	})
}

func (b *BoltDB) IgnoreCommand(command string) (err error) {
	err = b.db.Update(func(tx *bolt.Tx) error {
		return ignoreCommand(tx, command)
	})
	return
}

func (b *BoltDB) UnignoreCommand(command string) (err error) {
	err = b.db.Update(func(tx *bolt.Tx) error {
		return unignoreCommand(tx, command)
	})
	return
}

func ignoreCommand(tx *bolt.Tx, command string) (err error) {
	err = tx.Bucket(IgnoreBucket).Put([]byte(command), []byte{})
	return
}

func unignoreCommand(tx *bolt.Tx, command string) (err error) {
	err = tx.Bucket(IgnoreBucket).Delete([]byte(command))
	return
}

func shouldIgnore(tx *bolt.Tx, command string) (yes bool) {
	yes = tx.Bucket(IgnoreBucket).Get([]byte(command)) != nil
	return
}

var ErrAlreadyExists = errors.New("already exists")

func (b *BoltDB) Save(id []byte, instance interface{}, overwrite bool) (err error) {
	var bucket []byte
	switch instance.(type) {
	case models.HistoryRecord:
		bucket = HistoryBucket
	case models.DocumentedScript:
		bucket = SavedCmdsBucket
	default:
		bucket = IgnoreBucket
	}

	if id == nil || len(id) == 0 {
		id = xid.New().Bytes()
	}
	return b.db.Update(func(tx *bolt.Tx) error {
		if !overwrite && tx.Bucket(bucket).Get(id) != nil {
			return ErrAlreadyExists
		}
		return saveToBucket(tx, bucket, id, instance)
	})
}

func saveToBucket(tx *bolt.Tx, bucket []byte, id []byte, val interface{}) error {
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}
	return tx.Bucket(bucket).Put(id, data)
}

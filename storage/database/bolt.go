// Copyright (c) 2018, Adam Bright <brightam1@gmail.com>
// See LICENSE for licensing information

package database

import (
	"errors"
	"strings"
	"time"

	"github.com/buger/jsonparser"
	bolt "github.com/coreos/bbolt"
	"github.com/json-iterator/go"
	"github.com/rs/xid"

	models "github.com/alittlebrighter/coach/gen/proto"
)

const (
	Wildcard  = "?"
	TrashTag  = "coach.trash.983476" // just something arbitrary that is unlikely to be used by anything else
	FilePerms = 0660
)

var (
	HistoryBucket    = []byte("history")
	SavedCmdsBucket  = []byte("commands")
	IgnoreBucket     = []byte("ignore")
	IgnoreWordBucket = []byte("ignore-word")
	buckets          = [][]byte{
		HistoryBucket,
		SavedCmdsBucket,
		IgnoreBucket,
		IgnoreWordBucket,
	}

	json = jsoniter.ConfigCompatibleWithStandardLibrary

	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
)

type BoltDB struct {
	db *bolt.DB
}

func NewBoltDB(path string, readonly bool) (b *BoltDB, err error) {
	b = new(BoltDB)
	b.db, err = bolt.Open(path, FilePerms, &bolt.Options{Timeout: 2 * time.Second, ReadOnly: readonly})
	return
}

// Close closes the bolt db file.
func (b *BoltDB) Close() error {
	return b.db.Close()
}

func (b *BoltDB) Init() error {
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
		command = CleanseCommand(command)

		// Assume bucket exists and has keys
		c := tx.Bucket(HistoryBucket).Cursor()

		for k, v := c.Last(); count > 0 && k != nil; k, v = c.Prev() {
			fullCmd, err := jsonparser.GetUnsafeString(v, "fullCommand")
			if err == nil && fullCmd == command {
				count--
			}
		}

		countReached = count <= 0
		return nil
	})
	return
}

// CleanseCommand converts a command to what it would look like in the database
func CleanseCommand(command string) string {
	return strings.Replace(command, `"`, `\"`, -1)
}

// GetRecent retrieves the last count (arg) lines of history from specified tty (arg).
func (b *BoltDB) GetRecent(tty, username string, count int) ([]models.HistoryRecord, error) {
	records := []models.HistoryRecord{}
	b.db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		c := tx.Bucket(HistoryBucket).Cursor()

		for k, v := c.Last(); len(records) < int(count) && k != nil; k, v = c.Prev() {
			if tty != Wildcard {
				if lineTty, err := jsonparser.GetUnsafeString(v, "tty"); err != nil || lineTty != tty {
					continue
				}
			}
			if lineUser, err := jsonparser.GetUnsafeString(v, "user"); err != nil || lineUser != username {
				continue
			}

			var line models.HistoryRecord
			if err := json.Unmarshal(v, &line); err != nil {
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
	all := false
	for _, tag := range tags {
		if tag == Wildcard {
			all = true
			break
		}
	}

	err := b.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(SavedCmdsBucket).Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var savedCmd models.DocumentedScript
			shouldAdd := false
			if all && len(k) > 0 {
				shouldAdd = true
				goto add
			}

			for _, tag := range tags {
				tag = strings.ToLower(tag)

				if tag[0] == '~' && len(tag) > 1 && strings.Contains(strings.ToLower(string(k)), tag[1:]) ||
					strings.ToLower(string(k)) == tag {
					shouldAdd = true
					goto add
				}

				jsonparser.ArrayEach(v, func(scriptTag []byte, dataType jsonparser.ValueType, offset int, err error) {
					switch {
					case shouldAdd || err != nil:
						return
					case tag[0] == '~' && len(tag) > 1 && strings.Contains(strings.ToLower(string(scriptTag)), tag[1:]):
						fallthrough
					case strings.ToLower(string(scriptTag)) == tag:
						shouldAdd = true
					}
				}, "tags")
			}

		add:
			if shouldAdd {
				if err := json.Unmarshal(v, &savedCmd); err == nil {
					cmds = append(cmds, savedCmd)
				}
			}
		}

		return nil
	})
	return cmds, err
}

func (b *BoltDB) GetScript(alias []byte) (command *models.DocumentedScript) {
	b.db.View(func(tx *bolt.Tx) error {
		cmdData := tx.Bucket(SavedCmdsBucket).Get(alias)
		if cmdData == nil || len(cmdData) == 0 {
			return ErrNotFound
		}

		return json.Unmarshal(cmdData, &command)
	})
	return
}

func (b *BoltDB) DeleteScript(alias []byte) error {
	return b.db.Batch(func(tx *bolt.Tx) error {
		return tx.Bucket(SavedCmdsBucket).Delete(alias)
	})
}

func (b *BoltDB) IgnoreWord(word, username string) (err error) {
	w, u := []byte(word), []byte(username)
	err = b.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.Bucket(IgnoreWordBucket).CreateBucketIfNotExists(u)
		if err != nil {
			return err
		}
		return bucket.Put(w, []byte{})
	})
	return
}

func (b *BoltDB) UnignoreWord(word, username string) (err error) {
	w, u := []byte(word), []byte(username)
	err = b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(IgnoreWordBucket).Bucket(u)
		if bucket == nil {
			return nil
		}
		return bucket.Delete(w)
	})
	return
}

func (b *BoltDB) IgnoreCommand(command, username string) (err error) {
	c, u := []byte(command), []byte(username)
	err = b.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.Bucket(IgnoreBucket).CreateBucketIfNotExists(u)
		if err != nil {
			return err
		}
		return bucket.Put(c, []byte{})
	})
	return
}

func (b *BoltDB) UnignoreCommand(command, username string) (err error) {
	c, u := []byte(command), []byte(username)
	err = b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(IgnoreBucket).Bucket(u)
		if bucket == nil {
			return nil
		}
		return bucket.Delete(c)
	})
	return
}

func (b *BoltDB) ShouldIgnoreCommand(command, username string) (yes bool) {
	c, u := []byte(command), []byte(username)
	b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(IgnoreBucket).Bucket(u)
		if bucket == nil {
			return nil
		}
		if bucket.Get(c) != nil {
			yes = true
		}
		return nil
	})
	return
}

func (b *BoltDB) ShouldIgnoreWord(word, username string) (yes bool) {
	w, u := []byte(word), []byte(username)
	b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(IgnoreWordBucket).Bucket(u)
		if bucket == nil {
			return nil
		}
		if bucket.Get(w) != nil {
			yes = true
		}
		return nil
	})
	return
}

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
	return b.db.Batch(func(tx *bolt.Tx) error {
		if !overwrite && tx.Bucket(bucket).Get(id) != nil {
			return ErrAlreadyExists
		}
		return saveToBucket(tx, bucket, id, instance)
	})
}

func (b *BoltDB) SaveBatch(toSave <-chan HasID, bucket []byte) <-chan error {
	errs := make(chan error)

	go func() {
		b.db.Batch(func(tx *bolt.Tx) error {
			for instance := range toSave {
				saveToBucket(tx, bucket, instance.GetId(), instance)
			}
			return nil
		})
	}()

	return errs
}

func saveToBucket(tx *bolt.Tx, bucket []byte, id []byte, val interface{}) error {
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}
	return tx.Bucket(bucket).Put(id, data)
}

type HasID interface {
	GetId() []byte
}

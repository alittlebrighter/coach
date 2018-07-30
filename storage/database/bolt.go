package database

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/buger/jsonparser"

	bolt "github.com/coreos/bbolt"
	"github.com/rs/xid"

	// "github.com/buger/jsonparser" // for queries
	"github.com/json-iterator/go" // for full (de)serialization

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

func (d *BoltDB) SaveHistory(c *models.HistoryRecord) (err error) {
	if c.Id == nil || len(c.Id) == 0 {
		c.Id = xid.New().Bytes()
	}
	return d.db.Update(func(tx *bolt.Tx) error {
		return saveToBucket(tx, HistoryBucket, c)
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

// GetRecent retrieves the last count (arg) lines of history
func (b *BoltDB) GetRecent(tty string, count int) ([]models.HistoryRecord, error) {
	records := []models.HistoryRecord{}
	b.db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		c := tx.Bucket(HistoryBucket).Cursor()

		for k, v := c.Last(); len(records) <= count && k != nil; k, v = c.Prev() {
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

func (b *BoltDB) SaveDoc(sc *models.SavedCommand) (err error) {
	if sc.GetCommand() == nil || len(sc.GetCommand()) == 0 {
		return errors.New("no command attached to docs")
	}

	switch {
	case len(sc.Alias) > 0:
		sc.Id = []byte(sc.Alias)
	case sc.Id == nil || len(sc.Id) == 0:
		sc.Id = xid.New().Bytes()
	}

	err = b.db.Update(func(tx *bolt.Tx) error {
		err := saveToBucket(tx, SavedCmdsBucket, sc)
		if err != nil {
			return err
		}
		for _, command := range sc.GetCommand() {
			fullCommand := strings.Join(append([]string{command.GetCommand()}, command.GetArguments()...), " ")
			err = ignoreCommand(tx, fullCommand)
		}
		return err
	})
	return
}

func (b *BoltDB) QueryDoc(tags ...string) ([]models.SavedCommand, error) {
	cmds := []models.SavedCommand{}
	if len(tags) == 0 {
		return cmds, nil
	}
	err := b.db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		c := tx.Bucket(SavedCmdsBucket).Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			// this is an ugly abuse of scoping rules
			var savedCmd models.SavedCommand
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

func (b *BoltDB) GetSavedCmd(alias string) (command *models.SavedCommand) {
	b.db.View(func(tx *bolt.Tx) error {
		cmdData := tx.Bucket(SavedCmdsBucket).Get([]byte(alias))
		if cmdData == nil || len(cmdData) == 0 {
			fmt.Println("alias not found")
			return errors.New("not found")
		}

		return json.Unmarshal(cmdData, &command)
	})
	return
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

func saveToBucket(tx *bolt.Tx, bucket []byte, val HasId) error {
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}
	return tx.Bucket(bucket).Put(val.GetId(), data)
}

type HasId interface {
	GetId() []byte
}

package database

import (
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
	buckets         = [][]byte{
		HistoryBucket,
		SavedCmdsBucket,
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

func (d *BoltDB) Save(c *models.HistoryRecord) (err error) {
	if c.Id == nil || len(c.Id) == 0 {
		c.Id = xid.New().Bytes()
	}
	return d.db.Update(saveTxn(HistoryBucket, c))
}

// GetRecent)
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

func (d *BoltDB) Close() {
	d.db.Close()
}

func saveTxn(bucket []byte, val HasId) func(tx *bolt.Tx) error {
	return func(tx *bolt.Tx) error {
		data, err := json.Marshal(val)
		if err != nil {
			return err
		}
		tx.Bucket(bucket).Put(val.GetId(), data)
		return nil
	}
}

type HasId interface {
	GetId() []byte
}

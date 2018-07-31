package coach

import (
	"errors"
	"strings"

	"github.com/golang/protobuf/ptypes"

	"github.com/alittlebrighter/coach/gen/models"
	"github.com/alittlebrighter/coach/platforms"
)

func SaveHistory(line string, dupeCount int, store HistoryStore) (promptDoc bool, err error) {
	cmd := platforms.CleanupCommand(line)
	if len(cmd) == 0 {
		if lines, err := Shell.History(1); err == nil && len(lines) > 0 {
			cmd = lines[0]
		}
	}

	hLine := models.HistoryRecord{
		Id:          RandomID(),
		Timestamp:   ptypes.TimestampNow(),
		FullCommand: cmd,
		Tty:         Shell.GetTTY(),
	}
	if len(hLine.GetFullCommand()) == 0 || strings.HasPrefix(hLine.GetFullCommand(), "coach") {
		return
	}

	// size of one (1) just so we don't block
	enoughDupes := make(chan bool, 1)
	if dupeCount > 0 {
		go func() {
			enoughDupes <- store.CheckDupeCmds(hLine.GetFullCommand(), dupeCount-1)
		}()
	} else {
		enoughDupes <- false
	}

	err = store.Save(hLine.GetId(), hLine)

	promptDoc = <-enoughDupes

	return
}

func GetRecentHistory(n int, store HistoryStore) (lines []models.HistoryRecord, err error) {
	if n <= 0 {
		err = errors.New("invalid input")
	}

	lines, err = store.GetRecent(Shell.GetTTY(), n)
	return
}

type HistoryStore interface {
	Save(id []byte, value interface{}) error
	CheckDupeCmds(string, int) bool
	GetRecent(tty string, n int) ([]models.HistoryRecord, error)
}

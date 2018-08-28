// Copyright (c) 2018, Adam Bright <brightam1@gmail.com>
// See LICENSE for licensing information

package coach

import (
	"errors"
	"fmt"
	"os/user"
	"strings"

	"github.com/golang/protobuf/ptypes"
	"github.com/spf13/viper"

	"github.com/alittlebrighter/coach-pro/gen/models"
	"github.com/alittlebrighter/coach-pro/platforms"
	"github.com/alittlebrighter/coach-pro/storage/database"
)

func ImportHistory(store HistoryStore) error {
	user, err := user.Current()
	if err != nil {
		return err
	}
	now := ptypes.TimestampNow()
	tty := platforms.GetTTY()
	hLine := models.HistoryRecord{
		Timestamp: now,
		Tty:       tty,
		User:      user.Username,
	}

	lines, err := platforms.NativeHistory(viper.GetInt("history.maxlines"))
	if err != nil {
		return err
	}

	toSave := make(chan database.HasID)
	go func() {
		for err := range store.SaveBatch(toSave, database.HistoryBucket) {
			fmt.Println("ERROR:", err)
		}
	}()
	for line := range lines {
		if len(line) == 0 || strings.HasPrefix(line, "coach") ||
			strings.HasPrefix(line, "history") || strings.HasPrefix(line, "exit") {
			continue
		}

		hLine.Id = RandomID()
		hLine.FullCommand = line

		toSave <- &hLine
	}
	close(toSave)

	return nil
}

func SaveHistory(line string, dupeCount int, store HistoryStore) (promptDoc bool, err error) {
	cmd := platforms.CleanupCommand(line)
	if len(cmd) == 0 {
		if lines, err := platforms.NativeHistory(1); err == nil && len(lines) > 0 {
			for line := range lines {
				cmd = line
			}
		}
	}

	user, err := user.Current()
	if err != nil {
		return false, err
	}

	hLine := models.HistoryRecord{
		Id:          RandomID(),
		Timestamp:   ptypes.TimestampNow(),
		FullCommand: cmd,
		Tty:         platforms.GetTTY(),
		User:        user.Username,
	}
	if len(hLine.GetFullCommand()) == 0 || strings.HasPrefix(hLine.GetFullCommand(), "coach") {
		return
	}

	if err = store.Save(hLine.GetId(), hLine, true); err != nil {
		return
	}

	enoughDupes := make(chan bool, 1)
	ignore := make(chan bool, 1)
	if dupeCount > 0 {
		go func() {
			enoughDupes <- store.CheckDupeCmds(hLine.GetFullCommand(), dupeCount)
		}()
		go func() {
			ignore <- ShouldIgnore(hLine.GetFullCommand(), store)
		}()
	} else {
		enoughDupes <- false
	}

	var promptDocSet, ignoreSet, shouldIgnore bool
	for {
		select {
		case promptDoc = <-enoughDupes:
			promptDocSet = true
		case shouldIgnore = <-ignore:
			ignoreSet = true
		}

		if shouldIgnore {
			promptDoc = false
			break
		} else if promptDocSet && ignoreSet {
			break
		}
	}

	store.PruneHistory(viper.GetInt("history.maxlines"))
	return
}

func GetRecentHistory(n int, allSessions bool, store HistoryGetter) (lines []models.HistoryRecord, err error) {
	if n == 0 {
		err = errors.New("invalid input")
		return
	}

	var currentUser *user.User
	currentUser, err = user.Current()
	if err != nil {
		return
	}

	var tty string
	if allSessions {
		tty = database.Wildcard
	} else {
		tty = platforms.GetTTY()
	}

	lines, err = store.GetRecent(tty, currentUser.Username, n)
	return
}

type HistoryStore interface {
	Save(id []byte, value interface{}, overwrite bool) error
	SaveBatch(<-chan database.HasID, []byte) <-chan error
	CheckDupeCmds(string, int) bool
	PruneHistory(max int) error
	HistoryGetter
	IgnoreChecker
}

type HistoryGetter interface {
	GetRecent(tty string, username string, n int) ([]models.HistoryRecord, error)
}

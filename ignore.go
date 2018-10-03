package coach

import (
	"os/user"
	"strings"
	"time"

	models "github.com/alittlebrighter/coach/gen/proto"
	"github.com/alittlebrighter/coach/platforms"
)

func IgnoreHistory(lineCount int, allVariations, remove bool, store IgnoreStore) error {
	processLine := func(l string) string {
		return l
	}
	toRun := store.IgnoreCommand

	if allVariations {
		processLine = func(l string) string {
			return strings.Fields(l)[0]
		}
		toRun = store.IgnoreWord
	}
	if remove && allVariations {
		toRun = store.UnignoreWord
	} else if remove {
		toRun = store.UnignoreCommand
	}
	return toggleIgnore(lineCount, store, processLine, toRun)
}

func toggleIgnore(lineCount int, store HistoryGetter, processLine func(string) string, run func(string, string) error) error {
	currentUser, err := user.Current()
	if err != nil {
		return err
	}

	var hLines []models.HistoryRecord
	hLines, err = GetRecentHistory(lineCount, false, store)
	if err != nil {
		return err
	}
	for _, line := range hLines {
		run(processLine(line.GetFullCommand()), currentUser.Username)
	}
	return nil
}

func IgnoreCommand(command string, store IgnoreStore) error {
	currentUser, err := user.Current()
	if err != nil {
		return err
	}
	return store.IgnoreCommand(command, currentUser.Username)
}

func ShouldIgnore(command string, store IgnoreChecker) bool {
	if len(strings.TrimSpace(command)) == 0 {
		return true
	}

	user, _ := user.Current()

	cmdResponse := make(chan bool, 1)
	wordResponse := make(chan bool, 1)

	if platforms.IsCompoundStatement(command) {
		wordResponse <- false
	} else {
		go func() {
			wordResponse <- store.ShouldIgnoreWord(strings.Fields(command)[0], user.Username)
		}()
	}
	go func() {
		cmdResponse <- store.ShouldIgnoreCommand(command, user.Username)
	}()

	var word, cmd, wordSet, cmdSet bool
	timer := time.NewTimer(1 * time.Second)
	for {
		select {
		case word = <-wordResponse:
			wordSet = true
		case cmd = <-cmdResponse:
			cmdSet = true
		case <-timer.C:
			break
		}

		if word || cmd {
			return true
		} else if wordSet && cmdSet && !word && !cmd {
			return false
		}
	}
	return false
}

type IgnoreStore interface {
	IgnoreWord(word, username string) (err error)
	UnignoreWord(word, username string) (err error)
	IgnoreCommand(command, username string) (err error)
	UnignoreCommand(command, username string) (err error)
	IgnoreChecker
	HistoryGetter
}

type IgnoreChecker interface {
	ShouldIgnoreCommand(command, username string) (yes bool)
	ShouldIgnoreWord(word, username string) (yes bool)
}

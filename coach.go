package coach

import (
	"fmt"

	"github.com/alittlebrighter/coach/platforms/linux"
	"github.com/alittlebrighter/coach/types"
)

type CommandFetcher func(uint) ([]types.Cmd, error)

var fetch CommandFetcher

type HistoryParser func(string) (*types.Cmd, error)

var historyParser HistoryParser

func init() {
	// only supporting bash on linux so only one option for fetcher
	fetch = linux.LastCommands
	historyParser = linux.ParseHistory
}

func LastCommands(n uint) (*types.Cmd, error) {
	cmds, err := fetch(n)
	return &cmds[0], err
}

func ParseHistory(historyOut string) (*types.Cmd, error) {
	fmt.Println("history output arg:", historyOut)
	return historyParser(historyOut)
}

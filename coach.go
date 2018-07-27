//go:generate protoc -I protobuf/ --go_out=plugins=grpc:gen/models protobuf/coach.proto
package coach

import (
	"github.com/alittlebrighter/coach/gen/models"
	"github.com/alittlebrighter/coach/platforms/linux"
)

type (
	CommandFetcher func(uint) ([]models.Command, error)
	HistoryParser  func(string) (*models.Command, error)
)

var (
	fetch         CommandFetcher
	historyParser HistoryParser
)

func init() {
	// only supporting bash on linux so only one option for fetcher
	fetch = linux.LastCommands
	historyParser = linux.ParseHistory
}

func LastCommands(n uint) (*models.Command, error) {
	cmds, err := fetch(n)
	return &cmds[0], err
}

func ParseHistory(historyOut string) (*models.Command, error) {
	return historyParser(historyOut)
}

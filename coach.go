//go:generate protoc -I protobuf/ --go_out=plugins=grpc:gen/models protobuf/coach.proto
package coach

import (
	"strconv"
	"strings"

	"github.com/alittlebrighter/coach/gen/models"
	"github.com/alittlebrighter/coach/platforms/linux"
)

type (
	Platform interface {
		History(lineCount int) (lines []string, err error)
	}

	CommandFetcher func(uint) ([]models.Command, error)
	HistoryParser  func(string) (*models.Command, error)
)

var (
	fetch         CommandFetcher
	historyParser HistoryParser
)

func init() {
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

func CleanupCommand(cmd string) (clean string) {
	var err error
	parts := strings.Fields(cmd)
	historyLineNum := 0
	if len(parts) > 0 {
		historyLineNum, err = strconv.Atoi(parts[0])

	}
	switch {
	case (historyLineNum > 0 || err == nil) && len(parts) > 1:
		parts = parts[1:]
	case historyLineNum > 0 || err == nil:
		parts = []string{}
	}
	clean = strings.TrimSpace(strings.Join(parts, " "))
	return
}

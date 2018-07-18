package linux

import (
	"errors"
	"strings"

	"github.com/alittlebrighter/coach/types"
)

func LastCommands(n uint) (cmds []types.Cmd, err error) {
	return
}

func ParseHistory(historyOutput string) (cmd *types.Cmd, err error) {
	fields := strings.Fields(historyOutput)
	if len(fields) < 2 {
		err = errors.New("insufficient output to parse")
	}
	cmd = &types.Cmd{
		Command: fields[1],
	}
	if len(historyOutput) > 2 {
		cmd.Args = fields[2:]
	}
	return
}

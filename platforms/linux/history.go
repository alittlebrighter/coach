package linux

import (
	"errors"
	"strings"

	"github.com/alittlebrighter/coach/gen/models"
)

func LastCommands(n uint) (cmds []models.Command, err error) {
	return
}

func ParseHistory(historyOutput string) (cmd *models.Command, err error) {
	fields := strings.Fields(historyOutput)
	if len(fields) < 2 {
		err = errors.New("insufficient output to parse")
	}
	cmd = &models.Command{
		Command: fields[1],
	}
	if len(historyOutput) > 2 {
		cmd.Arguments = fields[2:]
	}
	return
}

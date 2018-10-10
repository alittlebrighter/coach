// Copyright (c) 2018, Adam Bright <brightam1@gmail.com>
// See LICENSE for licensing information

package platforms

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

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

type Bash struct{}

func (b *Bash) History(lineCount int) <-chan string {
	historyCmd := exec.Command("bash", "-i")
	historyCmd.Stdin = bytes.NewBuffer([]byte(fmt.Sprintf("history %d; exit;", lineCount)))
	output, _ := historyCmd.StdoutPipe()
	historyCmd.Start()
	lines := make(chan string)

	go func() {
		line := make([]byte, 128)
		command := ""
		for _, err := output.Read(line); err == nil; _, err = output.Read(line) {
			command += string(line)

			for {
				newlineIndex := strings.Index(command, "\n")
				if newlineIndex == -1 {
					break
				}

				lines <- CleanupCommand(command[:newlineIndex])

				if newlineIndex == len(command)-1 {
					command = ""
				} else {
					command = command[newlineIndex+1:]
				}
			}
		}
		close(lines)
		historyCmd.Wait()
	}()

	return lines
}

func (b *Bash) BuildCommand(ctx context.Context, script string, args []string) (*exec.Cmd, func(), error) {
	cmdArgs := append([]string{"-c", "( " + script + " )" /* HACK */, "''" /* END HACK */}, args...)
	return exec.CommandContext(ctx, "bash", cmdArgs...), nil, nil
}

func (b *Bash) FileExtension() string {
	return "sh"
}

func (b *Bash) LineComment() string {
	return "#"
}

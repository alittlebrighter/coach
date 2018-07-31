package platforms

import (
	"bytes"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Shell interface {
	History(lineCount int) (lines []string, err error)
	GetTTY() string
	GetPWD() string
	BuildCommand(script string) *exec.Cmd
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

type Bash struct{}

func (b *Bash) History(lineCount int) (lines []string, err error) {
	output := bytes.NewBuffer([]byte{})
	historyCmd := exec.Command("history", strconv.Itoa(lineCount))
	historyCmd.Stdout = output
	historyCmd.Run()

	for line, err := output.ReadString('\n'); err == nil; line, err = output.ReadString('\n') {
		lines = append(lines, CleanupCommand(line))
	}
	return
}

func (b *Bash) GetTTY() string {
	ttyCmd := exec.Command("tty")
	ttyCmd.Stdin = os.Stdin
	ttyBytes, _ := ttyCmd.Output()
	return strings.TrimSpace(string(ttyBytes))
}

func (b *Bash) GetPWD() string {
	path, _ := os.Getwd()
	return path
}

func (b *Bash) BuildCommand(script string) *exec.Cmd {
	return exec.Command("bash", "-c", "( "+script+" )")
}

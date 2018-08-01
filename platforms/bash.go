// Copyright (c) 2018, Adam Bright <brightam1@gmail.com>
// See LICENSE for licensing information

package platforms

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/rs/xid"
)

type Shell interface {
	History(lineCount int) (lines []string, err error)
	GetTTY() string
	GetPWD() string
	BuildCommand(script string) *exec.Cmd
	CreateTmpFile(contents []byte) (string, error)
	OpenEditor(filepath string) error
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
	historyCmd.Stdin = os.Stdin
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

func (b *Bash) CreateTmpFile(contents []byte) (string, error) {
	filename := "/tmp/coach-" + xid.New().String() + ".sh"
	return filename, ioutil.WriteFile(filename, contents, 0600)
}

func (b *Bash) OpenEditor(filename string) error {
	editor := os.Getenv("EDITOR")
	if len(editor) == 0 {
		editor = "nano"
	}
	return exec.Command(editor, filename).Run()
}

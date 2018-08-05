package platforms

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type Platform interface {
	History(lineCount int) (lines []string, err error)
	GetTTY() string
	GetPWD() string
	CreateTmpFile(contents []byte) (string, error)
	OpenEditor(filepath string) error
}

type Shell interface {
	BuildCommand(script string) (*exec.Cmd, func(), error)
}

func IdentifyShell() string {
	output, err := exec.Command("readlink", "/bin/sh").Output()
	if err != nil {
		output = []byte(DefaultShell)
	}
	return strings.TrimSpace(string(output))
}

func WriteTmpFile(script string) (string, func(), error) {
	tmpfile, err := ioutil.TempFile("", "coach")
	if err != nil {
		return "", nil, err
	}
	defer tmpfile.Close()
	cleanup := func() { os.Remove(tmpfile.Name()) }

	if _, err := tmpfile.Write([]byte(script)); err != nil {
		cleanup()
		return "", nil, err
	}

	return tmpfile.Name(), cleanup, nil
}

type AnyShell struct {
	Name string
}

func (a *AnyShell) BuildCommand(script string) (*exec.Cmd, func(), error) {
	filename, cleanup, err := WriteTmpFile(script)
	if err != nil {
		return nil, nil, err
	}

	return exec.Command(a.Name, filename), cleanup, nil
}

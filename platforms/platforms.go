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
	output, _ := exec.Command("echo", `"$0"`).Output()
	return strings.TrimSpace(string(output))
}

func GetShell(name string) Shell {
	if len(name) == 0 {
		name = IdentifyShell()
	}

	switch {
	case name == "dash":
		fallthrough
	case name == "bash":
		return new(Bash)
	default:
		return &AnyShell{Name: name}
	}
}

type AnyShell struct {
	Name string
}

func (a *AnyShell) BuildCommand(script string) (*exec.Cmd, func(), error) {
	tmpfile, err := ioutil.TempFile("", "coach")
	if err != nil {
		return nil, nil, err
	}
	defer tmpfile.Close()
	cleanup := func() { os.Remove(tmpfile.Name()) }

	if _, err := tmpfile.Write([]byte(script)); err != nil {
		cleanup()
		return nil, nil, err
	}

	return exec.Command(a.Name, tmpfile.Name()), cleanup, nil
}

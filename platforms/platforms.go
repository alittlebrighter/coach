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
	FileExtension() string
	LineComment() string
}

func IdentifyShell() string {
	output, err := exec.Command("readlink", "/bin/sh").Output()
	if err != nil {
		output = []byte(DefaultShell)
	}
	return strings.TrimSpace(string(output))
}

func GetShell(name string) Shell {
	if len(name) == 0 {
		name = IdentifyShell()
	}

	platformShell := GetPlatformShell(name)
	switch {
	case platformShell != nil:
		return platformShell
	case name == "bash":
		return new(Bash)
	case strings.Contains(strings.ToLower(name), "python"):
		return &Python{AnyShell: &AnyShell{Name: name}}
	case strings.Contains(strings.ToLower(name), "ruby"):
		return &Ruby{AnyShell: &AnyShell{Name: name}}
	case strings.Contains(strings.ToLower(name), "node"):
		return &Node{AnyShell: &AnyShell{Name: name}}
	default:
		return &AnyShell{Name: name}
	}
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

func BuildCommand(interpreter, script string) (*exec.Cmd, func(), error) {
	filename, cleanup, err := WriteTmpFile(script)
	if err != nil {
		return nil, nil, err
	}

	return exec.Command(interpreter, filename), cleanup, nil
}

type AnyShell struct {
	Name string
}

func (a *AnyShell) BuildCommand(script string) (*exec.Cmd, func(), error) {
	return BuildCommand(a.Name, script)
}

func (a *AnyShell) LineComment() string {
	return "#"
}

func (a *AnyShell) FileExtension() string {
	return ".script"
}

type Python struct {
	*AnyShell
}

func (p *Python) FileExtension() string {
	return "py"
}

type Node struct {
	*AnyShell
}

func (n *Node) FileExtension() string {
	return "js"
}

func (n *Node) LineComment() string {
	return "//"
}

type Ruby struct {
	*AnyShell
}

func (r *Ruby) FileExtension() string {
	return "rb"
}

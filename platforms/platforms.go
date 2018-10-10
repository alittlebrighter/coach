package platforms

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type Shell interface {
	BuildCommand(ctx context.Context, script string, args []string) (*exec.Cmd, func(), error)
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

func OpenEditor(filename string) error {
	cmd := exec.Command(GetEditorCmd(), filename)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	return cmd.Run()
}

func NativeHistory(lineCount int) (<-chan string, error) {
	shell := GetShell("")

	switch shell.(type) {
	case *Bash:
		return new(Bash).History(lineCount), nil
	default:
		return nil, errors.New("not implemented")
	}
}

func ShellNameFromExt(fileExt string) string {
	switch fileExt {
	case "sh":
		return "bash"
	case "py":
		return "python"
	case "js":
		return "nodejs"
	case "rb":
		return "ruby"
	case "ps1":
		return "powershell"
	case "bat":
		return "windowsCMD"
	default:
		return ""
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

func BuildCommand(ctx context.Context, interpreter, script string, args []string) (*exec.Cmd, func(), error) {
	filename, cleanup, err := WriteTmpFile(script)
	if err != nil {
		return nil, nil, err
	}
	cmdArgs := append([]string{filename}, args...)
	return exec.CommandContext(ctx, interpreter, cmdArgs...), cleanup, nil
}

type AnyShell struct {
	Name string
}

func (a *AnyShell) BuildCommand(ctx context.Context, script string, args []string) (*exec.Cmd, func(), error) {
	return BuildCommand(ctx, a.Name, script, args)
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

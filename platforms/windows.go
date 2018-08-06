// +build windows

package platforms

import (
	"os"
	"os/exec"
	"strings"

	"github.com/rs/xid"
)

const DefaultShell = "windowsCMD"

func Newline(count uint) string {
	if count <= 1 {
		return "\r\n"
	}
	newlines := ""
	for i := uint(0); i < count; i++ {
		newlines += "\r\n"
	}
	return newlines
}

func GetEditorCmd() string {
	editor := os.Getenv("EDITOR")

	if len(editor) == 0 {
		editor = "notepad"
	}
	return editor
}

func GetShell(name string) Shell {
	if len(name) == 0 {
		name = IdentifyShell()
	}

	switch {
	case strings.Contains(strings.ToLower(name), "powershell"):
		return new(PowerShell)
	case strings.ToLower(name) == "windowscmd":
		return new(WindowsCMD)
	case name == "bash":
		return new(Bash)
	default:
		return &AnyShell{Name: name}
	}
}

type PowerShell struct{}

func (p *PowerShell) BuildCommand(script string) (*exec.Cmd, func(), error) {
	tmpfile, err := os.OpenFile(os.TempDir()+"/coach"+xid.New().String()+".ps1", os.O_CREATE, 0600)
	if err != nil {
		return nil, nil, err
	}
	defer tmpfile.Close()
	cleanup := func() { os.Remove(tmpfile.Name()) }

	if _, err := tmpfile.Write([]byte(script)); err != nil {
		cleanup()
		return nil, nil, err
	}

	return exec.Command("PowerShell.exe", tmpfile.Name()), cleanup, nil
}

type WindowsCMD struct{}

func (c *WindowsCMD) BuildCommand(script string) (*exec.Cmd, func(), error) {
	tmpfile, err := os.OpenFile(os.TempDir()+"/coach"+xid.New().String()+".bat", os.O_CREATE, 0600)
	if err != nil {
		return nil, nil, err
	}
	defer tmpfile.Close()
	cleanup := func() { os.Remove(tmpfile.Name()) }

	if _, err := tmpfile.Write([]byte(script)); err != nil {
		cleanup()
		return nil, nil, err
	}

	return exec.Command(tmpfile.Name()), cleanup, nil
}
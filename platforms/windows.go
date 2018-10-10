// +build windows

package platforms

import (
	"context"
	"os"
	"os/exec"
	"strings"

	"github.com/rs/xid"
)

const DefaultShell = "windowsCMD"

func DefaultHomeDir() string {
	return strings.Replace(os.Getenv("ProgramFiles"), "\\", "/", -1)
}

func GetTTY() string {
	return "windows"
}

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

// TODO
func IsCompoundStatement(command string) bool {
	return false
}

func GetPlatformShell(name string) Shell {
	switch {
	case strings.Contains(strings.ToLower(name), "powershell"):
		return new(PowerShell)
	case strings.ToLower(name) == "windowscmd":
		return new(WindowsCMD)
	default:
		return nil
	}
}

func KillProcess(cmd *exec.Cmd) error {
	return cmd.Process.Kill()
}

type PowerShell struct{}

func (p *PowerShell) BuildCommand(ctx context.Context, script string, args []string) (*exec.Cmd, func(), error) {
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
	cmdArgs := append([]string{"-NoProfile", "-ExecutionPolicy", "Bypass", "-OutputFormat", "Text",
		"-File", tmpfile.Name()}, args...)
	return exec.CommandContext(ctx, "PowerShell.exe", cmdArgs...), cleanup, nil
}

func (p *PowerShell) LineComment() string {
	return "#"
}

func (p *PowerShell) FileExtension() string {
	return "ps1"
}

type WindowsCMD struct{}

func (c *WindowsCMD) BuildCommand(ctx context.Context, script string, args []string) (*exec.Cmd, func(), error) {
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

	return exec.CommandContext(ctx, tmpfile.Name(), args...), cleanup, nil
}

func (c *WindowsCMD) LineComment() string {
	return "REM"
}

func (c *WindowsCMD) FileExtension() string {
	return ".bat"
}

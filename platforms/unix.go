// +build !windows

package platforms

import (
	"os"
	"os/exec"
	"strings"
)

func GetTTY() string {
	ttyCmd := exec.Command("tty")
	ttyCmd.Stdin = os.Stdin
	ttyBytes, _ := ttyCmd.Output()
	return strings.TrimSpace(string(ttyBytes))
}

func Newline(count uint) string {
	if count <= 1 {
		return "\n"
	}
	newlines := ""
	for i := uint(0); i < count; i++ {
		newlines += "\n"
	}
	return newlines
}

func GetEditorCmd() string {
	editor := os.Getenv("EDITOR")

	if len(editor) == 0 {
		editor = "nano"
	}
	return editor
}

func GetPlatformShell(name string) Shell {
	switch {
	case name == "dash":
		return new(Bash)
	default:
		return nil
	}
}

const DefaultShell = "bash"

func DefaultHomeDir() string {
	return "/usr/local"
}

func IsCompoundStatement(command string) bool {
	return strings.ContainsAny(command, ";&|<>")
}

// +build !windows

package platforms

import (
	"os"
)

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

const DefaultShell = "bash"

func DefaultHomeDir() string {
	return "/opt"
}

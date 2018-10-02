package config

import (
	"os"
	"strings"

	"github.com/alittlebrighter/coach-pro/platforms"
	homedir "github.com/mitchellh/go-homedir"
)

const (
	ENV_PREFIX = "coach"
)

func HomeDir() string {
	var home string
	envSetHome := os.Getenv("COACH_HOME")
	defaultAppDir := platforms.DefaultHomeDir()
	_, sysErr := os.Stat(defaultAppDir + "/coach")
	switch {
	case len(envSetHome) > 0:
		home = envSetHome
	case sysErr == nil:
		home = defaultAppDir + "/coach"
	default:
		homeDir, _ := homedir.Dir()
		home = homeDir + "/.coach"
	}
	return strings.Replace(home, `\`, "/", -1)
}

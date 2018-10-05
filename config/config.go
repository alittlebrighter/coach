package config

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/alittlebrighter/coach/platforms"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"
)

const (
	ENV_PREFIX = "coach"
)

func AppConfiguration() {
	viper.SetTypeByDefaultValue(true)

	viper.SetEnvPrefix(ENV_PREFIX)
	viper.AutomaticEnv()

	viper.AddConfigPath(HomeDir())
	viper.SetConfigName("config")
}

func WriteConfiguration() {
	// if no config file is found, write the defaults to one
	if err := viper.ReadInConfig(); err != nil {
		defaults := viper.AllSettings()
		data, _ := yaml.Marshal(&defaults)
		ioutil.WriteFile(HomeDir()+"/config.yaml", data, 0600)
	}
}

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

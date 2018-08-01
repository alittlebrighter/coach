// Copyright (c) 2018, Adam Bright <brightam1@gmail.com>
// See LICENSE for licensing information

package main

import (
	"fmt"
	"io/ioutil"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var (
	home, dbpath string
)

func main() {
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	home = home + "/.coach"
	dbpath = home + "/coach.db"

	os.Mkdir(home, os.ModePerm)

	if err != nil {
		fmt.Println("Could not access database.  ERROR:", err)
		os.Exit(1)
	}

	rootCmd := &cobra.Command{
		Use:   "coach",
		Short: "A tool to help you save and document common commands executed on the command line.",
		Run:   appMain,
	}
	/*
		sessionCmd := &cobra.Command{
			Use:   "session",
			Short: "Initialize a terminal session",
			Run:   session,
		}
	*/
	historyCmd := &cobra.Command{
		Use:   "history",
		Short: "Store and query command history.",
		Run:   history,
	}
	historyCmd.Flags().StringP("record", "r", "", "Record command.  With `bash` you can use \"$(history 1)\"")

	docCmd := &cobra.Command{
		Use:     "doc",
		Short:   "Save and query commands with tags and documentation.  Default is to document the most recent command.",
		Example: "coach doc [alias] [tags] [comment] # empty alias represented by \"\", tag list must be quoted if it contains spaces",
		Run:     doc,
	}
	docCmd.Flags().StringP("query", "q", "", "Query your saved commands by tags.")
	docCmd.Flags().StringP("script", "s", "", "Quoted command that you would like to document and save.")
	docCmd.Flags().StringP("edit", "e", "", "Edit the script specified by alias.")
	docCmd.Flags().IntP("history", "l", 1, "Number of most recent lines in history to put into the script.")

	ignoreCmd := &cobra.Command{
		Use:   "ignore",
		Short: "Ignore this command when scanning for duplicates to prompt for documentation.  Defaults to the last run command.",
		Run:   ignore,
	}
	ignoreCmd.Flags().BoolP("remove", "r", false, "Remove a command from the ignore list.")

	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run a saved and documented command referenced by alias.",
		Run:   run,
	}
	runCmd.Flags().BoolP("confirm", "c", false, "Run the command immediately without review.")

	rootCmd.AddCommand(
		//	sessionCmd,
		historyCmd,
		docCmd,
		ignoreCmd,
		runCmd,
	)

	cobra.OnInitialize(initConfig)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetTypeByDefaultValue(true)
	viper.SetDefault("history.maxlines", 1000)
	viper.SetDefault("history.reps-pre-doc-prompt", 3)

	viper.AddConfigPath(home + "/.coach")
	viper.SetConfigName("config")

	viper.SetEnvPrefix("coach")
	viper.AutomaticEnv()

	// if no config file is found, write the defaults to one
	if err := viper.ReadInConfig(); err != nil {
		defaults := viper.AllSettings()
		data, _ := yaml.Marshal(&defaults)
		ioutil.WriteFile(home+"/.coach/config.yaml", data, 0600)
	}
}

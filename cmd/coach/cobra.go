// Copyright (c) 2018, Adam Bright <brightam1@gmail.com>
// See LICENSE for licensing information

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	coach "github.com/alittlebrighter/coach-pro"
	"github.com/alittlebrighter/coach-pro/storage/database"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var (
	home, dbpath string
)

func main() {
	// Find home directory.
	home = coach.HomeDir()
	os.Mkdir(home, os.ModePerm)
	dbpath = home + "/coach.db"

	rootCmd := &cobra.Command{
		Use:   "coach",
		Short: "A tool to help you save and document common commands executed on the command line.",
		Long: fmt.Sprintf("%s\nAuthor: %s\n\nConfiguration: %s\nScript DB: %s",
			"Coach PRO: Save, document, query, and run all of your scripts.",
			"Adam Bright <brightam1@gmail.com>",
			home+"/config.yaml",
			home+"/coach.db",
		),
		Run: appMain,
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
	historyCmd.Flags().BoolP("all", "a", false, "Retrieve history from all sessions.")

	docCmd := &cobra.Command{
		Use:     "doc",
		Short:   "Save and query commands with tags and documentation.  Default is to document the most recent command.",
		Example: "coach doc [alias] [tags] [comment] # empty alias represented by \"\", tag list must be quoted if it contains spaces",
		Run:     doc,
	}
	docCmd.Flags().StringP("query", "q", database.Wildcard, "Query your saved commands by tags.")
	docCmd.Flags().StringP("script", "s", "", "Quoted command that you would like to document and save.")
	docCmd.Flags().StringP("edit", "e", "", "Edit the script specified by alias.")
	docCmd.Flags().IntP("history-lines", "l", 1, "Number of most recent lines in history to put into the script.")
	docCmd.Flags().String("delete", "", "Delete a saved script.")

	ignoreCmd := &cobra.Command{
		Use:   "ignore",
		Short: "Ignore this command when scanning for duplicates to prompt for documentation.  Defaults to the last run command.",
		Run:   ignore,
	}
	ignoreCmd.Flags().BoolP("remove", "r", false, "Remove a command from the ignore list.")
	ignoreCmd.Flags().BoolP("all", "a", false, "Ignore all non-compound commands starting with the first word from the previous line.")
	ignoreCmd.Flags().IntP("history-lines", "l", 1, "Number of most recent lines in history to ignore.")

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

	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			fmt.Printf("ERROR: Something went wrong.  Do you have permissions to access '%s' and its contents?\n", home)
			os.Exit(1)
		}
	}()

	timer := time.NewTimer(2 * time.Second)
	runErr := make(chan error)
	go func() {
		runErr <- rootCmd.Execute()
	}()
	select {
	case err := <-runErr:
		if err != nil {
			handleErr(err)
			os.Exit(1)
		}
	case <-timer.C:
		// timed out
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetTypeByDefaultValue(true)
	viper.SetDefault("history.maxlines", 1000)
	viper.SetDefault("history.reps-pre-doc-prompt", 3)
	viper.SetDefault("timestampFormat", "01/02 03:04:05PM")

	viper.AddConfigPath(home + "/.coach")
	viper.SetConfigName("config")

	viper.SetEnvPrefix("coach")
	viper.AutomaticEnv()

	// if no config file is found, write the defaults to one
	if err := viper.ReadInConfig(); err != nil {
		defaults := viper.AllSettings()
		data, _ := yaml.Marshal(&defaults)
		ioutil.WriteFile(home+"/config.yaml", data, 0600)
	}
}

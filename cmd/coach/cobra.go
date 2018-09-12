// Copyright (c) 2018, Adam Bright <brightam1@gmail.com>
// See LICENSE for licensing information

package main

import (
	"fmt"
	"io/ioutil"
	"os"

	coach "github.com/alittlebrighter/coach-pro"
	"github.com/alittlebrighter/coach-pro/storage/database"
	"github.com/alittlebrighter/coach-pro/trial"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var (
	home string
)

func main() {
	// Find home directory.
	home = coach.HomeDir()
	os.Mkdir(home, os.ModePerm)
	coach.DBPath = home + "/coach.db"

	rootCmd := &cobra.Command{
		Use:   "coach",
		Short: "A tool to help you save and document common commands executed on the command line.",
		Long: fmt.Sprintf("Coach %s: %s\n%s\nConfiguration: %s\nScript DB: %s\n\nFor support contact support.coach@mg.alittlebrighter.io",
			trial.Version,
			"Save, document, query, and run all of your scripts.",
			trial.ExpireNotice,
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
	historyCmd.Flags().Bool("import", false, "Import command history.  Supports bash and PowerShell.")

	docCmd := &cobra.Command{
		Use:     "lib",
		Short:   "Save and query scripts/commands with tags and documentation.  Default is to save and document the most recent command.",
		Example: "coach doc [alias] [tags] [comment] # empty alias represented by \"\", tag list must be quoted if it contains spaces",
		Run:     doc,
	}
	docCmd.Flags().StringP("query", "q", database.Wildcard, "Query your saved commands by tags.")
	docCmd.Flags().StringP("script", "s", "", "Quoted command that you would like to document and save.")
	docCmd.Flags().StringP("edit", "e", "", "Edit the script specified by alias.")
	docCmd.Flags().IntP("history-lines", "l", 1, "Number of most recent lines in history to put into the script.")
	docCmd.Flags().String("delete", "", "Delete a saved script.")
	docCmd.Flags().String("restore", "", "Restore a deleted script.")
	docCmd.Flags().Bool("empty-trash", false, "Completely erase all deleted scripts.")

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
	runCmd.Flags().BoolP("check", "c", false, "Review the command documentation before running.")

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

	if err := rootCmd.Execute(); err != nil {
		handleErr(err)
		os.Exit(1)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if _, err := os.Stat(coach.DBPath); err != nil {
		coach.GetStore(false).Close()
	}

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
		ioutil.WriteFile(home+"/config.yaml", data, database.FilePerms)
	}
}

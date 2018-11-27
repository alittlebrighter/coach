// Copyright (c) 2018, Adam Bright <brightam1@gmail.com>
// See LICENSE for licensing information

package main

import (
	"fmt"
	"os"

	coach "github.com/alittlebrighter/coach"
	conf "github.com/alittlebrighter/coach/config"
	"github.com/alittlebrighter/coach/platforms"
	"github.com/alittlebrighter/coach/storage/database"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	home string
)

func main() {
	// Find home directory.
	home = conf.HomeDir()
	os.Mkdir(home, os.ModePerm)
	coach.DBPath = home + "/coach.db"

	rootCmd := &cobra.Command{
		Use:   "coach",
		Short: "A tool to help you save and document common commands executed on the command line.",
		Long: fmt.Sprintf("Coach: %s\n\nConfiguration: %s\nScript DB: %s",
			"Save, document, query, and run all of your scripts.",
			home+"/config",
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
		Short: "Store and query command history",
		Run:   history,
	}
	historyCmd.Flags().StringP("record", "r", "", "Record command.  With `bash` you can use \"$(history 1)\"")
	historyCmd.Flags().StringP("query", "q", "", "Query command history by regex.")
	historyCmd.Flags().BoolP("all", "a", false, "Retrieve history from all sessions")

	docCmd := &cobra.Command{
		Use:     "lib",
		Short:   "Save and query scripts/commands with tags and documentation",
		Example: "coach doc [alias] [tags] [comment] # empty alias represented by \"\", tag list must be quoted if it contains spaces",
		Run:     doc,
	}
	docCmd.Flags().StringP("query", "q", database.Wildcard, "Query your saved commands by tags")
	docCmd.Flags().StringP("script", "s", "", "Quoted command that you would like to document and save")
	docCmd.Flags().StringP("edit", "e", "", "Edit the script specified by alias")
	docCmd.Flags().IntP("history-lines", "l", 1, "Number of most recent lines in history to put into the script")
	docCmd.Flags().String("delete", "", "Delete a saved script")
	docCmd.Flags().String("restore", "", "Restore a deleted script")
	docCmd.Flags().Bool("empty-trash", false, "Completely erase all deleted scripts")

	ignoreCmd := &cobra.Command{
		Use:   "ignore",
		Short: "Ignore this command when scanning for duplicates to prompt for documentation.  Defaults to the last run command",
		Run:   ignore,
	}
	ignoreCmd.Flags().BoolP("remove", "r", false, "Remove a command from the ignore list")
	ignoreCmd.Flags().BoolP("all", "a", false, "Ignore all non-compound commands starting with the first word from the previous line")
	ignoreCmd.Flags().IntP("history-lines", "l", 1, "Number of most recent lines in history to ignore")

	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run a saved and documented command referenced by alias",
		Run:   run,
	}
	runCmd.Flags().BoolP("check", "c", false, "Review the command documentation before running")
	runCmd.Flags().DurationP("timeout", "t", 0, "Specify a maximum time this script should run")

	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Set default config values",
		Run:   config,
	}

	rootCmd.AddCommand(
		//	sessionCmd,
		historyCmd,
		docCmd,
		ignoreCmd,
		runCmd,
		configCmd,
	)

	cobra.OnInitialize(initConfig)

	/*
		defer func() {
			if r := recover(); r != nil {
				fmt.Println(r)
				fmt.Printf("ERROR: Something went wrong.  Do you have permissions to access '%s' and its contents?\n", home)
				os.Exit(1)
			}
		}()
	*/

	if err := rootCmd.Execute(); err != nil {
		handleErr(err)
		os.Exit(1)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if _, err := os.Stat(coach.DBPath); err != nil {
		store := coach.GetStore(false)
		store.Init()
		store.Close()
	}

	conf.AppConfiguration()

	viper.SetDefault("default_shell", platforms.DefaultShell)
	viper.SetDefault("history.max_lines", 1000)
	viper.SetDefault("history.reps_pre_doc_prompt", 3)
	viper.SetDefault("timestamp_format", "01/02 03:04:05PM")

	viper.ReadInConfig()

	conf.WriteConfiguration()
}

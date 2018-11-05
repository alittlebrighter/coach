package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/alittlebrighter/coach/config"
	"github.com/alittlebrighter/coach/storage/database"
)

var buildTimestamp string

var home, dbpath string

func main() {
	// Find home directory.
	home = config.HomeDir()
	os.Mkdir(home, os.ModePerm)
	dbpath = home + "/coach.db"

	importCmd := &cobra.Command{
		Use:   "coach-trader",
		Short: "Use this to import existing scripts into `coach`.",
		Long: fmt.Sprintf("%s\nAuthor: %s\n\nScript DB: %s",
			"coach-trader - import and export your existing scripts",
			"Adam Bright <brightam1@gmail.com>",
			home+"/coach.db",
		),
		Example: "coach-trader [dir] # import all scripts contained in dir\n" +
			"coach-trader --export [dir] # export all scripts saved in coach to dir",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				handleErr(errors.New("must specify directory for import/export"))
				os.Exit(1)
			}

			export, _ := cmd.Flags().GetBool("export")
			store, err := database.NewBoltDB(dbpath, export)
			if err != nil {
				handleErr(err)
				return
			}
			defer store.Close()

			if export {
				exportScripts(args[0], store)
			} else {
				importScripts(args[0], store)
			}
		},
	}
	importCmd.Flags().BoolP("export", "e", false, "Set this flag to export scripts from `coach`.")

	if err := importCmd.Execute(); err != nil {
		handleErr(err)
		os.Exit(1)
	}
}

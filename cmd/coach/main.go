package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/rs/xid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/alittlebrighter/coach"
	"github.com/alittlebrighter/coach/gen/models"
	"github.com/alittlebrighter/coach/storage/database"
)

func appMain(cmd *cobra.Command, args []string) {
	cmd.Help()
}

func session(cmd *cobra.Command, args []string) {

}

func history() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		record, rErr := cmd.Flags().GetBool("record")

		switch {
		case rErr == nil && record:
			store, err := database.NewBoltDB(dbpath, false)
			if err != nil {
				handleErr(err)
				return
			}
			defer store.Close()

			var historyOutput string
			if args != nil || len(args) > 0 {
				historyOutput = args[0]
			}

			dupeCount := viper.GetInt("history.reps-pre-doc-prompt")

			if enoughDupes, _ := coach.SaveHistory(historyOutput, dupeCount, store); enoughDupes {
				fmt.Printf("This command has been used %d+ times.\nRun `coach doc [alias] [tags] [comment]` to document this command.\n",
					dupeCount)
			}
		default:
			store, err := database.NewBoltDB(dbpath, true)
			if err != nil {
				handleErr(err)
				return
			}
			defer store.Close()

			count := 10
			if args != nil && len(args) >= 1 {
				if count, err = strconv.Atoi(args[0]); err != nil {
					count = 10
				}
			}
			lines, err := coach.GetRecentHistory(count, store)
			if err != nil {
				fmt.Println("Could not retrieve history for this session!  ERROR:", err)
				break
			}

			for _, line := range lines {
				id, _ := xid.FromBytes(line.GetId())
				fmt.Printf("%s %s - %s\n", id.Time().Format(time.RFC3339), id.String(),
					line.GetFullCommand())
			}
		}
	}
}

func doc(cmd *cobra.Command, args []string) {
	query, qErr := cmd.Flags().GetString("query")
	script, cErr := cmd.Flags().GetString("cmd")

	switch {
	case qErr == nil && len(query) > 0:
		store, err := database.NewBoltDB(dbpath, true)
		if err != nil {
			handleErr(err)
			return
		}
		defer store.Close()

		cmds, err := coach.QueryScripts(query, store)
		if err != nil {
			handleErr(err)
			return
		}

		for _, sCmd := range cmds {
			if sCmd.GetId() == nil || len(sCmd.GetId()) == 0 {
				continue
			}

			fmt.Printf("%14s: %s\n%14s: %s\n%14s: %s\n%14s: %s\n---\n",
				"Script", sCmd.GetScript().GetContent(),
				"Alias", sCmd.GetAlias(),
				"Tags", strings.Join(sCmd.GetTags(), ","),
				"Documentation", sCmd.GetDocumentation(),
			)
		}

	case len(args) >= 3:
		store, err := database.NewBoltDB(dbpath, false)
		if err != nil {
			handleErr(err)
			return
		}
		defer store.Close()

		if cErr != nil || len(script) == 0 {
			if hLines, err := coach.GetRecentHistory(1, store); err == nil && len(hLines) > 0 {
				script = hLines[0].GetFullCommand()
			}
		}

		err = coach.SaveScript(args[0], strings.Split(args[1], ","), strings.Join(args[2:], " "), script, store)
		if err != nil {
			handleErr(err)
			return
		}
	}
}

// should run this through "business logic" in coach package
func ignore(cmd *cobra.Command, args []string) {
	store, err := database.NewBoltDB(dbpath, false)
	if err != nil {
		handleErr(err)
		return
	}
	defer store.Close()

	var fullCmd string
	if hLines, err := coach.GetRecentHistory(1, store); err == nil && len(hLines) > 0 {
		fullCmd = hLines[0].GetFullCommand()
	}

	remove, rErr := cmd.Flags().GetBool("remove")
	if rErr == nil && remove {
		store.UnignoreCommand(fullCmd)
	} else {
		store.IgnoreCommand(fullCmd)
	}
}

func run(cmd *cobra.Command, args []string) {
	if args == nil || len(args) == 0 {
		fmt.Println("No alias specified.")
	}
	store, err := database.NewBoltDB(dbpath, false)
	if err != nil {
		handleErr(err)
		return
	}
	defer store.Close()

	toRun := store.GetScript(args[0])
	if toRun == nil {
		fmt.Println("No script found by that alias.")
		return
	}

	fmt.Printf("Command '%s' found:\n# %s\n%s\n", toRun.GetAlias(), toRun.GetDocumentation(), toRun.GetScript().GetContent())

	if confirmed, cErr := cmd.Flags().GetBool("confirm"); cErr == nil && confirmed {
		fmt.Println("Running now...")
	} else {
		fmt.Print("Run now? [y/n] ")
		in, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil || (len(in) >= 1 && in[0] != byte('y')) {
			fmt.Println("Not running command.")
			return
		}
	}

	if err := coach.RunScript(*toRun); err != nil {
		handleErr(err)
	}

}

type DocStore interface {
	SaveDoc(*models.DocumentedScript) error
	QueryDoc(...string) ([]models.DocumentedScript, error)
}

func GetTTY() string {
	ttyCmd := exec.Command("tty")
	ttyCmd.Stdin = os.Stdin
	ttyBytes, _ := ttyCmd.Output()
	return strings.TrimSpace(string(ttyBytes))
}

func handleErr(e error) {
	if e != nil {
		fmt.Println("ERROR:", e)
	}
}

func readTags(val string) []string {
	stringReader := strings.NewReader(val)
	csvReader := csv.NewReader(stringReader)
	tags, _ := csvReader.Read()
	return tags
}

// Copyright (c) 2018, Adam Bright <brightam1@gmail.com>
// See LICENSE for licensing information

package main

import (
	"bufio"
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
	fmt.Println("\nAuthor: Adam Bright <brightam1@gmail.com>")
}

func session(cmd *cobra.Command, args []string) {

}

func history(cmd *cobra.Command, args []string) {
	record, rErr := cmd.Flags().GetString("record")

	switch {
	case rErr == nil && len(record) > 0:
		store, err := database.NewBoltDB(dbpath, false)
		if err != nil {
			handleErr(err)
			return
		}
		defer store.Close()

		dupeCount := viper.GetInt("history.reps-pre-doc-prompt")

		if enoughDupes, _ := coach.SaveHistory(record, dupeCount, store); enoughDupes {
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

func doc(cmd *cobra.Command, args []string) {
	query, qErr := cmd.Flags().GetString("query")
	script, cErr := cmd.Flags().GetString("cmd")
	edit, eErr := cmd.Flags().GetString("edit")
	hLines, _ := cmd.Flags().GetInt("history")

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

			scriptStr := sCmd.GetScript().GetContent()
			if len(scriptStr) > 21 {
				scriptStr = scriptStr[:21] + "..."
			}
			fmt.Printf("%14s: %s\n%14s: %s\n%14s: %s\n%14s: %s\n---\n",
				"Script", scriptStr,
				"Alias", sCmd.GetAlias(),
				"Tags", strings.Join(sCmd.GetTags(), ","),
				"Documentation", sCmd.GetDocumentation(),
			)
		}
	case eErr == nil && len(edit) > 0:
		store, err := database.NewBoltDB(dbpath, false)
		if err != nil {
			handleErr(err)
			return
		}
		defer store.Close()

		if err := coach.EditScript(edit, store); err != nil {
			handleErr(err)
		}
		return
	case len(args) >= 3:
		store, err := database.NewBoltDB(dbpath, false)
		if err != nil {
			handleErr(err)
			return
		}
		defer store.Close()

		if cErr != nil || len(script) == 0 {
			if lines, err := coach.GetRecentHistory(hLines, store); err == nil && len(lines) > 0 {
				for _, line := range lines {
					script += line.GetFullCommand() + "\n"
				}
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

	scriptStr := toRun.GetScript().GetContent()
	if len(scriptStr) > 21 {
		scriptStr = scriptStr[:21] + "..."
	}
	fmt.Printf("Command '%s' found:\n# %s\n%s\n", toRun.GetAlias(), toRun.GetDocumentation(), scriptStr)

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

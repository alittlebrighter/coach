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
	"github.com/alittlebrighter/coach/platforms"
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
			fmt.Printf("\n---\nThis command has been used %d+ times.\nRun `coach doc [alias] "+
				"[tags] [comment...]` to save and document this command.\nRun `coach ignore` to silence "+
				"this output for this command.\n",
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
	script, cErr := cmd.Flags().GetString("script")
	edit, eErr := cmd.Flags().GetString("edit")
	hLines, _ := cmd.Flags().GetInt("history")
	delete, _ := cmd.Flags().GetString("delete")

	switch {
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

		shell := platforms.IdentifyShell()
		if len(shell) == 0 {
			fmt.Println("You're shell could not be identified.  Using 'bash' for now.\nRun `coach doc -e " + args[0] + "` to edit.")
			shell = "bash"
		}
		err = coach.SaveScript(models.DocumentedScript{
			Alias:         args[0],
			Tags:          strings.Split(args[1], ","),
			Documentation: strings.Join(args[2:], " "),
			Script:        &models.Script{Content: script, Shell: shell}},
			false, store)
		if err != nil {
			handleErr(err)
			return
		}
	case eErr == nil && len(edit) > 0:
		store, err := database.NewBoltDB(dbpath, false)
		if err != nil {
			handleErr(err)
			return
		}
		defer store.Close()

		newScript, err := coach.EditScript(edit, store)
		if err != nil {
			handleErr(err)
			return
		}

		overwrite := true
		if newScript.GetAlias() != edit {
			overwrite = false
		}

		save := func(ovrwrt bool) error {
			err := coach.SaveScript(*newScript, ovrwrt, store)

			if newScript.GetAlias() != edit && err == nil {
				store.DeleteScript([]byte(edit))
			}
			return err
		}

		for err = save(overwrite); err == database.ErrAlreadyExists; err = save(overwrite) {
			fmt.Printf("The alias '%s' already exists.\n", newScript.GetAlias())
			fmt.Printf("Enter '%s' again to overwrite, or try something else: ", newScript.GetAlias())
			in, inErr := bufio.NewReader(os.Stdin).ReadString('\n')
			if inErr != nil || len(strings.TrimSpace(in)) == 0 {
				overwrite = false
				continue
			}
			input := strings.Fields(in)[0]

			if input == newScript.GetAlias() || input == edit {
				overwrite = true
				continue
			} else {
				newScript.Alias = input
				overwrite = false
			}
		}

		if err != nil {
			handleErr(err)
			return
		}
	case len(delete) > 0:
		store, err := database.NewBoltDB(dbpath, false)
		if err != nil {
			handleErr(err)
			return
		}
		defer store.Close()

		if store.GetScript([]byte(strings.TrimSpace(delete))) == nil {
			handleErr(database.ErrNotFound)
			return
		}

		fmt.Printf("Type '%s' again to delete: ", delete)
		in, err := bufio.NewReader(os.Stdin).ReadString('\n')
		input := strings.Fields(in)
		if err != nil || len(input) == 0 {
			fmt.Println("Not deleting.")
			return
		}

		fmt.Printf("Deleting '%s' now.\n", input[0])
		if err := store.DeleteScript([]byte(input[0])); err != nil {
			handleErr(err)
			return
		}
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

			fmt.Printf("%14s: %s\n%14s: %s\n%14s: %s\n%14s: %s\n%17s\n",
				"Script", Slugify(sCmd.GetScript().GetContent(), 48),
				"Alias", sCmd.GetAlias(),
				"Tags", strings.Join(sCmd.GetTags(), ","),
				"Documentation", sCmd.GetDocumentation(),
				"---",
			)
		}
	}
	return
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

	toRun := store.GetScript([]byte(args[0]))
	if toRun == nil {
		fmt.Println("No script found by that alias.")
		return
	}

	if confirmed, cErr := cmd.Flags().GetBool("confirm"); cErr == nil && confirmed {
		fmt.Println("Running '" + toRun.GetAlias() + "'...")
	} else {
		fmt.Printf("Command '%s' found:\n###\n%s\n###\n$ %s\n\n", toRun.GetAlias(), toRun.GetDocumentation(), Slugify(toRun.GetScript().GetContent(), 48))
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

func Slugify(content string, length uint) string {
	lines := strings.Split(content, "\n")
	scriptStr := strings.TrimSpace(lines[0])
	if len(scriptStr) > int(length) {
		scriptStr = scriptStr[:length] + "..."
	} else if len(lines) > 1 {
		scriptStr += "..."
	}
	return scriptStr
}

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

	"github.com/golang/protobuf/ptypes"
	"github.com/rs/xid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/alittlebrighter/coach"
	"github.com/alittlebrighter/coach/gen/models"
	"github.com/alittlebrighter/coach/storage/database"
)

func appMain(cmd *cobra.Command, args []string) {

}

func session(cmd *cobra.Command, args []string) {

}

func history() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		record, rErr := cmd.Flags().GetBool("record")

		tty := GetTTY()

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
			hLine := models.HistoryRecord{
				Timestamp:   ptypes.TimestampNow(),
				FullCommand: coach.CleanupCommand(historyOutput),
				Tty:         tty,
			}
			if len(hLine.GetFullCommand()) == 0 || strings.HasPrefix(hLine.GetFullCommand(), "coach") {
				return
			}

			enoughDupes := make(chan bool)
			dupeCount := viper.GetInt("history.reps-pre-doc-prompt")
			go func(line *models.HistoryRecord, count int, reporter chan bool) {
				reporter <- store.CheckDupeCmds(line.GetFullCommand(), dupeCount)
			}(&hLine, dupeCount, enoughDupes)
			go store.SaveHistory(&hLine)

			if <-enoughDupes {
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
			lines, err := store.GetRecent(tty, count)
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

type HistoryStore interface {
	SaveHistory(*models.HistoryRecord) error
	GetRecent(string, int) ([]models.HistoryRecord, error)
	Close()
}

func doc(cmd *cobra.Command, args []string) {
	query, qErr := cmd.Flags().GetString("query")
	command, cErr := cmd.Flags().GetString("cmd")

	switch {
	case qErr == nil && len(query) > 0:
		store, err := database.NewBoltDB(dbpath, true)
		if err != nil {
			handleErr(err)
			return
		}
		defer store.Close()

		cmds, err := store.QueryDoc(readTags(query)...)
		if err != nil {
			handleErr(err)
			return
		}

		for _, sCmd := range cmds {
			if sCmd.GetId() == nil || len(sCmd.GetId()) == 0 {
				continue
			}

			cmdString := ""
			for _, c := range sCmd.GetCommand() {
				cmdString += c.GetCommand() + " " + strings.Join(c.GetArguments(), " ") + "\n"
			}
			fmt.Printf("%-8s: %s%-8s: %s\n%-8s: %s\n%-8s: %s\n---\n",
				"Command", cmdString,
				"Alias", sCmd.GetAlias(),
				"Tags", strings.Join(sCmd.GetTags(), ","),
				"Comment", sCmd.GetComment(),
			)
		}

	case len(args) >= 3:
		savedCmd := models.SavedCommand{
			Alias:   args[0],
			Tags:    strings.Split(args[1], ","),
			Comment: strings.Join(args[2:], " "),
			Command: []*models.Command{},
		}

		if len(savedCmd.GetAlias()) > 0 {
			savedCmd.Id = []byte(savedCmd.GetAlias())
		}

		store, err := database.NewBoltDB(dbpath, false)
		if err != nil {
			handleErr(err)
			return
		}
		defer store.Close()

		if cErr == nil && len(command) > 0 {
			// HACK
			cmdIn, err := coach.ParseHistory("1 " + command)
			if err == nil {
				savedCmd.Command = append(savedCmd.Command, cmdIn)
			}
		} else if hLine, err := store.GetRecent(GetTTY(), 1); err == nil && len(hLine) > 0 {
			// HACK
			hCommand, _ := coach.ParseHistory("1 " + hLine[0].GetFullCommand())
			savedCmd.Command = append(savedCmd.Command, hCommand)
		}

		if err := store.SaveDoc(&savedCmd); err != nil {
			handleErr(err)
			return
		}
	}
}

func ignore(cmd *cobra.Command, args []string) {
	store, err := database.NewBoltDB(dbpath, false)
	if err != nil {
		handleErr(err)
		return
	}
	defer store.Close()

	tty := GetTTY()

	var fullCmd string
	if hLines, err := store.GetRecent(tty, 1); err == nil && len(hLines) > 0 {
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

	toRun := store.GetSavedCmd(args[0])
	if toRun == nil {
		fmt.Println("No command found by that alias.")
		return
	}

	formattedCommand := ""
	for _, cmdLine := range toRun.Command {
		formattedCommand += strings.Join(append([]string{cmdLine.GetCommand()}, cmdLine.GetArguments()...), " ") + "\n"
	}

	fmt.Printf("Command '%s' found:\n# %s\n%s", toRun.GetAlias(), toRun.GetComment(), formattedCommand)

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

	for _, run := range toRun.GetCommand() {
		runCmd := exec.Command(run.GetCommand(), run.GetArguments()...)
		runCmd.Stdin = os.Stdin
		runCmd.Stdout = os.Stdout
		runCmd.Stderr = os.Stderr
		runCmd.Run()
	}
}

type DocStore interface {
	SaveDoc(*models.SavedCommand) error
	QueryDoc(...string) ([]models.SavedCommand, error)
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

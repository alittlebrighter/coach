package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/rs/xid"
	"github.com/spf13/cobra"

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
		//ignore, iErr := cmd.Flags().GetString("ignore")

		// this needs to be injected
		ttyCmd := exec.Command("tty")
		ttyCmd.Stdin = os.Stdin
		ttyBytes, _ := ttyCmd.Output()
		tty := strings.TrimSpace(string(ttyBytes))

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
			hLine := models.HistoryRecord{Timestamp: ptypes.TimestampNow()}
			hLine.Command, err = coach.ParseHistory(historyOutput)
			if err != nil {
				handleErr(err)
				break
			}
			hLine.Tty = tty
			store.Save(&hLine)
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
				fmt.Printf("(%s) %s | %s\n", id.String(), id.Time().Format(time.RFC3339),
					strings.Join(append([]string{line.Command.Command}, line.Command.Arguments...), " "))
			}
		}
	}
}

type HistoryStore interface {
	Save(*models.HistoryRecord) error
	GetRecent(string, int) ([]models.HistoryRecord, error)
	Close()
}

func doc(cmd *cobra.Command, args []string) {

}

func handleErr(e error) {
	if e != nil {
		fmt.Println("ERROR:", e)
	}
}

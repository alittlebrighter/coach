package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/alittlebrighter/coach"
)

func appMain(cmd *cobra.Command, args []string) {
	fmt.Println("coach called")
}

func observeRun(cmd *cobra.Command, args []string) {
	var historyOutput string
	if args != nil || len(args) > 0 {
		historyOutput = args[0]
	}
	command, err := coach.ParseHistory(historyOutput)
	fmt.Println("command:", command, "err:", err)
}

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	coach "github.com/alittlebrighter/coach-pro"
	"github.com/alittlebrighter/coach-pro/platforms"
	"github.com/alittlebrighter/coach-pro/storage/database"
)

const header = "exported from COACH (https://github.com/alittlebrighter/coach)"

func importScripts(dir string, store *database.BoltDB) {

}

func exportScripts(dir string, store *database.BoltDB) {
	scripts, err := coach.QueryScripts(database.Wildcard, store)
	if err != nil {
		handleErr(err)
		return
	}

	for _, script := range scripts {
		shell := platforms.GetShell(script.GetScript().GetShell())
		shebang := "#!/usr/bin/env " + script.GetScript().GetShell()

		path := filepath.Join(append([]string{dir}, strings.Split(script.GetAlias(), ".")...)...)
		fullPath := ""
		for _, subdir := range strings.Split(path[:strings.LastIndex(path, "/")], "/") {
			fullPath += subdir + "/"
			os.Mkdir(fullPath, 0700)
		}

		file, err := os.OpenFile(path+"."+shell.FileExtension(), os.O_CREATE|os.O_RDWR, os.ModePerm)
		handleErr(err)

		file.WriteString(shebang + platforms.Newline(2))
		file.WriteString(shell.LineComment() + " " + header + platforms.Newline(2))
		file.Write(coach.MarshalEdit(script))
		file.Close()
	}
}

func handleErr(e error) {
	if e != nil {
		fmt.Println("ERROR:", e)
	}
}

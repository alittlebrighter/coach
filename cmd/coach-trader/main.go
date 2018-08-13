package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	coach "github.com/alittlebrighter/coach-pro"
	"github.com/alittlebrighter/coach-pro/gen/models"
	"github.com/alittlebrighter/coach-pro/platforms"
	"github.com/alittlebrighter/coach-pro/storage/database"
	"github.com/rs/xid"
)

func importScripts(dir string, store *database.BoltDB) {
	_, err := os.Stat(dir)
	if err != nil {
		handleErr(err)
		return
	}

	base := dir
	lastSlash := strings.LastIndex(dir, "/")
	if lastSlash >= 0 {
		base = dir[lastSlash+1:]
	}

	filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		script, err := ParseFile(p, base, info.Size())
		if err != nil {
			handleErr(err)
			return err
		}
		for err = coach.SaveScript(*script, false, store); err == database.ErrAlreadyExists; err = coach.SaveScript(*script, false, store) {
			script.Alias += "-" + xid.New().String()
		}
		if err != nil {
			handleErr(err)
			return err
		}

		return nil
	})
}

func ParseFile(path, base string, size int64) (*models.DocumentedScript, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	script := new(models.DocumentedScript)

	var ext string
	parts := strings.Split(path, ".")
	if len(parts) < 2 {
		ext = database.Wildcard
	} else {
		ext = parts[1]
	}
	fromBase := parts[0][strings.Index(parts[0], base):]
	pathParts := strings.Split(fromBase, "/")
	script.Alias = strings.Join(pathParts, ".")
	script.Tags = pathParts
	script.Script = &models.Script{Shell: platforms.ShellNameFromExt(ext)}

	data := make([]byte, size)
	_, err = f.Read(data)
	contents := string(data)
	if len(strings.TrimSpace(contents)) == 0 {
		return nil, errors.New("no content")
	}

	shell := platforms.GetShell(script.GetScript().GetShell())
	for _, line := range strings.Split(contents, platforms.Newline(1)) {
		trimmed := strings.TrimSpace(line)

		switch {
		case strings.HasPrefix(trimmed, "#!"):
			parts := strings.Split(trimmed, " ")
			interpreter := parts[0][strings.LastIndex(parts[0], "/")+1:]
			if interpreter == "env" {
				script.Script.Shell = parts[1]
				shell = platforms.GetShell(parts[1])
			} else {
				script.Script.Shell = interpreter
				shell = platforms.GetShell(interpreter)
			}
		case coach.UnmarshalLine(trimmed, script):
			// no-op
		case strings.HasPrefix(trimmed, shell.LineComment()):
			docLine := strings.TrimLeft(line, " \t"+shell.LineComment()) + platforms.Newline(1)
			if len(strings.TrimSpace(docLine)) > 0 {
				script.Documentation += docLine
			}
			fallthrough
		default:
			script.Script.Content += line + platforms.Newline(1)
		}
	}
	script.Documentation = strings.TrimSpace(script.GetDocumentation())

	return script, nil
}

func exportScripts(dir string, store *database.BoltDB) {
	scripts, err := coach.QueryScripts(database.Wildcard, store)
	if err != nil {
		handleErr(err)
		return
	}

	for _, script := range scripts {
		shell := platforms.GetShell(script.GetScript().GetShell())

		var shebang string
		if !strings.Contains(strings.ToLower(script.GetScript().GetShell()), "powershell") ||
			strings.ToLower(script.GetScript().GetShell()) != "windowscmd" {
			shebang = "#!/usr/bin/env " + script.GetScript().GetShell()
		}

		path := filepath.Join(append([]string{dir}, strings.Split(script.GetAlias(), ".")...)...)
		fullPath := ""
		lastSlash := strings.LastIndex(path, "/")
		allButLastDir := path
		if lastSlash != -1 {
			allButLastDir = path[:lastSlash]
		}
		for _, subdir := range strings.Split(allButLastDir, "/") {
			fullPath += subdir + "/"
			os.Mkdir(fullPath, 0700)
		}

		file, err := os.OpenFile(path+"."+shell.FileExtension(), os.O_CREATE|os.O_RDWR, os.ModePerm)
		handleErr(err)

		file.WriteString(shebang + platforms.Newline(2))
		file.WriteString(shell.LineComment() + " " + coach.Header + platforms.Newline(2))
		file.Write(coach.MarshalEdit(script))
		file.Close()
	}
}

func handleErr(e error) {
	if e != nil {
		fmt.Println("ERROR:", e)
	}
}

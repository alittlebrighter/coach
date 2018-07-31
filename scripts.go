// Copyright (c) 2018, Adam Bright <brightam1@gmail.com>
// See LICENSE for licensing information

package coach

import (
	"errors"
	"os"
	"strings"

	"github.com/alittlebrighter/coach/gen/models"
	"github.com/alittlebrighter/coach/platforms"
)

func QueryScripts(query string, store ScriptStore) (scripts []models.DocumentedScript, err error) {
	if len(query) == 0 {
		err = errors.New("no query specified")
		return
	}

	scripts, err = store.QueryScripts(strings.Split(strings.TrimSpace(query), ",")...)
	return
}

func SaveScript(alias string, tags []string, documentation string, script string, store ScriptStore) (err error) {
	if len(script) == 0 {
		return errors.New("no script to save")
	}

	toSave := models.DocumentedScript{
		Alias:         alias,
		Tags:          tags,
		Documentation: documentation,

		// TODO: parse variables out of script
		Script: &models.Script{Content: platforms.CleanupCommand(script)},
	}

	if len(toSave.GetAlias()) > 0 {
		toSave.Id = []byte(toSave.GetAlias())
	} else {
		toSave.Id = RandomID()
		toSave.Alias = string(toSave.GetId())
	}

	if err = store.Save(toSave.GetId(), toSave); err != nil {
		return
	}

	// save to ignore list
	err = store.IgnoreCommand(toSave.GetScript().GetContent())
	return
}

func RunScript(script models.DocumentedScript) error {
	toRun := Shell.BuildCommand(script.GetScript().GetContent())
	toRun.Stdin = os.Stdin
	toRun.Stdout = os.Stdout
	toRun.Stderr = os.Stderr
	toRun.Run()

	return nil
}

type ScriptStore interface {
	Save(id []byte, value interface{}) error
	QueryScripts(...string) ([]models.DocumentedScript, error)
	IgnoreCommand(string) error
}

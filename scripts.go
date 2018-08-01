// Copyright (c) 2018, Adam Bright <brightam1@gmail.com>
// See LICENSE for licensing information

package coach

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/alittlebrighter/coach/gen/models"
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
		Script: &models.Script{Content: script},
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

func EditScript(alias string, store ScriptStore) error {
	script := store.GetScript(alias)
	if script == nil {
		return errors.New("notfound")
	}

	tmpfile, ferr := ioutil.TempFile("", "coach")
	if ferr != nil {
		return ferr
	}
	var err error
	defer func() {
		if err == nil {
			os.Remove(tmpfile.Name()) // clean up
		} else {
			fmt.Println("There was an issue saving your edited script to the database.\nYou can find your draft here:", tmpfile.Name())
		}
	}()

	if _, err := tmpfile.Write([]byte(script.GetScript().GetContent())); err != nil {
		return err
	}
	if err := tmpfile.Close(); err != nil {
		return err
	}

	if err := Shell.OpenEditor(tmpfile.Name()); err != nil {
		return err
	}
	newContents, err := ioutil.ReadFile(tmpfile.Name())
	if err != nil {
		return err
	}

	script.Script.Content = string(newContents)
	if err := store.Save(script.GetId(), *script); err != nil {
		return err
	}
	return nil
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
	GetScript(alias string) *models.DocumentedScript
	QueryScripts(...string) ([]models.DocumentedScript, error)
	IgnoreCommand(string) error
}

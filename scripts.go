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

func SaveScript(toSave models.DocumentedScript, overwrite bool, store ScriptStore) (err error) {
	if toSave.GetScript() == nil || len(toSave.GetScript().GetContent()) == 0 {
		return errors.New("no script to save")
	}

	toSave.Alias = strings.TrimSpace(toSave.GetAlias())
	for i := range toSave.GetTags() {
		toSave.Tags[i] = strings.TrimSpace(toSave.GetTags()[i])
	}
	toSave.Documentation = strings.TrimSpace(toSave.GetDocumentation())
	// TODO: parse variables out of script
	toSave.Script.Content = strings.TrimSpace(toSave.GetScript().GetContent()) + "\n"
	toSave.Script.Shell = strings.TrimSpace(toSave.GetScript().GetShell())

	if len(toSave.GetAlias()) == 0 {
		toSave.Alias = string(RandomID())
	}
	toSave.Id = []byte(toSave.GetAlias())

	if err = store.Save(toSave.GetId(), toSave, overwrite); err != nil {
		return
	}

	// save to ignore list
	err = store.IgnoreCommand(toSave.GetScript().GetContent())
	return
}

func EditScript(alias string, store ScriptStore) (*models.DocumentedScript, error) {
	script := store.GetScript([]byte(alias))
	if script == nil {
		script = &models.DocumentedScript{
			Id:     []byte(alias),
			Alias:  alias,
			Tags:   []string{alias},
			Script: &models.Script{Shell: platforms.IdentifyShell()},
		}
	}

	tmpfile, ferr := ioutil.TempFile("", "coach")
	if ferr != nil {
		return nil, ferr
	}
	var err error
	defer func() {
		if err == nil {
			os.Remove(tmpfile.Name()) // clean up
		} else {
			fmt.Println("There was an issue saving your edited script to the database.\nYou can find your draft here:", tmpfile.Name())
		}
	}()

	if _, err := tmpfile.Write(MarshalEdit(*script)); err != nil {
		return nil, err
	}
	if err := tmpfile.Close(); err != nil {
		return nil, err
	}

	if err := Shell.OpenEditor(tmpfile.Name()); err != nil {
		return nil, err
	}
	newContents, err := ioutil.ReadFile(tmpfile.Name())
	if err != nil {
		return nil, err
	}

	var newScript models.DocumentedScript
	if newScript, err = UnmarshalEdit(string(newContents)); err != nil {
		return nil, err
	}

	return &newScript, nil
}

func RunScript(script models.DocumentedScript) error {
	shell := platforms.GetShell(script.GetScript().GetShell())

	toRun, cleanup, err := shell.BuildCommand(script.GetScript().GetContent())
	if cleanup != nil {
		defer cleanup()
	}
	if err != nil {
		return err
	}
	toRun.Stdin = os.Stdin
	toRun.Stdout = os.Stdout
	toRun.Stderr = os.Stderr
	toRun.Run()

	return nil
}

type ScriptStore interface {
	Save(id []byte, value interface{}, overwrite bool) error
	GetScript(id []byte) *models.DocumentedScript
	QueryScripts(...string) ([]models.DocumentedScript, error)
	DeleteScript(id []byte) error
	IgnoreCommand(string) error
}

const doNotEditLine = "!DO NOT EDIT THIS LINE!"

func MarshalEdit(s models.DocumentedScript) []byte {
	var contents strings.Builder
	contents.WriteString("-ALIAS- = " + s.GetAlias() + "\n")
	contents.WriteString("-TAGS- = " + strings.Join(s.GetTags(), ",") + "\n")
	contents.WriteString("-SHELL- = " + s.GetScript().GetShell() + "\n\n")
	contents.WriteString("-DOCUMENTATION- " + doNotEditLine + "\n")
	contents.WriteString(s.GetDocumentation() + "\n\n")
	contents.WriteString("-SCRIPT- " + doNotEditLine + "\n")
	contents.WriteString(s.GetScript().GetContent())
	return []byte(contents.String())
}

func UnmarshalEdit(contents string) (ds models.DocumentedScript, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("could not parse file")
		}
	}()

	ds.Script = new(models.Script)

	parts := strings.Split(contents, "\n")
	var inDoc, docStarted, inScript, scriptStarted bool
	for _, p := range parts {
		part := strings.TrimSpace(p)
		if !docStarted && !scriptStarted && len(part) == 0 {
			continue
		}

		switch {
		case strings.HasPrefix(part, "-ALIAS- ="):
			ds.Alias = strings.TrimSpace(strings.Split(part, "=")[1])

			inDoc, docStarted = false, false
			inScript, scriptStarted = false, false
		case strings.HasPrefix(part, "-TAGS- ="):
			tags := strings.Split(strings.Split(part, "=")[1], ",")
			for i := range tags {
				tags[i] = strings.TrimSpace(tags[i])
			}
			ds.Tags = tags

			inDoc, docStarted = false, false
			inScript, scriptStarted = false, false
		case strings.HasPrefix(part, "-SHELL- ="):
			ds.Script.Shell = strings.TrimSpace(strings.Split(part, "=")[1])

			inDoc, docStarted = false, false
			inScript, scriptStarted = false, false
		case strings.HasPrefix(part, "-DOCUMENTATION-"):
			inDoc, docStarted = true, false
			inScript, scriptStarted = false, false
			continue
		case strings.HasPrefix(part, "-SCRIPT-"):
			inDoc, docStarted = false, false
			inScript, scriptStarted = true, false
			continue
		case inDoc:
			docStarted = true
		case inScript:
			scriptStarted = true
		}

		if inDoc {
			ds.Documentation += strings.TrimRight(p, "\t ") + "\n"
		} else if inScript {
			ds.Script.Content += p + "\n"
		}
	}
	ds.Documentation = strings.TrimSpace(ds.GetDocumentation())
	ds.Script.Content = strings.TrimSpace(ds.GetScript().GetContent()) + "\n"
	return
}

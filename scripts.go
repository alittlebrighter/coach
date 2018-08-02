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
	script := store.GetScript([]byte(alias))
	if script == nil {
		return errors.New("not found")
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

	if _, err := tmpfile.Write(MarshalEdit(*script)); err != nil {
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

	var newScript models.DocumentedScript
	if newScript, err = UnmarshalEdit(string(newContents)); err != nil {
		return err
	}
	if len(newScript.GetAlias()) == 0 {
		newScript.Alias = string(RandomID())
	}
	newScript.Id = []byte(newScript.GetAlias())
	if newScript.GetAlias() != script.GetAlias() {
		store.DeleteScript(script.GetId())
	}
	if err := store.Save(newScript.GetId(), newScript); err != nil {
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
	GetScript(id []byte) *models.DocumentedScript
	QueryScripts(...string) ([]models.DocumentedScript, error)
	DeleteScript(id []byte) error
	IgnoreCommand(string) error
}

const doNotEditLine = "!DO NOT EDIT THIS LINE!"

func MarshalEdit(s models.DocumentedScript) []byte {
	var contents strings.Builder
	contents.WriteString("-ALIAS- = " + s.GetAlias() + "\n")
	contents.WriteString("-TAGS- = " + strings.Join(s.GetTags(), ",") + "\n\n")
	contents.WriteString("-DOCUMENTATION- " + doNotEditLine + "\n")
	contents.WriteString(s.GetDocumentation() + "\n\n")
	contents.WriteString("-SCRIPT- " + doNotEditLine + "\n")
	contents.WriteString(s.GetScript().GetContent() + "\n")
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
	var inDoc, inScript bool
	for _, p := range parts {
		part := strings.TrimSpace(p)
		if len(part) == 0 {
			continue
		}

		switch {
		case strings.HasPrefix(part, "-ALIAS- ="):
			ds.Alias = strings.TrimSpace(strings.Split(part, "=")[1])

			inDoc = false
			inScript = false
		case strings.HasPrefix(part, "-TAGS- ="):
			tags := strings.Split(strings.Split(part, "=")[1], ",")
			for i := range tags {
				tags[i] = strings.TrimSpace(tags[i])
			}
			ds.Tags = tags

			inDoc = false
			inScript = false
		case strings.HasPrefix(part, "-DOCUMENTATION-"):
			inDoc = true
			inScript = false
			continue
		case strings.HasPrefix(part, "-SCRIPT-"):
			inScript = true
			inDoc = false
			continue
		}

		if inDoc {
			ds.Documentation = ds.GetDocumentation() + part + "\n"
		} else if inScript {
			ds.Script.Content += part + "\n"
		}
	}
	ds.Documentation = strings.TrimSpace(ds.GetDocumentation())
	ds.Script.Content = strings.TrimSpace(ds.GetScript().GetContent())
	return
}

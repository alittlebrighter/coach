// Copyright (c) 2018, Adam Bright <brightam1@gmail.com>
// See LICENSE for licensing information

package coach

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/golang/protobuf/ptypes"
	"github.com/rs/xid"

	models "github.com/alittlebrighter/coach-pro/gen/proto"
	"github.com/alittlebrighter/coach-pro/platforms"
	"github.com/alittlebrighter/coach-pro/storage/database"
)

const Header = "exported from COACH - https://github.com/alittlebrighter/coach"

func QueryScripts(query string, store ScriptStore) (scripts []models.DocumentedScript, err error) {
	if len(query) == 0 {
		err = errors.New("no query specified")
		return
	}

	scripts, err = store.QueryScripts(strings.Split(strings.TrimSpace(query), ",")...)
	return
}

func SaveScript(toSave models.DocumentedScript, overwrite bool, store ScriptStore) (err error) {
	if toSave.GetScript() == nil || len(strings.TrimSpace(toSave.GetScript().GetContent())) == 0 {
		return errors.New("no script to save")
	}

	if toSave.GetAuditLog() == nil {
		toSave.AuditLog = new(models.AuditLog)
	}

	currentUser, _ := user.Current()
	now := ptypes.TimestampNow()
	if len(toSave.GetAuditLog().GetCreatedBy()) == 0 {
		toSave.AuditLog.Created = now
		toSave.AuditLog.CreatedBy = currentUser.Username
	}

	toSave.AuditLog.Updated = now
	toSave.AuditLog.UpdatedBy = currentUser.Username

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
	err = IgnoreCommand(toSave.GetScript().GetContent(), store)
	return
}

func EditScript(alias string, store ScriptStore) (*models.DocumentedScript, error) {
	script := store.GetScript([]byte(alias))
	store.Close()
	if script == nil {
		script = &models.DocumentedScript{
			Id:     []byte(alias),
			Alias:  alias,
			Tags:   strings.Split(alias, "."),
			Script: &models.Script{Shell: platforms.IdentifyShell()},
		}
	}

	shell := platforms.GetShell(script.GetScript().GetShell())

	var tmpfile *os.File
	var err error
	for i := 0; i < 10; i++ {
		name := filepath.Join(os.TempDir(), "coach"+xid.New().String()+"."+shell.FileExtension())
		tmpfile, err = os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
		if os.IsExist(err) {
			continue
		}
		break
	}
	if err != nil {
		return nil, err
	}

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

	if err := platforms.OpenEditor(tmpfile.Name()); err != nil {
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

func RunScript(script models.DocumentedScript, args []string, configureIO func(*exec.Cmd)) error {
	shell := platforms.GetShell(script.GetScript().GetShell())

	toRun, cleanup, err := shell.BuildCommand(script.GetScript().GetContent(), args)
	if cleanup != nil {
		defer cleanup()
	}
	if err != nil {
		return err
	}
	configureIO(toRun)
	toRun.Run()

	return nil
}

func DeleteScript(alias string, store ScriptStore) error {
	trash := store.GetScript([]byte(alias))
	if trash == nil {
		return database.ErrNotFound
	}
	trash.Alias = database.TrashTag + "." + trash.GetAlias()
	trash.Documentation = "-TAGS- = " + strings.Join(trash.GetTags(), ",") + platforms.Newline(1) +
		trash.GetDocumentation()
	trash.Tags = []string{database.TrashTag}

	store.DeleteScript(trash.GetId())
	return SaveScript(*trash, true, store)
}

func RestoreScript(alias string, store ScriptStore) (*models.DocumentedScript, error) {
	restore := store.GetScript([]byte(database.TrashTag + "." + alias))
	if restore == nil {
		return nil, database.ErrNotFound
	}
	store.DeleteScript(restore.GetId())
	restore.Alias = strings.TrimPrefix(restore.GetAlias(), database.TrashTag+".")
	firstNewline := strings.Index(restore.GetDocumentation(), platforms.Newline(1))
	UnmarshalLine(restore.GetDocumentation()[:firstNewline], restore) // restore Tags
	restore.Documentation = restore.GetDocumentation()[firstNewline:]

	return restore, SaveScript(*restore, false, store)
}

type ScriptStore interface {
	Save(id []byte, value interface{}, overwrite bool) error
	GetScript(id []byte) *models.DocumentedScript
	QueryScripts(...string) ([]models.DocumentedScript, error)
	DeleteScript(id []byte) error
	IgnoreStore
	Closable
}

const doNotEditLine = "!DO NOT EDIT THIS LINE!"

func MarshalEdit(s models.DocumentedScript) []byte {
	c := platforms.GetShell(s.GetScript().GetShell()).LineComment()
	nl := platforms.Newline(1)

	var contents strings.Builder
	contents.WriteString(c + "-ALIAS- = " + s.GetAlias() + nl)
	contents.WriteString(c + " -TAGS- = " + strings.Join(s.GetTags(), ",") + nl)
	contents.WriteString(c + "-SHELL- = " + s.GetScript().GetShell() + nl + c + nl)
	contents.WriteString(c + "-DOCUMENTATION- " + doNotEditLine + nl)
	for _, line := range strings.Split(s.GetDocumentation(), nl) {
		contents.WriteString(c + " " + strings.TrimRight(line, "\t ") + nl)
	}
	contents.WriteString(c + nl)
	contents.WriteString(c + "-SCRIPT- " + doNotEditLine + nl)
	contents.WriteString(strings.TrimSpace(s.GetScript().GetContent()) + nl)
	return []byte(contents.String())
}

func UnmarshalEdit(contents string) (ds models.DocumentedScript, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("could not parse file")
		}
	}()

	ds.Script = new(models.Script)

	parts := strings.Split(contents, platforms.Newline(1))
	var inDoc, docStarted, inScript, scriptStarted bool
	for _, p := range parts {
		part := strings.TrimSpace(p)
		if !docStarted && !scriptStarted && len(part) == 0 {
			continue
		}

		switch {
		case strings.Contains(part, "-ALIAS- ="):
			ds.Alias = strings.TrimSpace(strings.Split(part, "=")[1])

			inDoc, docStarted = false, false
			inScript, scriptStarted = false, false
		case strings.Contains(part, "-TAGS- ="):
			tags := strings.Split(strings.Split(part, "=")[1], ",")
			for i := range tags {
				tags[i] = strings.TrimSpace(tags[i])
			}
			ds.Tags = tags

			inDoc, docStarted = false, false
			inScript, scriptStarted = false, false
		case strings.Contains(part, "-SHELL- ="):
			ds.Script.Shell = strings.TrimSpace(strings.Split(part, "=")[1])

			inDoc, docStarted = false, false
			inScript, scriptStarted = false, false
		case strings.Contains(part, "-DOCUMENTATION- "+doNotEditLine):
			inDoc, docStarted = true, false
			inScript, scriptStarted = false, false
			continue
		case strings.Contains(part, "-SCRIPT- "+doNotEditLine):
			inDoc, docStarted = false, false
			inScript, scriptStarted = true, false
			continue
		case inDoc:
			docStarted = true
		case inScript:
			scriptStarted = true
		}

		if inDoc {
			ds.Documentation += strings.TrimRight(strings.Replace(strings.TrimLeft(p, "/#"), " ", "", 1), "\t ") + platforms.Newline(1)
		} else if inScript {
			ds.Script.Content += p + platforms.Newline(1)
		}
	}
	ds.Documentation = strings.TrimSpace(ds.GetDocumentation())
	ds.Script.Content = strings.TrimSpace(ds.GetScript().GetContent()) + platforms.Newline(1)
	return
}

func UnmarshalLine(line string, ds *models.DocumentedScript) (processed bool) {
	switch part := strings.TrimSpace(line); {
	case strings.Contains(part, "-ALIAS- ="):
		ds.Alias = strings.TrimSpace(strings.Split(part, "=")[1])

		processed = true
	case strings.Contains(part, "-TAGS- ="):
		tags := strings.Split(strings.Split(part, "=")[1], ",")
		for i := range tags {
			tags[i] = strings.TrimSpace(tags[i])
		}
		ds.Tags = tags
		processed = true
	case strings.Contains(part, "-SHELL- ="):
		ds.Script.Shell = strings.TrimSpace(strings.Split(part, "=")[1])
		processed = true
	case strings.Contains(part, Header):
		fallthrough
	case strings.Contains(part, "-DOCUMENTATION- "+doNotEditLine):
		fallthrough
	case strings.Contains(part, "-SCRIPT- "+doNotEditLine):
		processed = true
	default:
		processed = false
	}
	return
}

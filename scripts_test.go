package coach

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alittlebrighter/coach-pro/gen/models"
)

var (
	scriptFile = `-ALIAS- = myscript
 -TAGS- = cheese,whiskey,mashed potatoes
-SHELL- = bash

-DOCUMENTATION- ` + doNotEditLine + `
A simple script that does my things.

	A second line.

-SCRIPT- ` + doNotEditLine + `
ls -a
	cat ~/.bashrc | grep "BILLY BOB"

	restart critical service
`

	scriptStruct = models.DocumentedScript{
		Alias: "myscript",
		Tags:  []string{"cheese", "whiskey", "mashed potatoes"},
		Documentation: `A simple script that does my things.

	A second line.`,
		Script: &models.Script{
			Content: `ls -a
	cat ~/.bashrc | grep "BILLY BOB"

	restart critical service
`,
			Shell: "bash",
		},
	}
)

func TestMarshalEdit(t *testing.T) {
	assert.Equal(t, scriptFile, string(MarshalEdit(scriptStruct)))
}

func TestUnmarshalEdit(t *testing.T) {
	newScript, _ := UnmarshalEdit(scriptFile)
	assert.Equal(t, scriptStruct.GetAlias(), newScript.GetAlias())
	for i := range scriptStruct.GetTags() {
		assert.Equal(t, scriptStruct.GetTags()[i], newScript.GetTags()[i])
	}
	assert.Equal(t, scriptStruct.GetDocumentation(), newScript.GetDocumentation())
	assert.Equal(t, scriptStruct.GetScript().GetContent(), newScript.GetScript().GetContent())
}

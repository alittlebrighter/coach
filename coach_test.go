package coach

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCleanupCommand(t *testing.T) {
	// input is the key, expected output is the value
	tests := map[string]string{
		"4 echo 'hello'": "echo 'hello'",
		"cat   	one/two/three  	": "cat one/two/three",
		" exec  	mycmd": "exec mycmd",
		"  	890   ": "",
	}

	for in, expectedOut := range tests {
		assert.Equal(t, expectedOut, CleanupCommand(in))
	}
}

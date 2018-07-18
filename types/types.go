package types

import (
	"time"
)

type (
	Cmd struct {
		Command string
		Args    []string
		Context *CmdContext
	}

	SessionContext struct {
		OS       string
		Shell    string
		Hostname string
	}

	CmdContext struct {
		*SessionContext
		Directory string // FK
		ExecTime  time.Time
	}
)

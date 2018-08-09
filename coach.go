// Copyright (c) 2018, Adam Bright <brightam1@gmail.com>
// See LICENSE for licensing information

//go:generate protoc -I protobuf/ --go_out=plugins=grpc:gen/models protobuf/coach.proto
package coach

import (
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/rs/xid"

	"github.com/alittlebrighter/coach-pro/platforms"
)

var (
	Shell platforms.Platform
)

func init() {
	Shell = &platforms.Bash{}
}

func RandomID() (id []byte) {
	id = xid.New().Bytes()
	return
}

type Closable interface {
	Close() error
}

func HomeDir() string {
	var home string
	envSetHome := os.Getenv("COACH_HOME")
	defaultAppDir := platforms.DefaultHomeDir()
	_, sysErr := os.Stat(defaultAppDir + "/coach")
	switch {
	case len(envSetHome) > 0:
		home = envSetHome
	case sysErr == nil:
		home = defaultAppDir + "/coach"
	default:
		homeDir, _ := homedir.Dir()
		home = homeDir + "/.coach"
	}
	return home
}

// Copyright (c) 2018, Adam Bright <brightam1@gmail.com>
// See LICENSE for licensing information

//go:generate protoc -I protobuf/ --go_out=plugins=grpc:gen/models protobuf/coach.proto
package coach

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/rs/xid"

	"github.com/alittlebrighter/coach-pro/platforms"
	"github.com/alittlebrighter/coach-pro/storage/database"
)

var (
	Shell  platforms.Platform
	DBPath string
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

func GetStore(readonly bool) *database.BoltDB {
	store, err := database.NewBoltDB(DBPath, readonly)
	if err != nil {
		fmt.Println("could not access db:", err)
		os.Exit(1)
	}
	return store
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

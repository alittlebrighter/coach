// Copyright (c) 2018, Adam Bright <brightam1@gmail.com>
// See LICENSE for licensing information

//go:generate protoc -I protobuf/ --go_out=plugins=grpc:gen/proto protobuf/coach.proto protobuf/service.proto

package coach

import (
	"fmt"
	"os"

	"github.com/rs/xid"

	"github.com/alittlebrighter/coach-pro/storage/database"
)

var (
	DBPath string
)

func RandomID() (id []byte) {
	id = xid.New().Bytes()
	return
}

type Closable interface {
	Close() error
}

// TODO: this needs to return an interface
func GetStore(readonly bool) *database.BoltDB {
	store, err := database.NewBoltDB(DBPath, readonly)
	if err != nil {
		fmt.Println("could not access db:", err)
		os.Exit(1)
	}
	return store
}

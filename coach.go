// Copyright (c) 2018, Adam Bright <brightam1@gmail.com>
// See LICENSE for licensing information

//go:generate protoc -I protobuf/ --go_out=plugins=grpc:gen/models protobuf/coach.proto
package coach

import (
	"github.com/rs/xid"

	"github.com/alittlebrighter/coach/platforms"
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

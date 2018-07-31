//go:generate protoc -I protobuf/ --go_out=plugins=grpc:gen/models protobuf/coach.proto
package coach

import (
	"github.com/rs/xid"

	"github.com/alittlebrighter/coach/platforms"
)

var (
	Shell platforms.Shell
)

func init() {
	Shell = &platforms.Bash{}
}

func RandomID() (id []byte) {
	id = xid.New().Bytes()
	return
}

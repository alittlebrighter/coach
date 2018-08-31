package main

import (
	"context"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync"

	"google.golang.org/grpc"

	"github.com/alittlebrighter/coach-pro"
	pb "github.com/alittlebrighter/coach-pro/gen/proto"
	"github.com/alittlebrighter/coach-pro/platforms"
	"github.com/alittlebrighter/coach-pro/storage/database"
)

const EOF = "!!!EOF!!!"

func main() {
	// Find home directory.
	home := coach.HomeDir()
	os.Mkdir(home, os.ModePerm)
	coach.DBPath = home + "/coach.db"

	svc := &CoachRPC{GetStore: coach.GetStore}

	// listen on 8326 (TEAM on a number pad)
	listen, err := net.Listen("tcp", ":8326")
	if err != nil {
		log.Fatalf("failed to listen for tcp connections: %v", err)
	}

	rpcServer := grpc.NewServer()
	pb.RegisterCoachRPCServer(rpcServer, svc)
	if err := rpcServer.Serve(listen); err != nil {
		log.Fatalf("CoachRPC failed to serve connections: %v", err)
	}
}

type CoachRPC struct {
	GetStore func(bool) *database.BoltDB
}

func (c *CoachRPC) Scripts(ctx context.Context, query *pb.ScriptsQuery) (*pb.GetScriptsResponse, error) {
	store := c.GetStore(true)
	scripts, err := coach.QueryScripts(query.TagQuery, store)
	store.Close()

	response := &pb.GetScriptsResponse{Scripts: []*pb.DocumentedScript{}}
	for i := range scripts {
		response.Scripts = append(response.Scripts, &scripts[i])
	}

	return response, err
}

func (c *CoachRPC) RunScript(streams pb.CoachRPC_RunScriptServer) error {
	log.Println("started RunScript")
	initEvent, err := streams.Recv()
	if err != nil {
		return err
	}

	var alias string
	args := []string{}
	cmdArgs := strings.Fields(initEvent.GetInput())
	switch {
	case len(cmdArgs) > 1:
		args = cmdArgs[1:]
		fallthrough
	case len(cmdArgs) == 1:
		alias = cmdArgs[0]
	default:
		return database.ErrNotFound
	}

	store := coach.GetStore(true)
	toRun := store.GetScript([]byte(alias))
	store.Close()
	if toRun == nil {
		return database.ErrNotFound
	}

	IOStreams := new(ioStreams)
	ioWG := new(sync.WaitGroup)
	ioWG.Add(1)

	// run the script
	runErr := make(chan error)
	go func() {
		runErr <- coach.RunScript(*toRun, args, c.configureCmdIO(IOStreams, ioWG))
	}()

	// wait for stdin/out/err to be initialized
	ioWG.Wait()

	// proxy the output
	outClose := make(chan bool)
	output := make(chan string)
	go func() {
		buffer := make([]byte, initEvent.GetResponseSize())
		for {
			readCount, err := IOStreams.Stdout.Read(buffer)
			switch {
			case err == nil:
				output <- string(buffer[:readCount-1])
			case err == io.EOF && readCount > 0:
				output <- string(buffer[:readCount-1])
				fallthrough
			case err == io.EOF && readCount == 0:
				output <- EOF
				return
			default:
				return
			}
		}
	}()
	go func() {
		for {
			select {
			case out := <-output:
				if out == EOF {
					streams.Send(&pb.RunEventOut{Output: EOF})
				}
				if len(out) > 0 {
					streams.Send(&pb.RunEventOut{Output: out})
				}
			case <-outClose:
				return
			}
		}
	}()

	// proxy the input
	inClose := make(chan bool)
	input := make(chan string)
	go func() {
		for {
			inEvent, err := streams.Recv()
			if err == io.EOF {
				return
			}
			input <- inEvent.GetInput()
		}
	}()
	go func() {
		for {
			select {
			case in := <-input:
				IOStreams.Stdin.Write([]byte(strings.TrimSpace(in) + platforms.Newline(1)))
			case <-inClose:
				IOStreams.Stdin.Close()
				return
			}
		}
	}()

	err = <-runErr
	outClose <- true
	inClose <- true
	log.Println("finished RunScript")
	return err
}

func (c *CoachRPC) configureCmdIO(streams *ioStreams, lock *sync.WaitGroup) func(*exec.Cmd) error {
	return func(cmd *exec.Cmd) error {
		defer lock.Done()

		var err error

		if streams.Stdin, err = cmd.StdinPipe(); err != nil {
			return err
		}
		if streams.Stdout, err = cmd.StdoutPipe(); err != nil {
			return err
		}
		if streams.Stderr, err = cmd.StderrPipe(); err != nil {
			return err
		}
		return nil
	}
}

type ioStreams struct {
	Stdin  io.WriteCloser
	Stdout io.ReadCloser
	Stderr io.ReadCloser
}

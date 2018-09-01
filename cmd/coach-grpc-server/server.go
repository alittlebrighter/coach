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
	wg := new(sync.WaitGroup)
	wg.Add(1)

	// run the script
	runErr := make(chan error)
	go func() {
		runErr <- coach.RunScript(*toRun, args, c.configureCmdIO(IOStreams, wg))
	}()

	// wait for stdin/out/err to be initialized
	wg.Wait()

	// proxy the output
	output := make(chan string)
	wg.Add(1)
	go func() {
		parseStream(IOStreams.Stdout, output, initEvent.GetResponseSize())
		log.Println("closed stdout")
	}()

	stderr := make(chan string)
	wg.Add(1)
	go func() {
		parseStream(IOStreams.Stderr, stderr, initEvent.GetResponseSize())
		log.Println("closed stderr")
	}()

	go func() {
		for out := range output {
			if len(out) == 0 {
				continue
			}

			streams.Send(&pb.RunEventOut{Output: out})
		}
		err := streams.Send(&pb.RunEventOut{Output: EOF})
		wg.Done()
		log.Println("sent EOF to Output", err)
	}()

	go func() {
		for err := range stderr {
			if len(err) == 0 {
				continue
			}

			streams.Send(&pb.RunEventOut{Error: err})
		}
		err := streams.Send(&pb.RunEventOut{Error: EOF})
		wg.Done()
		log.Println("sent EOF to Error", err)
	}()

	// proxy the input
	input := make(chan string)
	go func() {
	main:
		for {
			inEvent, err := streams.Recv()
			if err == io.EOF || inEvent.GetInput() == EOF {
				close(input)
				break main
			}
			input <- inEvent.GetInput()
		}
		log.Println("input stream closed")
	}()
	go func() {
		for in := range input {
			IOStreams.Stdin.Write([]byte(strings.TrimSpace(in) + platforms.Newline(1)))
		}
		IOStreams.Stdin.Close()
		log.Println("closed stdin stream")
	}()

	err = <-runErr
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

func parseStream(stream io.ReadCloser, output chan string, buffSize uint32) {
	buffer := make([]byte, buffSize)
parseLoop:
	for {
		readCount, err := stream.Read(buffer)
		switch {
		case err == nil:
			output <- string(buffer[:readCount-1])
		case err == io.EOF && readCount > 0:
			output <- string(buffer[:readCount-1])
			fallthrough
		case err == io.EOF && readCount == 0:
			output <- EOF
			fallthrough
		default:
			break parseLoop
		}
	}
	close(output)
	stream.Close()
}

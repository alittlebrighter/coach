package coachrpc

import (
	"context"
	"io"
	"log"
	"os/exec"
	"strings"
	"sync"

	coach "github.com/alittlebrighter/coach"
	pb "github.com/alittlebrighter/coach/gen/proto"
	"github.com/alittlebrighter/coach/platforms"
	"github.com/alittlebrighter/coach/storage/database"
)

const (
	EOF           = "!!!EOF!!!"
	CancelCommand = "coach::cancelRun"
)

type CoachRPC struct {
	GetStore func(bool) *database.BoltDB
}

func (c *CoachRPC) QueryScripts(ctx context.Context, query *pb.ScriptsQuery) (*pb.GetScriptsResponse, error) {
	store := c.GetStore(true)
	scripts, err := coach.QueryScripts(query.GetQuery(), store)
	store.Close()

	response := &pb.GetScriptsResponse{Scripts: []*pb.DocumentedScript{}}
	for i := range scripts {
		response.Scripts = append(response.Scripts, &scripts[i])
	}

	return response, err
}

func (c *CoachRPC) GetScript(ctx context.Context, alias *pb.ScriptsQuery) (script *pb.DocumentedScript, err error) {
	store := c.GetStore(true)
	defer store.Close()
	script = store.GetScript([]byte(alias.GetQuery()))
	if script == nil {
		err = database.ErrNotFound
	}
	return
}

func (c *CoachRPC) SaveScript(ctx context.Context, script *pb.SaveScriptRequest) (*pb.Response, error) {
	store := c.GetStore(false)
	err := coach.SaveScript(*script.GetScript(), script.GetOverwrite(), store)
	store.Close()

	response := &pb.Response{Success: err == nil}
	if err != nil {
		response.Error = err.Error()
	}

	return response, err
}

func (c *CoachRPC) RunScript(streams pb.CoachRPC_RunScriptServer) error {
	log.Println("started RunScript")
	initEvent, err := streams.Recv()
	if err != nil {
		return err
	}

	shutdown := false

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
	defer IOStreams.Close()

	wg := new(sync.WaitGroup)
	wg.Add(1)

	// run the script
	runErr := make(chan error)
	ctx, cancelRun := context.WithCancel(context.Background())
	go func() {
		runErr <- coach.RunScript(ctx, *toRun, args, c.configureCmdIO(IOStreams, wg))
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
			log.Println("sent StdOut:", out)
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
			log.Println("sent StdErr:", err)
		}
		err := streams.Send(&pb.RunEventOut{Error: EOF})
		wg.Done()
		log.Println("sent EOF to Error", err)
	}()

	// proxy the input
	input := make(chan string)
	go func() {
	main:
		for !shutdown {
			inEvent, err := streams.Recv()
			if err == io.EOF || inEvent.GetInput() == EOF {
				break main
			}
			input <- inEvent.GetInput()
		}
		close(input)
		log.Println("input stream closed")
	}()
	go func() {
		for in := range input {
			if in == CancelCommand {
				cancelRun()
				return
			}
			IOStreams.Stdin.Write([]byte(strings.TrimSpace(in) + platforms.Newline(1)))
		}
		log.Println("closed stdin stream")
	}()

	err = <-runErr
	shutdown = true
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

func (s *ioStreams) Close() {
	s.Stdin.Close()
	s.Stderr.Close()
	s.Stdout.Close()
}

func parseStream(stream io.ReadCloser, output chan string, buffSize uint32) {
	buffer := make([]byte, buffSize)
parseLoop:
	for {
		readCount, err := stream.Read(buffer)

		switch {
		case err == nil:
			for _, line := range processOutput(buffer[:readCount]) {
				output <- line
			}
		case err == io.EOF && readCount > 0:
			for _, line := range processOutput(buffer[:readCount]) {
				output <- line
			}
			fallthrough
		case err == io.EOF && readCount == 0:
			output <- EOF
			fallthrough
		default:
			break parseLoop
		}
	}
	close(output)
}

func processOutput(data []byte) []string {
	str := string(data)
	return strings.Split(str, platforms.Newline(1))
}

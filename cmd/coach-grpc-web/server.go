package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"

	pb "github.com/alittlebrighter/coach-pro/gen/proto"
)

const EOF = "!!!EOF!!!"

func main() {
	conn, err := grpc.Dial(":4444", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewCoachRPCClient(conn)

	coachScripts, err := client.Scripts(context.Background(), &pb.ScriptsQuery{TagQuery: "coach"})
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Scripts found:")
	for _, script := range coachScripts.GetScripts() {
		fmt.Println(script.GetAlias(), "-", script.GetDocumentation())
	}

	streams, err := client.RunScript(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	incoming := make(chan *pb.RunEventOut)
	go func() {
		for {
			outEvent, err := streams.Recv()
			if err == io.EOF || outEvent.GetOutput() == EOF {
				if len(outEvent.GetOutput()) > 0 && outEvent.GetOutput() != EOF {
					incoming <- outEvent
				}
				incoming <- &pb.RunEventOut{Output: EOF}
				break
			} else if err != nil {
				log.Println("ERROR:", err)
				incoming <- &pb.RunEventOut{Output: EOF}
				break
			} else {
				incoming <- outEvent
			}
		}
		close(incoming)
		wg.Done()
	}()

	fmt.Println("running hello-world")
	streams.Send(&pb.RunEventIn{Input: "hello-world", ResponseSize: 128})
	stdinReader := bufio.NewReader(os.Stdin)

inputLoop:
	for {
		timer := time.NewTimer(500 * time.Millisecond)
		select {
		case event, ok := <-incoming:
			if event.GetOutput() == EOF || !ok {
				break inputLoop
			}
			fmt.Println(event.GetOutput())
		case <-timer.C:
			in, err := stdinReader.ReadString('\n')
			if len(strings.TrimSpace(in)) > 0 && err == nil {
				streams.Send(&pb.RunEventIn{Input: in})
			}
		}
	}

	streams.CloseSend()

	wg.Wait()
}

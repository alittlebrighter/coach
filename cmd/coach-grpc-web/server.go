package main

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/json-iterator/go"
	"github.com/rs/xid"
	"google.golang.org/grpc"

	pb "github.com/alittlebrighter/coach-pro/gen/proto"
)

const EOF = "!!!EOF!!!"

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func main() {
	appCtx, err := NewAppContext()
	if err != nil {
		log.Fatalln(err)
	}
	defer appCtx.CloseRPC()

	startWebsocketServer(appCtx)
}

const webUri = "localhost:8327"

func startWebsocketServer(appCtx *appContext) {
	http.HandleFunc("/rpc", appCtx.rpc)
	log.Fatal(http.ListenAndServe(webUri, nil))
}

type appContext struct {
	rpcConn   *grpc.ClientConn
	rpcClient pb.CoachRPCClient

	ActiveInputs map[string]chan *RPCCall
}

func NewAppContext() (*appContext, error) {
	ctx := &appContext{ActiveInputs: map[string]chan *RPCCall{}}

	var err error
	ctx.rpcConn, err = grpc.Dial(":8326", grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	ctx.rpcClient = pb.NewCoachRPCClient(ctx.rpcConn)

	return ctx, nil
}

func (a *appContext) CloseRPC() error {
	return a.rpcConn.Close()
}

func (a *appContext) GetScripts(req *RPCCall, out chan *RPCCall) {
	if len(req.Input) == 0 {
		req.Input = "?"
	}

	scripts, err := a.rpcClient.Scripts(context.Background(), &pb.ScriptsQuery{TagQuery: req.Input})
	if err != nil {
		req.Error = "rpc client: " + err.Error()
		out <- req
		return
	}

	req.Output = scripts
	out <- req
	return
}

func (a *appContext) RunScript(req *RPCCall, in, out chan *RPCCall) {
	log.Println("started runScript")
	streams, err := a.rpcClient.RunScript(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	incoming := make(chan *pb.RunEventOut)
	stdoutClosed, stderrClosed := false, false
	go func() {
	main:
		for !stdoutClosed || !stderrClosed {
			outEvent, err := streams.Recv()
			if outEvent != nil {
				incoming <- outEvent
			}

			if err != nil {
				log.Println("incoming stream:", err)
				incoming <- &pb.RunEventOut{Output: EOF}
				incoming <- &pb.RunEventOut{Error: EOF}
				break main
			}
		}
		close(incoming)
		wg.Done()
		log.Println("stopped receiving response stream")
	}()

	streams.Send(&pb.RunEventIn{Input: req.Input, ResponseSize: 256})

	for !stdoutClosed || !stderrClosed {
		select {
		case event, chanOk := <-incoming:
			// hackish way of copying req and getting pointer to the copy
			response := &(*req)

			stdoutClosed = stdoutClosed || event.GetOutput() == EOF
			stderrClosed = stderrClosed || event.GetError() == EOF
			switch {
			case !chanOk:
				log.Println("!chanOk")
				fallthrough
			case len(event.GetOutput()) == 0 && len(event.GetError()) == 0:
				log.Println("setting both closed")
				stdoutClosed, stderrClosed = true, true
			}

			response.Output = event.GetOutput()
			response.Error = event.GetError()

			out <- response
		case input := <-in:
			streams.Send(&pb.RunEventIn{Input: input.Input})
		}
	}

	out <- &RPCCall{Id: req.Id, Method: req.Method, Output: EOF, Error: EOF}
	streams.Send(&pb.RunEventIn{Input: EOF})
	streams.CloseSend()
	log.Println("sent CloseSend")
	wg.Wait()

	stdoutClosed, stderrClosed = true, true

	close(in)
	delete(a.ActiveInputs, req.Id)
	log.Println("finished runScript")
}

var upgrader = websocket.Upgrader{} // use default options

func (a *appContext) rpc(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("ws upgrade:", err)
		return
	}
	defer c.Close()

	wsOut := make(chan *RPCCall, 5)

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		for {
			req := new(RPCCall)
			if err := c.ReadJSON(&req); err != nil {
				log.Println("ws read:", err)
				break
			}

			if len(req.Id) == 0 {
				req.Id = xid.New().String()

				if req.Method == "runScript" {
					input := make(chan *RPCCall, 5)
					a.ActiveInputs[req.Id] = input
					go a.RunScript(req, input, wsOut)
					continue
				}
			}

			switch req.Method {
			case "getScripts":
				go a.GetScripts(req, wsOut)
			case "runScript":
				input, exists := a.ActiveInputs[req.Id]
				if !exists {
					continue
				}

				input <- req
			}
		}
		wg.Done()
	}()

	go func() {
		for out := range wsOut {
			err = c.WriteJSON(out)
			if err != nil {
				log.Println("ws write:", err)
				break
			}
		}
		close(wsOut)
		wg.Done()
	}()

	wg.Wait()
}

type RPCCall struct {
	Id     string      `json:"id"`
	Method string      `json:"method"`
	Input  string      `json:"input"`
	Output interface{} `json:"output"`
	Error  string      `json:"error"`
}

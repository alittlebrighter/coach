package main

import (
	"context"
	libJson "encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gobuffalo/packr"
	"github.com/gorilla/websocket"
	"github.com/json-iterator/go"
	"github.com/rs/xid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/alittlebrighter/coach/config"
	pb "github.com/alittlebrighter/coach/gen/proto"
	"github.com/alittlebrighter/coach/grpc"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func appMain(cmd *cobra.Command, args []string) {
	webUri := viper.GetString("web.host")
	rpcUri := viper.GetString("rpc.host")

	appCtx, err := NewAppContext(rpcUri)
	if err != nil {
		log.Fatalln(err)
	}
	defer appCtx.CloseRPC()

	startWebsocketServer(appCtx, webUri)
}

func startWebsocketServer(appCtx *appContext, serveAt string) {
	box := packr.NewBox("./web/coach-ui/dist")

	http.Handle("/", http.FileServer(box))
	http.HandleFunc("/rpc", appCtx.rpc)

	log.Println("serving coach-ui at " + serveAt)
	log.Fatal(http.ListenAndServe(serveAt, nil))
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "coach-grpc-web",
		Short: "Coach web UI.",
		Run:   appMain,
	}
	rootCmd.Flags().String("host", "localhost:26224", "Address to serve on.")
	rootCmd.Flags().String("rpc-host", ":8326", "URI for the coach gRPC server.")

	configure(rootCmd)

	rootCmd.Execute()
}

func configure(cmd *cobra.Command) {
	config.AppConfiguration()

	viper.SetDefault("rpc.host", ":8326")
	viper.SetDefault("web.host", "localhost:26224")
	viper.SetDefault("security.certificate_filepath", config.HomeDir()+"/security/coach.pem")

	viper.BindPFlag("rpc.host", cmd.Flags().Lookup("rpc-host"))
	viper.BindPFlag("web.host", cmd.Flags().Lookup("host"))
}

type appContext struct {
	rpcConn   *grpc.ClientConn
	rpcClient pb.CoachRPCClient

	ActiveInputs map[string]chan *RPCCall
}

func NewAppContext(rpcUri string) (*appContext, error) {
	ctx := &appContext{ActiveInputs: map[string]chan *RPCCall{}}

	var err error
	dialOpts := []grpc.DialOption{}
	if creds, err := credentials.NewClientTLSFromFile(viper.GetString("security.certificate_filepath"), ""); err == nil {
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(creds))
	} else {
		dialOpts = append(dialOpts, grpc.WithInsecure())
		log.Println("tls certificate:", err)
		log.Println("WARNING: dialing server running without TLS, all IO over the network is being transmitted in plaintext")
	}

	ctx.rpcConn, err = grpc.Dial(rpcUri, dialOpts...)
	if err != nil {
		return nil, err
	}
	ctx.rpcClient = pb.NewCoachRPCClient(ctx.rpcConn)

	return ctx, nil
}

func (a *appContext) CloseRPC() error {
	return a.rpcConn.Close()
}

func (a *appContext) QueryScripts(req *RPCCall, out chan *RPCCall) {
	if len(req.Input) == 0 {
		req.Input = []byte("?")
	}
	scripts, err := a.rpcClient.QueryScripts(context.Background(), &pb.ScriptsQuery{Query: BytesToString(req.Input)})
	if err != nil {
		req.Error = "rpc client: " + err.Error()
		out <- req
		return
	}

	req.Output = scripts
	out <- req
	return
}

func (a *appContext) GetScript(req *RPCCall, out chan *RPCCall) {
	script, err := a.rpcClient.GetScript(context.Background(), &pb.ScriptsQuery{Query: BytesToString(req.Input)})
	if err != nil {
		req.Error = "rpc client: " + err.Error()
		out <- req
		return
	}

	req.Output = script
	out <- req
	return
}

func (a *appContext) SaveScript(req *RPCCall, out chan *RPCCall) {
	response := &(*req)

	var script pb.DocumentedScript
	if err := json.Unmarshal(req.Input, &script); err != nil {
		log.Println("savescript:", err)
		response.Error = err.Error()
	}

	var err error
	response.Output, err = a.rpcClient.SaveScript(context.Background(), &pb.SaveScriptRequest{Script: &script, Overwrite: true})
	if err != nil {
		response.Error = err.Error()
	}

	response.Input = nil

	out <- response
}

func (a *appContext) RunScript(req *RPCCall, in, out chan *RPCCall) {
	streams, err := a.rpcClient.RunScript(context.Background())
	if err != nil {
		log.Println("runscript starting:", err)
		return
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
				incoming <- &pb.RunEventOut{Output: coachrpc.EOF}
				incoming <- &pb.RunEventOut{Error: coachrpc.EOF}
				break main
			}
		}
		close(incoming)
		wg.Done()
		log.Println("stopped receiving response stream")
	}()

	streams.Send(&pb.RunEventIn{Input: BytesToString(req.Input), ResponseSize: 256})

	for !stdoutClosed || !stderrClosed {
		select {
		case event, chanOk := <-incoming:
			// hackish way of copying req and getting pointer to the copy
			response := &RPCCall{Id: req.Id, Method: req.Method}

			stdoutClosed = stdoutClosed || event.GetOutput() == coachrpc.EOF
			stderrClosed = stderrClosed || event.GetError() == coachrpc.EOF
			switch {
			case !chanOk:
				fallthrough
			case len(event.GetOutput()) == 0 && len(event.GetError()) == 0:
				stdoutClosed, stderrClosed = true, true
			default:
				response.Output = event.GetOutput()
				response.Error = event.GetError()

				out <- response
			}
		case input := <-in:
			streams.Send(&pb.RunEventIn{Input: BytesToString(input.Input)})
		}
	}

	out <- &RPCCall{Id: req.Id, Method: req.Method, Output: coachrpc.EOF, Error: coachrpc.EOF}
	streams.Send(&pb.RunEventIn{Input: coachrpc.EOF})
	streams.CloseSend()
	wg.Wait()

	stdoutClosed, stderrClosed = true, true

	close(in)
	delete(a.ActiveInputs, req.Id)
}

var upgrader = websocket.Upgrader{} // use default options

func (a *appContext) rpc(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("ws upgrade:", err)
		return
	}
	defer c.Close()

	log.Println("new ws connection")

	wsOut := make(chan *RPCCall, 5)

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		for {
			req := new(RPCCall)
			if err := c.ReadJSON(&req); err != nil {
				log.Println("websocket read:", err)
				break
			}

			if len(req.Id) == 0 {
				req.Id = xid.New().String()
			}

			if _, exists := a.ActiveInputs[req.Id]; req.Method == "runScript" && !exists {
				input := make(chan *RPCCall, 5)
				a.ActiveInputs[req.Id] = input
				go a.RunScript(req, input, wsOut)
				continue
			}

			switch req.Method {
			case "queryScripts":
				go a.QueryScripts(req, wsOut)
			case "getScript":
				go a.GetScript(req, wsOut)
			case "runScript":
				input, exists := a.ActiveInputs[req.Id]
				if !exists {
					continue
				}

				input <- req
			case "saveScript":
				go a.SaveScript(req, wsOut)
			}
		}
		wg.Done()
	}()

	go func() {
		for out := range wsOut {
			err = c.WriteJSON(out)
			if err != nil {
				log.Println("websocket write:", err)
				break
			}
		}
		close(wsOut)
		wg.Done()
	}()

	wg.Wait()
}

type RPCCall struct {
	Id     string             `json:"id"`
	Method string             `json:"method"`
	Input  libJson.RawMessage `json:"input"`
	Output interface{}        `json:"output"`
	Error  string             `json:"error"`
}

func BytesToString(data []byte) string {
	return strings.Trim(string(data), `"`)
}

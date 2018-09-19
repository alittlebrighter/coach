package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	"github.com/alittlebrighter/coach-pro"
	pb "github.com/alittlebrighter/coach-pro/gen/proto"
	"github.com/alittlebrighter/coach-pro/grpc"
	"github.com/alittlebrighter/coach-pro/trial"
)

func appMain(cmd *cobra.Command, args []string) {
	rpcUri, _ := cmd.Flags().GetString("host")
	webUri, _ := cmd.Flags().GetString("web-host")

	svc := &coachrpc.CoachRPC{GetStore: coach.GetStore}

	rpcServer := grpc.NewServer()
	pb.RegisterCoachRPCServer(rpcServer, svc)

	log.Println(trial.ExpireNotice)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		log.Println("grpc server running at", rpcUri)
		defer wg.Done()

		// listen on 8326 (TEAM on a number pad)
		listen, err := net.Listen("tcp", rpcUri)
		if err != nil {
			log.Fatalf("failed to listen for tcp connections: %v\n", err)
		}

		if err := rpcServer.Serve(listen); err != nil {
			log.Fatalf("CoachRPC failed to serve connections: %v\n", err)
		}
	}()

	if len(webUri) > 0 {
		wg.Add(1)
		go func() {
			log.Println("grpc web server running at", webUri)
			defer wg.Done()
			log.Fatal(http.ListenAndServe(webUri, rpcServer))
		}()
	}

	wg.Wait()
}

func main() {
	// Find home directory.
	home := coach.HomeDir()
	os.Mkdir(home, os.ModePerm)
	coach.DBPath = home + "/coach.db"

	if _, err := os.Stat(coach.DBPath); err != nil {
		store := coach.GetStore(false)
		store.Init()
		store.Close()
	}

	rootCmd := &cobra.Command{
		Use:   "coach-grpc-server",
		Short: "Coach script library functions available over a gRPC interface.",
		Run:   appMain,
	}
	rootCmd.Flags().String("host", "localhost:8326", "URL to host GRPC server.")
	rootCmd.Flags().String("web-host", "", "URL to host the gRPC web server.  "+
		"Web server will not start if this value is blank.")

	rootCmd.Execute()
}

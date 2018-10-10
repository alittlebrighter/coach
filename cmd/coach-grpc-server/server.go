package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/alittlebrighter/coach"
	"github.com/alittlebrighter/coach/config"
	pb "github.com/alittlebrighter/coach/gen/proto"
	"github.com/alittlebrighter/coach/grpc"
)

func appMain(cmd *cobra.Command, args []string) {
	rpcUri := viper.GetString("rpc.host")
	webUri := viper.GetString("rpc.web-host")

	svc := &coachrpc.CoachRPC{GetStore: coach.GetStore}

	grpcOpts := []grpc.ServerOption{}
	if tlsOpt, err := credentials.NewServerTLSFromFile(viper.GetString("security.certificate_filepath"), viper.GetString("security.key_filepath")); err == nil {
		grpcOpts = append(grpcOpts, grpc.Creds(tlsOpt))
	} else {
		log.Println("tls certificate:", err)
		log.Println("WARNING: server running without TLS, all IO over the network is being transmitted in plaintext")
	}
	rpcServer := grpc.NewServer(grpcOpts...)
	pb.RegisterCoachRPCServer(rpcServer, svc)

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
			rpcServer.GracefulStop()
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
	home := config.HomeDir()
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

	configure(rootCmd)
	rootCmd.Execute()
}

func configure(cmd *cobra.Command) {
	config.AppConfiguration()

	viper.SetDefault("rpc.host", "localhost:8326")
	viper.SetDefault("rpc.web-host", "")
	viper.SetDefault("security.certificate_filepath", config.HomeDir()+"/security/coach.pem")
	viper.SetDefault("security.key_filepath", config.HomeDir()+"/security/coach_key.pem")

	viper.BindPFlag("rpc.host", cmd.Flags().Lookup("host"))
	viper.BindPFlag("rpc.web-host", cmd.Flags().Lookup("web-host"))
}

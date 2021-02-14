package main

import (
	"io/ioutil"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"

	"github.com/atreya2011/go-grpc-laughing-broccoli/gateway"
	pbExample "github.com/atreya2011/go-grpc-laughing-broccoli/proto"
	"github.com/atreya2011/go-grpc-laughing-broccoli/server"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"

	// Static files
	_ "github.com/atreya2011/go-grpc-laughing-broccoli/statik"
)

func main() {
	// Adds gRPC internal logs. This is quite verbose, so adjust as desired!
	log := grpclog.NewLoggerV2(os.Stdout, ioutil.Discard, ioutil.Discard)
	grpclog.SetLoggerV2(log)

	addr := "0.0.0.0:10000"
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(grpc_auth.UnaryServerInterceptor(server.ExampleAuthFunc)))

	pbExample.RegisterUserServiceServer(s, server.New())

	// Serve gRPC Server
	log.Info("Serving gRPC on https://", addr)
	go func() {
		log.Fatal(s.Serve(lis))
	}()

	err = gateway.Run("dns:///" + addr)
	log.Fatalln(err)
}

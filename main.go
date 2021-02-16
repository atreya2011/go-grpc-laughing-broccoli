package main

import (
	"io/ioutil"
	"net"
	"net/http"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"

	pbExample "github.com/atreya2011/go-grpc-laughing-broccoli/proto"
	"github.com/atreya2011/go-grpc-laughing-broccoli/server"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/improbable-eng/grpc-web/go/grpcweb"

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

	wrappedGrpc := grpcweb.WrapServer(s, grpcweb.WithOriginFunc(func(origin string) bool {
		return origin == "http://localhost:3005"
	}))
	grpcWebHandler := http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		wrappedGrpc.ServeHTTP(resp, req)
	})

	pbExample.RegisterUserServiceServer(s, server.New())

	// Serve gRPC Server
	log.Info("Serving gRPC on https://", addr)
	go func() {
		log.Fatal(s.Serve(lis))
	}()

	log.Fatalln(http.ListenAndServe(":10001", grpcWebHandler))

	// err = gateway.Run("dns:///" + addr)
	// log.Fatalln(err)
}

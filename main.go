package main

import (
	"fmt"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

func RunRPCServer(port int, register func(*grpc.Server, ...interface{}), myServer ...interface{}) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
		os.Exit(1)
	}
	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)
	register(grpcServer, myServer...)

	grpcServer.Serve(lis)
}

func main() {

}

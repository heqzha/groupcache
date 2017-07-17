package rpc

import (
	"fmt"
	"net"
	"os"

	"github.com/heqzha/dcache/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

func runRPCServer(port int, register func(*grpc.Server, ...interface{}), services ...interface{}) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
		os.Exit(1)
	}
	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)
	register(grpcServer, services...)

	grpcServer.Serve(lis)
}

func RunRPCServer(port int) {
	runRPCServer(port, func(grpc *grpc.Server, services ...interface{}) {
		for _, s := range services {
			switch s.(type) {
			case *DCacheService:
				pb.RegisterCacheServServer(grpc, s.(*DCacheService))
			}
		}
	}, new(DCacheService))
}

type DCacheService struct{}

func (s *DCacheService) Get(ctx context.Context, in *pb.GetReq) (*pb.GetRes, error) {
	return nil, nil
}

func (s *DCacheService) Set(ctx context.Context, in *pb.SetReq) (*pb.SetRes, error) {
	return nil, nil
}

func (s *DCacheService) Del(ctx context.Context, in *pb.DelReq) (*pb.DelRes, error) {
	return nil, nil
}

func (s *DCacheService) Register(ctx context.Context, in *pb.RegisterReq) (*pb.RegisterRes, error) {
	return nil, nil
}

func (s *DCacheService) Unregister(ctx context.Context, in *pb.UnregisterReq) (*pb.UnregisterRes, error) {
	return nil, nil
}

func (s *DCacheService) SyncSrvGroup(ctx context.Context, in *pb.SyncSrvGroupReq) (*pb.SyncSrvGroupRes, error) {
	return nil, nil
}

func (s *DCacheService) Ping(ctx context.Context, in *pb.PingReq) (*pb.PingRes, error) {
	return nil, nil
}

package core

import (
	"context"
	"fmt"
	"time"

	"github.com/heqzha/dcache/pb"
	"google.golang.org/grpc"
)

type CacheServClient struct {
	conn *grpc.ClientConn
	cli  pb.CacheServClient
}

func (c *CacheServClient) NewRPCClient(host string, port int, timeout time.Duration) error {
	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := grpc.Dial(addr, grpc.WithBlock(), grpc.WithTimeout(timeout), grpc.WithInsecure())
	if err != nil {
		return err
	}
	c.conn = conn
	c.cli = pb.NewCacheServClient(conn)
	return nil
}

func (c *CacheServClient) Get(group, key string) (*pb.GetRes, error) {
	return c.cli.Get(context.Background(), &pb.GetReq{
		Group: group,
		Key:   key,
	})
}

func (c *CacheServClient) Set(group, key string, value []byte) (*pb.SetRes, error) {
	return c.cli.Set(context.Background(), &pb.SetReq{
		Group: group,
		Key:   key,
		Value: value,
	})
}

func (c *CacheServClient) Del(group, key string) (*pb.DelRes, error) {
	return c.cli.Del(context.Background(), &pb.DelReq{
		Group: group,
		Key:   key,
	})
}

func (c *CacheServClient) Register(group, addr string) (*pb.RegisterRes, error) {
	return c.cli.Register(context.Background(), &pb.RegisterReq{
		Group: group,
		Addr:  addr,
	})
}

func (c *CacheServClient) Unregister(group, addr string) (*pb.UnregisterRes, error) {
	return c.cli.Unregister(context.Background(), &pb.UnregisterReq{
		Group: group,
		Addr:  addr,
	})
}

func (c *CacheServClient) SyncSrvGroups(srvgroups []byte) (*pb.SyncSrvGroupsRes, error) {
	return c.cli.SyncSrvGroups(context.Background(), &pb.SyncSrvGroupsReq{
		SrvGroups: srvgroups,
	})
}

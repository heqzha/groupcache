package test

import (
	"context"
	"testing"

	"github.com/heqzha/dcache/rpc"
)

func TestDCacheService(t *testing.T) {
	ser := rpc.DCacheService{}
	ser.Get(context.Background(), nil)
	ser.Set(context.Background(), nil)
	ser.Del(context.Background(), nil)
	ser.Register(context.Background(), nil)
	ser.Unregister(context.Background(), nil)
	ser.SyncSrvGroups(context.Background(), nil)
}

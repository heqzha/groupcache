package test

import (
	"testing"

	"github.com/heqzha/dcache/utils"
)

func TestRPCClient(t *testing.T) {
	pool := utils.GetCliPoolInst()
	cli, err := pool.Add("127.0.0.1:11000")
	if err != nil {
		t.Error(err)
		return
	}
	res, err := cli.Register("test1", "127.0.0.1:11001")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("cli.Register: %t", res.Status)
}

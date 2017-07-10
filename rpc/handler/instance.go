package handler

import (
	"sync"

	"github.com/heqzha/dcache/conf"
	"github.com/heqzha/dcache/core"
)

var sgmInst *core.SGM
var once sync.Once

func GetSGMInst() *core.SGM {
	once.Do(func() {
		if conf.HostPort == "" {
			panic("missing host_port in config.yaml")
		}
		sgmInst.Init(conf.HostPort)
	})
	return sgmInst
}

package instance

import (
	"sync"

	"github.com/heqzha/dcache/core"
)

// Config instance
var confInst *core.Config
var confInstOnce sync.Once

func GetConfInst() *core.Config {
	confInstOnce.Do(func() {
		confInst.Init()
	})
	return confInst
}

// SGM instance
var sgmInst *core.SGM
var sgmInstOnce sync.Once

func GetSGMInst() *core.SGM {
	conf := GetConfInst()
	sgmInstOnce.Do(func() {
		if conf.Addr == "" {
			panic("missing addr in config.yaml")
		}
		sgmInst.Init(conf.Addr)
	})
	return sgmInst
}

// SGHash instance
var sghInst *core.SGHash
var sghInstOnce sync.Once

func GetSGHInst() *core.SGHash {
	sghInstOnce.Do(func() {
		sghInst.Init(3, nil)
	})
	return sghInst
}

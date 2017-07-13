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
		confInst = new(core.Config)
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
		sgmInst = new(core.SGM)
		sgmInst.Init(conf.Addr)
	})
	return sgmInst
}

// SGHash instance
var sghInst *core.SGHash
var sghInstOnce sync.Once

func GetSGHInst() *core.SGHash {
	sghInstOnce.Do(func() {
		sghInst = new(core.SGHash)
		sghInst.Init(3, nil)
	})
	return sghInst
}

//Message queue instance
var msgQInst *core.MessageQueue
var msgQInstOnce sync.Once

func GetMsgQInst() *core.MessageQueue {
	msgQInstOnce.Do(func() {
		msgQInst = new(core.MessageQueue)
		msgQInst.Init()
	})
	return msgQInst
}

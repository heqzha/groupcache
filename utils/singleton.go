package utils

import (
	"sync"

	"github.com/heqzha/dcache/core"
	"github.com/heqzha/dcache/rpc"
)

// Config instance
var confInst *Config
var confInstOnce sync.Once

func GetConfInst() *Config {
	confInstOnce.Do(func() {
		confInst = new(Config)
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
			panic("missing addr in config.yml")
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

//RPC client pool instance
var cliPoolInst *rpc.CSClientPool
var cliPoolInstOnce sync.Once

func GetCliPoolInst() *rpc.CSClientPool {
	cliPoolInstOnce.Do(func() {
		cliPoolInst = new(rpc.CSClientPool)
	})
	return cliPoolInst
}

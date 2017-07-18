package handler

import (
	"fmt"

	"github.com/heqzha/dcache/core"
	"github.com/heqzha/dcache/utils"
	"github.com/heqzha/goutils/logger"
)

var (
	sgm  = utils.GetSGMInst()
	msgQ = utils.GetMsgQInst()
)

func Register(group, addr string) error {
	sgm.Register(group, addr)
	msgQ.Push("srvgroup", map[string]interface{}{
		"type": "sync",
	})
	return nil
}

func Unregister(group, addr string) error {
	sgm.Unregister(group, addr)
	msgQ.Push("srvgroup", map[string]interface{}{
		"type": "sync",
	})
	return nil
}

func SyncSrvGroups(srvgroups []byte) (core.Condition, *core.SGM, error) {
	tmpSGM := core.SGM{}
	tmpSGM.Load(srvgroups)
	logger.Info("SGM.Merge", fmt.Sprintf("Before Merge: %s", sgm.CompareReadable(tmpSGM)))
	cond := sgm.Merge(tmpSGM)
	logger.Info("SGM.Merge", fmt.Sprintf("After Merge: %s", sgm.CompareReadable(tmpSGM)))
	return cond, sgm, nil
}

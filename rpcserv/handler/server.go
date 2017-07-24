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
	logger.Debug("SGM.Register", fmt.Sprintf("%s: %s", group, addr))
	sgm.Register(group, addr)
	msgQ.Push("srvgroup", map[string]interface{}{
		"type": "sync",
	})
	return nil
}

func Unregister(group, addr string) error {
	logger.Debug("SGM.Unregister", fmt.Sprintf("%s: %s", group, addr))
	sgm.Unregister(group, addr)
	msgQ.Push("srvgroup", map[string]interface{}{
		"type": "sync",
	})
	return nil
}

func SyncSrvGroups(srvgroups []byte) (core.Condition, []byte, error) {
	tmpSGM := core.SGM{}
	tmpSGM.Load(srvgroups)
	logger.Debug("SGM.Merge", fmt.Sprintf("Before Merge: %s", sgm.CompareReadable(tmpSGM)))
	cond := sgm.Merge(tmpSGM)
	logger.Debug("SGM.Merge", fmt.Sprintf("After Merge: %s", sgm.CompareReadable(tmpSGM)))
	dump, err := sgm.Dump()
	if err != nil {
		return -1, dump, err
	}
	return cond, dump, nil
}

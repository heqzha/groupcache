package process

import (
	"fmt"
	"time"

	"github.com/heqzha/dcache/core"
	"github.com/heqzha/dcache/instance"
	"github.com/heqzha/goutils/flow"
)

var sgMsgQ *core.MessageQueue
var (
	fh = flow.FlowNewHandler()
)

func MaintainSvrGroups() error {
	sgMsgQ = instance.GetMsgQInst()
	l, err := fh.NewLine(Receive, Handle, Reload)
	if err != nil {
		return err
	}
	fh.Start(l, flow.Params{})
	return nil
}

func Receive(c *flow.Context) {
	//TODO Receive event from channel and do some preparations
	for {
		msg := sgMsgQ.Pop("srvgroup")
		if msg == nil {
			time.Sleep(500 * time.Millisecond)
			continue
		}
		c.Set("msg", msg)
		c.Pass()
	}
}

func Handle(c *flow.Context) {
	//TODO Handle event concurrently
	msg := c.MustGet("msg").(map[string]interface{})
	c.Set("type", msg["type"])
	switch msg["type"].(string) {
	case "sync":
		fmt.Println("Handle: sync")
		c.Next()
	case "ping":
		fmt.Println("Handle: ping")
	}
}

func Reload(c *flow.Context) {
	//TODO reload sghash
	fmt.Println("Reload", c.MustGet("type"))
}

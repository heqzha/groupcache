package test

import (
	"fmt"
	"testing"

	"github.com/heqzha/dcache/instance"
)

func TestMessageQueue(t *testing.T) {
	q := instance.GetMsgQInst()
	for index := 0; index < 10; index++ {
		q.Push("test1", index)
	}

	for (*q)["test1"].Len() != 0 {
		fmt.Println(q.Pop("test1"))
	}
}

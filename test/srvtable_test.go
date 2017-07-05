package test

import (
	"testing"

	"fmt"

	"strconv"

	"github.com/heqzha/dcache/core"
)

func TestSrvTable(t *testing.T) {
	tb1 := core.SrvTable{
		"127.0.0.1:1001": true,
		"127.0.0.1:1002": true,
		"127.0.0.1:1003": true,
	}
	tb2 := core.SrvTable{
		"127.0.0.1:1004": true,
		"127.0.0.1:1005": true,
	}
	tb1.Add(tb2)
	fmt.Println(tb1.String())

	tb2 = tb1.Clone()
	fmt.Println(tb2.String())
}

func TestSTM(t *testing.T) {
	var (
		stm  core.STM
		stm2 core.STM
	)
	stm.Init("127.0.0.1:1001")
	stm2.Init("127.0.0.2:1001")

	for index := 1001; index < 1004; index++ {
		stm.Register("127.0.0.1:" + strconv.Itoa(index))
	}
	tb := stm.GetTable()
	clk := stm.GetClock()
	fmt.Println(tb.String())
	fmt.Println(clk.ReturnVCString())

	for index := 1001; index < 1003; index++ {
		stm2.Register("127.0.0.2:" + strconv.Itoa(index))
	}
	tb = stm2.GetTable()
	clk = stm2.GetClock()
	fmt.Println(tb.String())
	fmt.Println(clk.ReturnVCString())

	cond := stm2.Merge(stm)
	fmt.Println("Condition:", cond)
	tb = stm2.GetTable()
	clk = stm2.GetClock()
	fmt.Println(tb.String())
	fmt.Println(clk.ReturnVCString())

	cond = stm.Merge(stm2)
	fmt.Println("Condition:", cond)
	tb = stm.GetTable()
	clk = stm.GetClock()
	fmt.Println(tb.String())
	fmt.Println(clk.ReturnVCString())

	fmt.Println(stm.CompareReadable(stm2))
}

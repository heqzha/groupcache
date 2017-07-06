package test

import (
	"math/rand"
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

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func TestSTHash(t *testing.T) {
	var (
		sth core.STHash
		stm core.STM
	)
	stm.Init("127.0.0.1:1001")

	for index := 1001; index < 1010; index++ {
		stm.Register("127.0.0.1:" + strconv.Itoa(index))
	}

	sth.Init(3, nil)
	sth.Load(stm.GetTable())

	for index := 0; index < 20; index++ {
		key := RandStringBytes(10)
		srv := sth.Pick(key)
		fmt.Println(key, srv)
	}
}

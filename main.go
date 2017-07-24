package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/heqzha/dcache/process"
	"github.com/heqzha/dcache/rpcserv"
	"github.com/heqzha/dcache/utils"

	"github.com/heqzha/goutils/logger"
)

var (
	conf *utils.Config
)

func init() {
	conf = utils.GetConfInst()
	logger.Config(conf.LogDir, logger.LOG_LEVEL_DEBUG)
}

func CreatePID(name string) int {
	wd, _ := os.Getwd()
	pidFile, err := os.OpenFile(filepath.Join(wd, fmt.Sprintf("%s.pid", name)), os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("failed to create pid file: %s", err.Error())
		os.Exit(1)
	}
	defer pidFile.Close()

	pid := os.Getpid()
	pidFile.WriteString(strconv.Itoa(pid))
	return pid
}

func main() {
	pid := CreatePID("dcache")
	process.MaintainSvrGroups()
	defer process.StopAll()
	fmt.Printf("Start to Serving :%d with pid %d\n", conf.ServPort, pid)
	rpcserv.Run(conf.ServPort)
}

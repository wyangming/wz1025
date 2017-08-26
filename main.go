package main

import (
	"fmt"
	"runtime"
	"wz1025/http"

	"sync"
	"wz1025/zzdemo"

	"github.com/jakecoffman/cron"
)

func main() {
	//cpu密集型项目时充分利用cpu性能
	runtime.GOMAXPROCS(runtime.NumCPU())
	var wg sync.WaitGroup
	wg.Add(1)
	//http服务器
	http.Init()
	go zzdemo.HttpDemo()
	wg.Wait()
}

func conDemo() {
	spec := "*/5 * * * * ?" //每5s执行一次
	cronJob := cron.New()
	cronJob.AddFunc(spec, conFun, "cronFun")
	cronJob.Start()
}

func conFun() {
	fmt.Println("this is conFun Test")
}

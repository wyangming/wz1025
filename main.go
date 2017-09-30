package main

import (
	"fmt"
	"runtime"
	"sync"
	_"wz1025/module/http"
	"wz1025/zzdemo"

	"github.com/jakecoffman/cron"
)

type Item struct {
	Foo string
	Bar string
}

func main() {
	run()
}

func run() {
	//cpu密集型项目时充分利用cpu性能
	runtime.GOMAXPROCS(runtime.NumCPU())
	var wg sync.WaitGroup
	wg.Add(1)
	go zzdemo.HttpDemo()
	//go conDemo()
	wg.Wait()
}

func conDemo() {
	spec := "*/5 * * * * ?" //每5s执行一次
	cronJob := cron.New()
	//cronJob.RemoveJob()//要删除任务使用这个方法
	cronJob.AddFunc(spec, conFun, "cronFun")
	cronJob.Start()
}

func conFun() {
	fmt.Println("this is conFun Test")
}

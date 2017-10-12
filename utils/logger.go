package utils

import (
	"fmt"
	"time"
	"github.com/astaxie/beego"
)

//错误日志打印
func ErrorLog(info string, err error) {
	fmt.Println("time is ", time.Now())
	fmt.Println(info)
	beego.Error(err)
}

//信息日志打印
func InfoLog(format string, args ...interface{}) {
	fmt.Println("time is ", time.Now())
	if len(args) > 0 {
		fmt.Println(fmt.Sprintf(format, args))
	} else {
		fmt.Println(format)
	}
}

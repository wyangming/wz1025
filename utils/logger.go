package utils

import (
	"fmt"
	"time"
)

//错误日志打印
func ErrorLog(info string, obj interface{}) {
	fmt.Println("time is ", time.Now())
	fmt.Println(info,obj)
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

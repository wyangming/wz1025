package utils

import (
	"fmt"

	"github.com/astaxie/beego"
)

//错误日志打印
func ErrorLog(info string, err error) {
	fmt.Println(info)
	beego.Error(err)
}

//信息日志打印
func InfoLog(format string, args ...interface{}) {
	if len(args) > 0 {
		fmt.Println(fmt.Sprintf(format, args))
	} else {
		fmt.Println(format)
	}
}

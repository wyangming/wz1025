// init.go
package http

import (
	"wz1025/http/controllers"

	"github.com/astaxie/beego"
)

func Init() {
	go run()
}
func run() {
	//路由设置
	init_Router()
	//http服务器启动
	beego.Run()
}
func init_Router() {
	mainController := &controllers.MainController{}
	beego.Router("/", mainController, "*:Get")
	beego.Router("/video", mainController, "*:Video")
}

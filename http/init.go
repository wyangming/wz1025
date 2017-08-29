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
	//首页
	beego.Router("/", mainController, "*:Get")
	//视频页面
	beego.Router("/video", mainController, "*:Video")
	//登录页面
	beego.Router("/login", mainController, "*:Login")
}

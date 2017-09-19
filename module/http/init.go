// init.go
package http

import (
	"wz1025/module/http/controllers"
	"wz1025/module/http/define"
	"wz1025/module/http/filter"

	"github.com/astaxie/beego"
)

func Init() {
	go run()
}
func run() {

	init_filter()

	init_Router()
	//http服务器启动
	beego.Run()
}

//过滤器设置
func init_filter() {
	//会员过滤器设置
	beego.InsertFilter(define.FILTER_MEMBER, beego.BeforeRouter, filter.MemberAuthor)
	beego.InsertFilter(define.FILTER_MEMBER_ALL, beego.BeforeRouter, filter.MemberAuthor)
}

//路由设置
func init_Router() {
	//主接口
	mainController := &controllers.MainController{}
	//首页
	beego.Router(define.URL_INDEX, mainController, "*:Get")
	//会员登录页面
	beego.Router(define.URL_LOGIN, mainController, "get,post:Login")
	//会员注册页面
	beego.Router(define.URL_REG, mainController, "get,Post:Reg")

	//会员
	memberController := &controllers.MemberController{}
	//会员主页面
	beego.Router(define.URL_MEMBER, memberController, "*:Get")
	//会员信息
	beego.Router(define.URL_MEMBER_INFO, memberController, "*:Info")
	//会员视频页面
	beego.Router(define.URL_MEMBER_VIDEO, memberController, "*:Video")
	//得到视频解析信息
	beego.Router(define.URL_MEMBER_EXPLAIN_INFO, memberController, "*:AjaxExplainInfo")
}

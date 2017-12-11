// init.go
package http

import (
	"wz1025/module/http/controllers"
	"wz1025/module/http/define"
	"wz1025/module/http/filter"

	"github.com/astaxie/beego"
)

//初始化http内容
func init() {
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

	//admin过滤器
	beego.InsertFilter(define.FILTER_ADMIN, beego.BeforeRouter, filter.AdminAuthor)
	beego.InsertFilter(define.FILTER_ADMIN_ALL, beego.BeforeRouter, filter.AdminAuthor)
}

//路由设置
func init_Router() {
	//主接口
	mainController := &controllers.MainController{}
	//首页
	beego.Router(define.URL_INDEX, mainController, define.CON_MAIN_GET_METHOD)
	//会员登录页面
	beego.Router(define.URL_LOGIN, mainController, define.CON_MAIN_LOGIN_METHOD)
	//会员注册页面
	beego.Router(define.URL_REG, mainController, define.CON_MAIN_REG_METHOD)
	//admin登录页面
	beego.Router(define.URL_ADMINLOGIN, mainController, define.CON_MAIN_ADMINLOGIN_METHOD)

	//会员
	memberController := &controllers.MemberController{}
	//会员主页面
	beego.Router(define.URL_MEMBER, memberController, define.CON_MEMBER_GET_METHOD)
	//会员信息
	beego.Router(define.URL_MEMBER_INFO, memberController, define.CON_MEMBER_INFO_METHOD)
	//会员视频页面
	beego.Router(define.URL_MEMBER_VIDEO, memberController, define.CON_MEMBER_VIDEO_METHOD)
	//得到视频解析信息
	beego.Router(define.URL_MEMBER_EXPLAIN_INFO, memberController, define.CON_MEMBER_AJAXEXPLAININFO_METHOD)

	//管理员
	adminController := &controllers.AdminController{}
	//管理员主页
	beego.Router(define.URL_ADMIN, adminController, define.CON_ADMIN_GET_METHOD)
	//解析地址
	beego.Router(define.URL_ADMIN_EXPLAIN, adminController, define.CON_ADMIN_EXPLAIN_METHOD)
	//解析地址列表信息
	beego.Router(define.URL_ADMIN_EXPLAIN_LIST, adminController, define.CON_ADMIN_EXPLAINLIST_METHOD)
	//修改解析地址状态
	beego.Router(define.URL_ADMIN_EXPLAIN_ACTIVEUPDATE, adminController, define.CON_ADMIN_EXPLAINACTIVEUPDATE_METHOD)
	//会员管理
	beego.Router(define.URL_ADMIN_MEMBER, adminController, define.CON_ADMIN_MEMBER_METHOD)
	//会员列表信息
	beego.Router(define.URL_ADMIN_MEMBER_LIST, adminController, define.CON_ADMIN_MEMBERlIST_METHOD)
	//修改会员状态
	beego.Router(define.URL_ADMIN_MEMBER_ACTIVEUPDATE, adminController, define.CON_ADMIN_MEMBERACTIVEUPDATE_METHOD)
	//修改会员过期时间
	beego.Router(define.URL_ADMIN_MEMBER_EXPIREUPDATE, adminController, define.CON_ADMIN_MEMBEREXPIREUPDATE_METHOD)
	//手动更新解析地址
	beego.Router(define.URL_ADMIN_EXPLAIN_SPIDERUPDATE, adminController, define.CON_ADMIN_EXPLAINSPIDERUPDATE_METHOD)

	//微信
	wxController := &controllers.WxController{}
	//微信接口
	beego.Router(define.URL_WECHAT, wxController)
}

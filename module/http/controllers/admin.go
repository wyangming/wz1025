package controllers

import (
	"github.com/astaxie/beego"
	"wz1025/module/http/define"
	db "wz1025/db/http"
	"fmt"
)

type AdminController struct {
	beego.Controller
}

//在方法调用之前，beego会自动调用
//主要把session里的数据放入Controller
func (self *AdminController) Prepare() {
	//当前用户信息
	member_info, member_info_has := self.GetSession(define.SESSION_ADMIN_INFO).(map[string]string)
	if member_info_has {
		self.Data[define.SESSION_ADMIN_INFO] = member_info
	}
}

//首页面
func (self *AdminController) Get() {
	self.TplName = define.CON_ADMIN_MAIN_PAGE
}

//解析地址
func (self *AdminController) Explain() {
	self.TplName = define.CON_ADMIN_EXPLAIN
}

//解析地址信息
func (self *AdminController) Explain_List() {
	//参数处理
	self.Ctx.Request.ParseForm()
	form_values := self.Ctx.Request.Form
	fmt.Println("前台参数信息为：")
	//limit, page
	for param_key, param_val := range form_values {
		fmt.Println(param_key, param_val)
	}
	args := map[string]interface{}{
		"startLine": 0,
		"endLine":   20,
		"active":    1,
	}

	//数据
	explains := db.NewAdminDbFun().FindVideoExplain(args)
	//设置必要参数
	explains["code"] = 0
	explains["msg"] = ""

	self.Data["json"] = explains
	self.ServeJSON()
}

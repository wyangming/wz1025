package controllers

import (
	"github.com/astaxie/beego"
	"wz1025/module/http/define"
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

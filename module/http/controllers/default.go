package controllers

import (
	"fmt"
	"strings"
	"wz1025/db"
	"wz1025/module/http/define"
	//"wz1025/utils"

	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

//首页
func (self *MainController) Get() {
	self.TplName = define.CON_MAIN_INDEX_PAGE
	fmt.Println("index method")
}

//登录页面
func (self *MainController) Login() {
	//非post请求直接跳转到登录页面不做处理
	if "POST" != self.Ctx.Request.Method {
		self.TplName = define.CON_MAIN_LOGIN_PAGE
		return
	}

	//登录逻辑
	form_values := self.Ctx.Request.Form
	args := make(map[string]interface{}, len(form_values))
	for k, v := range form_values {
		if len(v) > 0 {
			args[k] = strings.TrimSpace(v[0])
		}
	}

	user_info := db.NewHttpDbFun().Login(args)
	if user_info == nil {
		self.Data[define.CON_MAIN_LOGIN_STATUS] = "false"
		self.TplName = define.CON_MAIN_LOGIN_PAGE
		return
	}

	//存入session
	self.SetSession(define.SESSION_MEMBER_INFO, user_info)

	self.Data[define.CON_MAIN_LOGIN_STATUS] = "true"
	self.TplName = define.CON_MEMBER_MAIN_PAGE
}

//注册
func (self *MainController) Reg() {
	//非post请求直接跳转到登录页面不做处理
	if "POST" != self.Ctx.Request.Method {
		self.TplName = define.CON_MAIN_REG_PAGE
		return
	}

	//注册逻辑
	form_values := self.Ctx.Request.Form
	args := make(map[string]interface{}, len(form_values))
	for k, v := range form_values {
		//v是数组类型判断v是否有值
		if len(v) > 0 {
			args[k] = strings.TrimSpace(v[0])
		}
	}

	if db.NewHttpDbFun().RegMember(args) {
		self.Data[define.CON_MAIN_REG_STATUS] = "true"
		self.TplName = define.CON_MAIN_LOGIN_PAGE
		return
	}
	self.Data[define.CON_MAIN_REG_STATUS] = "false"
	self.TplName = define.CON_MAIN_REG_PAGE
}

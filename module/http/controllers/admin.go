package controllers

import (
	"github.com/astaxie/beego"
	"wz1025/module/http/define"
	db "wz1025/db/http"
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
	explains := map[string]interface{}{}
	explains["code"] = 1

	//参数处理
	limit, limit_err := self.GetInt("limit", 10)
	if limit_err != nil {
		explains["msg"] = "分页参数不正确"
		self.Data["json"] = explains
		self.ServeJSON()
		return
	}
	page, page_err := self.GetInt("page", 1)
	if page_err != nil {
		explains["msg"] = "分页参数不正确"
		self.Data["json"] = explains
		self.ServeJSON()
		return
	}

	args := map[string]interface{}{
		"startLine": limit * (page - 1),
		"endLine":   limit * page,
		"active":    1,
	}

	//数据
	db.NewAdminDbFun().FindVideoExplain(args, explains)
	//设置必要参数
	explains["code"] = 0
	explains["msg"] = ""

	self.Data["json"] = explains
	self.ServeJSON()
}

//添加解析地址
func (self *AdminController) Explain_Add() {
	explains := map[string]interface{}{}
	explains["code"] = 1

	//参数处理
	url_addr := self.GetString("url_addr")
	if url_addr == "" {
		explains["msg"] = "参数不正确"
		self.Data["json"] = explains
		self.ServeJSON()
		return
	}
	type_int, type_err := self.GetInt("type", 0)
	if type_err != nil {
		explains["msg"] = "参数不正确"
		self.Data["json"] = explains
		self.ServeJSON()
		return
	}

	args := map[string]interface{}{
		"url_addr": url_addr,
		"type":     type_int,
	}

	if !db.NewAdminDbFun().AddVideoExplainUrl(args) {
		explains["msg"] = "添加失败"
		self.Data["json"] = explains
		self.ServeJSON()
		return
	}

	explains["code"] = 0
	explains["msg"] = ""
	self.Data["json"] = explains
	self.ServeJSONP()
}

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

//解析地址管理
func (self *AdminController) Explain() {
	self.TplName = define.CON_ADMIN_EXPLAIN_PAGE
}

//解析地址信息
func (self *AdminController) Explain_List() {
	explains := map[string]interface{}{}
	explains["code"] = 1

	//分页参数处理
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
	}

	//默认查询所有
	type_val, _ := self.GetUint8("type", 255)
	//默认查询激活的
	active_val, _ := self.GetUint8("active", 1)
	if type_val < 255 {
		args["type"] = type_val
	}
	if active_val < 255 {
		args["active"] = active_val
	}

	//数据
	db.NewAdminDbFun().FindVideoExplain(args, explains)
	//设置必要参数
	explains["code"] = 0
	explains["msg"] = ""

	self.Data["json"] = explains
	self.ServeJSON()
}

//修改解析地址状态
func (self *AdminController) Explain_ActiveUpdate() {
	explains := map[string]interface{}{}
	explains["code"] = 1

	//参数处理
	ids := self.GetString("ids")
	if ids == "" {
		explains["msg"] = "参数不正确"
		self.Data["json"] = explains
		self.ServeJSON()
		return
	}
	action_type := self.GetString("action_type")
	if action_type == "" {
		explains["msg"] = "参数不正确"
		self.Data["json"] = explains
		self.ServeJSON()
		return
	}

	if !db.NewAdminDbFun().UpdateVideoExplainActive(ids, action_type) {
		explains["msg"] = "修改状态失败"
	} else {
		explains["code"] = 0
		explains["msg"] = ""
	}
	self.Data["json"] = explains
	self.ServeJSON()
}

//会员管理
func (self *AdminController) Member() {
	self.TplName = define.CON_ADMIN_MEMBER_PAGE
}

//会员信息
func (self *AdminController) Member_List() {
	explains := map[string]interface{}{}
	explains["code"] = 1

	//分页参数处理
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
	}

	//默认查询正常的
	active_val, _ := self.GetUint8("active", 1)
	if active_val < 255 {
		args["active"] = active_val
	}
	phone_num := self.GetString("phone_num", "")
	if phone_num != "" {
		args["phone_num"] = phone_num
	}
	nick_name := self.GetString("nick_name", "")
	if nick_name != "" {
		args["nick_name"] = nick_name
	}

	//数据
	db.NewAdminDbFun().FindMember(args, explains)
	//设置必要参数
	explains["code"] = 0
	explains["msg"] = ""

	self.Data["json"] = explains
	self.ServeJSON()
}

//修改会员状态
func (self *AdminController) Member_ActiveUpdate() {
	explains := map[string]interface{}{}
	explains["code"] = 1

	//参数处理
	ids := self.GetString("ids")
	if ids == "" {
		explains["msg"] = "参数不正确"
		self.Data["json"] = explains
		self.ServeJSON()
		return
	}
	action_type := self.GetString("action_type")
	if action_type == "" {
		explains["msg"] = "参数不正确"
		self.Data["json"] = explains
		self.ServeJSON()
		return
	}

	if !db.NewAdminDbFun().UpdateMemberActive(ids, action_type) {
		explains["msg"] = "修改状态失败"
	} else {
		explains["code"] = 0
		explains["msg"] = ""
	}
	self.Data["json"] = explains
	self.ServeJSON()
}

//修改会员过期时间
func (self *AdminController) Member_ExpireUpdate() {
	explains := map[string]interface{}{}
	explains["code"] = 1

	//参数处理
	field := self.GetString("field", "")
	if field == "" {
		explains["msg"] = "参数不正确"
		self.Data["json"] = explains
		self.ServeJSON()
		return
	}
	id_val, id_err := self.GetUint64("id", 0)
	if id_err != nil {
		explains["msg"] = "参数不正确"
		self.Data["json"] = explains
		self.ServeJSON()
		return
	}
	date_val := self.GetString("value", "")
	if date_val == "" {
		explains["msg"] = "参数不正确"
		self.Data["json"] = explains
		self.ServeJSON()
		return
	}

	args := map[string]interface{}{
		"field": field,
		"id":    id_val,
		"value": date_val,
	}

	if(!db.NewAdminDbFun().UpdateMemberExpire(args)){
		explains["msg"] = "延长时间失败"
	} else {
		explains["code"] = 0
		explains["msg"] = ""
	}

	self.Data["json"] = explains
	self.ServeJSON()
}

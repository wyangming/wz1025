package controllers

import (
	"github.com/astaxie/beego"
	"wz1025/module/http/define"
)

type MemberController struct {
	beego.Controller
}

//在方法调用之前，beego会自动调用
//主要把session里的数据放入Controller
func (self *MemberController) Prepare() {
	//当前用户信息
	member_info, member_info_has := self.GetSession(define.SESSION_MEMBER_INFO).(map[string]interface{})
	if member_info_has {
		self.Data[define.SESSION_MEMBER_INFO] = member_info
	}
}

//会员首页
func (self *MemberController) Get() {
	self.TplName = define.CON_MEMBER_MAIN_PAGE
}

//会员信息
func (self *MemberController) Info() {
	self.TplName = define.CON_MEMBER_INFO_PAGE
}

//视频页面
func (self *MemberController) Video() {
	self.TplName = define.CON_MEMBER_VIDEO_PAGE

}

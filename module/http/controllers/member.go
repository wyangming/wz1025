package controllers

import (
	"github.com/astaxie/beego"
	"strings"
	"time"
	db "wz1025/db/http"
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

	//参数是否非法
	type_val := self.GetString("type", "0")
	if "0" == type_val {
		self.Data["status"] = 1
		return
	}
	self.Data["type"] = type_val

	//是否过期
	if !self.isExpire(type_val) {
		self.Data["status"] = 2
		return
	}

	self.Data["status"] = 0
}

//判断用户vip是否过期
func (self *MemberController) isExpire(type_int_str string) bool {
	//过期时间
	type_str := "aiqiyi_expire"
	switch type_int_str {
	case "2":
		type_str = "youku_expire"
	case "3":
		type_str = "letv_expire"
	case "4":
		type_str = "tentcent_expire"
	default:
	}
	member_info_obj := self.Data[define.SESSION_MEMBER_INFO]
	member_info, _ := member_info_obj.(map[string]interface{})

	//判断过期
	expire_val, _ := member_info[type_str]
	expire_time, ok := expire_val.(time.Time)
	if !ok {
		return false
	}

	now_time := time.Now()
	return now_time.Before(expire_time)
}

//请求视频解析数据
func (self *MemberController) AjaxExplainInfo() {
	res := map[string]string{
		"result": "false",
		"msg":    "has error this request",
	}
	self.Data["json"] = res

	//验证信息
	//url
	url_val := self.GetString("url", "")
	if url_val == "" {
		res["msg"] = "播放地址不可以为空"
		self.ServeJSON()
		return
	}
	//视频类型
	type_val, type_has := self.GetUint8("type", 0)
	if type_has != nil {
		res["msg"] = "非法参数"
		self.ServeJSON()
		return
	}
	//是否过期
	if !self.isExpire(string(type_val)) {
		res["msg"] = "会员过期"
		self.ServeJSON()
		return
	}

	//获得解析地址
	explainUrl := db.NewMemberDbFun().VideoExplainUrlByType(type_val)
	if len(explainUrl) < 1 {
		res["msg"] = "地址解析失败"
		self.ServeJSON()
		return
	}

	//拼装返回信息并返回
	res["result"] = "true"
	res["msg"] = "request is success"
	res["info"] = strings.Join([]string{"<iframe src='", explainUrl, url_val, "' id='player' name='player' width='100%' height='600px' allowtransparency='true' frameborder='0' scrolling='no'></iframe>"}, "")
	self.ServeJSON()
}

func (self *MemberController) ToSearchVideoRec() {
	self.TplName = define.CON_MEMBER_SEARCHVIDEOREC

	video_name := self.GetString("video_name")
	if len(video_name) < 1 {
		self.Data["first"] = "true"
		return
	}
	self.Data["first"] = "false"
	self.Data["video_name"] = video_name

	ret := db.NewMemberDbFun().FindSpiderVideoRecs(video_name)
	self.Data["Recs"] = *ret
}

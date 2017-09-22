package filter

import (
	"github.com/astaxie/beego/context"
	"wz1025/module/http/define"
)

//会员后台判断是否登录过滤器
func MemberAuthor(ctx *context.Context) {
	_, member_info_has := ctx.Input.Session(define.SESSION_MEMBER_INFO).(map[string]interface{})
	if !member_info_has {
		ctx.Redirect(302, define.URL_LOGIN)
	}
}

//admin后台判断是否登录过滤器
func AdminAuthor(ctx *context.Context) {
	_, admin_info_has := ctx.Input.Session(define.SESSION_ADMIN_INFO).(map[string]string)
	if !admin_info_has {
		ctx.Redirect(302, define.URL_ADMINLOGIN)
	}
}

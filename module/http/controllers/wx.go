package controllers

import (
	"github.com/astaxie/beego"
	"wz1025/utils/wx"
)

type WxController struct {
	beego.Controller
	wxToken string
}

func (this *WxController) getWxToken() string {
	if len(this.wxToken) < 1 {
		this.wxToken = beego.AppConfig.String("wxToken")
	}
	return this.wxToken
}

//消息接收
func (this *WxController) Post() {
	var str_msg string = string(this.Ctx.Input.RequestBody)
	recmsg := wx.InitRecMsg(str_msg)
	if recmsg.MsgType == wx.MsgTypeText {
		this.Ctx.WriteString(wx.ReplyText("不支持关键词回复", recmsg))
	} else {
		this.Ctx.WriteString(wx.ReplyText("不支持此类消息", recmsg))
	}
}

//微信认证
func (this *WxController) Get() {
	timestamp, nonce, signatureIn, echostr := this.GetString("timestamp"), this.GetString("nonce"), this.GetString("signature"), this.GetString("echostr")
	this.Ctx.WriteString(wx.WxAuth(timestamp, nonce, signatureIn, echostr, this.getWxToken()))
}

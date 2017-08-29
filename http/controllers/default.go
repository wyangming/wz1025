package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

//首页
func (this *MainController) Get() {
	this.TplName = "index.html"
}

//视频页面
func (this *MainController) Video() {
	this.TplName = "video.html"
}

//登录页面
func (this *MainController) Login() {
	this.TplName = "login.html"
}

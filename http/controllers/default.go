package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (this *MainController) Get() {
	this.TplName = "index.html"
}
func (this *MainController) Video() {
	this.TplName = "video.html"
}

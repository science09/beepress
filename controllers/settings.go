package controllers

import (
	"beepress/models"
	"github.com/astaxie/beego"
)

type Settings struct {
	BaseController
}

func (this *Settings) NestPrepare() {
	this.requireAdmin()
}

func (this *Settings) Index() {
	flash := beego.NewFlash()
	settings, err := models.GetSettings()
	if err != nil {
		flash.Error("数据库获取setting失败!")
		flash.Store(&this.Controller)
		return
	}
	this.Data["title"] = "设置项"
	this.Data["settings"] = settings
	this.TplName = "settings/index.html"
}

func (this *Settings) Edit() {
	key := this.Ctx.Input.Param(":key")
	setting := models.FindSettingByKey(key)
	this.Data["setting"] = setting
	this.Data["title"] = "修改设置"
	this.TplName = "settings/edit.html"
}

func (this *Settings) Update() {
	flash := beego.NewFlash()
	key := this.Ctx.Input.Param(":key")
	setting := models.FindSettingByKey(key)
	setting.Val = this.GetString("val")
	this.Data["setting"] = setting
	if err := models.UpdateSetting(setting); err != nil {
		flash.Error("更新失败!")
		flash.Store(&this.Controller)
		this.Redirect("/settings/edit.html")
	}
	flash.Success("设置更新成功")
	flash.Store(&this.Controller)
	this.Redirect("/settings")
}

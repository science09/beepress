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
	settings := []models.Setting{}
	//DB.Model(Setting{}).Order("`key` desc").Find(&settings)

	this.Data["title"] = "设置项"
	this.Data["settings"] = settings
	this.TplName = "settings/index.html"
}

func (this *Settings) Edit(key string) {
	setting := models.FindSettingByKey(key)
	this.Data["setting"] = setting
	this.Data["title"] = "修改设置"
	this.TplName = "settings/edit.html"
}

func (this *Settings) Update(key string) {
	setting := models.FindSettingByKey(key)
	//this.Params.Bind(&setting.Val, "val")
	this.Data["setting"] = setting
	if err := models.UpdateSetting(setting); err != nil {
		//	return this.Render("settings/edit.html")
		return
	}
	beego.NewFlash().Success("设置更新成功")
	this.Redirect("/settings")
}

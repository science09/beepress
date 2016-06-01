package controllers

import (
	"beepress/models"
	"github.com/astaxie/beego"
)

type Users struct {
	BaseController
	user models.User
}

func (this *Users) NestPrepare() {
	login := this.Ctx.Input.Param(":login")
	beego.Info("UserController", login)
	var err error
	this.user, err = models.FindUserByLogin(login)
	if err != nil {
		this.Abort("404")
	}
	this.Data["user"] = this.user
}

func (this *Users) Show() {
	//recentTopics := []models.Topic{}
	//DB.Order("id desc").Where("user_id = ?", this.user.Id).Limit(10).Find(&recentTopics)
	recentTopics, _ := models.GetTopicByUserId(this.user.Id)
	this.Data["title"] = "社区"
	this.Data["controller_name"] = "User"
	this.Data["recent_topics"] = recentTopics
	this.TplName = "users/show.html"
}

func (this *Users) Topics(login string) {
}

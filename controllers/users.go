package controllers

import (
	"github.com/science09/beepress/models"
)

type Users struct {
	BaseController
	user *models.User
}

func (this *Users) NestPrepare() {
	login := this.Ctx.Input.Param(":login")
	var err error
	this.user, err = models.FindUserByLogin(login)
	if err != nil {
		this.Abort("404")
	}
	this.Data["user"] = this.user
}

func (this *Users) Show() {
	recentTopics, _ := models.GetRecentTopics(this.user.Id)
	this.Data["title"] = "社区"
	this.Data["controller_name"] = "User"
	this.Data["recent_topics"] = recentTopics
	this.TplName = "users/show.html"
}

func (this *Users) Topics(login string) {
}

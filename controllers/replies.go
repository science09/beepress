package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/science09/beepress/models"
	"strconv"
)

type Replies struct {
	BaseController
	topic models.Topic
}

func (this *Replies) Create() {
	var err error
	flash := beego.NewFlash()
	this.requireUser()
	reply := &models.Reply{Body: this.GetString("body")}
	topicId, _ := strconv.Atoi(this.Ctx.Input.Param(":id"))
	this.topic, err = models.GetTopicById(int32(topicId))
	this.Data["topic"] = this.topic
	this.TplName = "topics/show.html"
	if err != nil {
		beego.Error("id:", topicId, err)
		flash.Error(err.Error())
		flash.Store(&this.Controller)
		return
	}

	reply.Topic = &this.topic
	reply.User = &this.currentUser
	err = models.CreateReply(reply)
	if err != nil {
		beego.Error(err)
		flash.Error(err.Error())
		flash.Store(&this.Controller)
		this.Redirect(fmt.Sprintf("/topics/%v", this.topic.Id))
	}
	replies, _ := models.GetReplyByTopicId(int32(topicId))
	this.Data["replies"] = replies
	this.Redirect(fmt.Sprintf("/topics/%v#reply%v", this.topic.Id, this.topic.RepliesCount))
}

func (this *Replies) Update() {
	flash := beego.NewFlash()
	this.requireUser()
	replyId, _ := strconv.Atoi(this.Ctx.Input.Param(":id"))
	reply, err := models.GetReplyById(int32(replyId))
	this.TplName = "replies/edit.html"
	if err != nil {
		flash.Error("")
		flash.Store(&this.Controller)
		return
	}
	if !this.isOwner(reply) {
		flash.Error("不允许修改他人的评论!")
		flash.Store(&this.Controller)
		return
	}
	reply.Body = this.GetString("body")
	err = models.UpdateReply(reply)
	if err != nil {
		flash.Error("评论修改失败!")
		flash.Store(&this.Controller)
		return
	}
	this.Redirect(fmt.Sprintf("/topics/%v", reply.Topic.Id))
}

//修改评论
func (this *Replies) Edit() {
	flash := beego.NewFlash()
	this.requireUser()
	replyId, _ := strconv.Atoi(this.Ctx.Input.Param(":id"))
	reply, err := models.GetReplyById(int32(replyId))
	this.Data["title"] = "修改回帖"
	this.Data["reply"] = reply
	this.TplName = "replies/edit.html"
	if err != nil {
		flash.Error("该评论已经不存在!")
		flash.Store(&this.Controller)
		return
	}
	if !this.isOwner(reply) {
		flash.Error("不允许修改他人的评论!")
		flash.Store(&this.Controller)
		return
	}
}

func (this *Replies) Delete() {
	flash := beego.NewFlash()
	this.requireUser()
	replyId, _ := strconv.Atoi(this.Ctx.Input.Param(":id"))
	reply, err := models.GetReplyById(int32(replyId))
	if err != nil {
		flash.Error("该条评论已经不存在!")
		flash.Store(&this.Controller)
		return
	}
	if !this.isOwner(reply) {
		flash.Error("不允许修改他人的评论!")
		flash.Store(&this.Controller)
		return
	}

	if err = reply.Del(); err != nil {
		flash.Error(err.Error())
	} else {
		flash.Success("回帖删除成功")
	}
	flash.Store(&this.Controller)
	this.Redirect(fmt.Sprintf("/topics/%v", reply.Topic.Id))
}

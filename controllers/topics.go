package controllers

import (
	"beepress/help"
	"beepress/models"
	"fmt"
	"github.com/astaxie/beego"
	"strconv"
	"strings"
)

type Topics struct {
	BaseController
}

func (this *Topics) NestPrepare() {

}

func (this *Topics) Index() {
	var page, nodeId int
	channel := ""

	//this.Ctx.Input.Bind(&page, "page")
	//this.Ctx.Input.Bind(&nodeId, "node_id")
	//beego.Info("page,nodeID:", this.GetString("page"), this.GetString("node_id"))

	if strings.EqualFold(channel, "node") {
		nodeId = 2
		node, _ := models.GetNodeById(int32(nodeId))
		this.Data["node"] = node
	}
	topics, pageInfo := models.FindTopicPages(channel, nodeId, page, 20)
	pageInfo.Path = "" //this.Request.URL.Path
	this.Data["title"] = "社区"
	this.Data["channel"] = channel
	this.Data["topics"] = topics
	this.Data["page_info"] = pageInfo
	this.TplName = "topics/index.html"
}

func (this *Topics) TopicNode() {
	var page int
	channel := "node"
	nodeId, _ := strconv.Atoi(this.Ctx.Input.Param(":node_id"))
	topics, pageInfo := models.FindTopicPages(channel, nodeId, page, 20)
	Node, _ := models.GetNodeById(int32(nodeId))
	pageInfo.Path = "" //this.Request.URL.Path
	this.Data["title"] = "社区"
	this.Data["node"] = Node
	this.Data["channel"] = channel
	this.Data["topics"] = topics
	this.Data["page_info"] = pageInfo
	this.TplName = "topics/index.html"
}
func (this *Topics) Popular() {
	var page int
	channel := "popular"
	topics, pageInfo := models.FindTopicPages(channel, 0, page, 20)
	pageInfo.Path = "" //this.Request.URL.Path
	this.Data["title"] = "社区"
	this.Data["channel"] = channel
	this.Data["topics"] = topics
	this.Data["page_info"] = pageInfo
	this.TplName = "topics/index.html"
}

func (this *Topics) Recent() {
	var page int
	channel := "recent"
	topics, pageInfo := models.FindTopicPages(channel, 0, page, 20)
	pageInfo.Path = "" //this.Request.URL.Path
	this.Data["title"] = "社区"
	this.Data["channel"] = channel
	this.Data["topics"] = topics
	this.Data["page_info"] = pageInfo
	this.TplName = "topics/index.html"
}

func (this *Topics) Feed() {
	topics, _ := models.FindTopicPages("recent", 0, 1, 20)
	this.Data["topics"] = topics
	this.Layout = "topics/feed.html"
	//this.TplName = "topics/feed.html"
	this.ServeXML()
}

func (this *Topics) New() {
	this.requireUser()
	t := &models.Topic{}
	this.Data["title"] = "发表新话题"
	this.Data["nodes"] = models.FindAllNodes()
	this.Data["topic"] = t
	this.TplName = "topics/new.html"
}

func (this *Topics) Create() {
	this.requireUser()
	var nodeId int32
	flash := beego.NewFlash()
	this.Ctx.Input.Bind(&nodeId, "node_id")
	beego.Info("Topic Create:", nodeId, this.GetString("node_id"))
	node, _ := models.GetNodeById(int32(nodeId))
	t := &models.Topic{
		Title: this.GetString("title"),
		Body:  this.GetString("body"),
		Node:  &node,
	}

	t.User = &this.currentUser
	err := models.CreateTopic(t)
	if err != nil {
		this.Data["topic"] = t
		this.Data["nodes"] = models.FindAllNodes()
		flash.Error(err.Error())
		flash.Store(&this.Controller)
		this.TplName = "topics/new.html"
		return
	}
	this.Redirect(fmt.Sprintf("/topics/%v", t.Id))
}

func (this *Topics) Show() {
	topicId, _ := strconv.Atoi(this.Ctx.Input.Param(":id"))
	topic, err := models.GetTopicById(int32(topicId))
	if err != nil {
		beego.Info(err)
		this.Abort("404")
	}
	replies, err := models.GetReplyByTopicId(topic.Id)
	if err != nil {
		beego.Info(err)
		this.Abort("404")
	}
	topic.RepliesCount = int32(len(replies))
	this.Data["topic"] = topic
	this.Data["title"] = topic.Title
	this.Data["replies"] = replies
	this.TplName = "topics/show.html"
}

func (this *Topics) Edit() {
	this.requireUser()
	t, _ := models.GetTopicById(help.StrToInt32(this.Ctx.Input.Param(":id")))
	if !this.isOwner(t) {
		beego.NewFlash().Error("没有修改的权限")
		this.Redirect("/")
	}
	this.Data["title"] = "修改话题"
	this.Data["topic"] = &t
	this.Data["nodes"] = models.FindAllNodes()
	this.TplName = "topics/edit.html"
}

func (this *Topics) Update() {
	this.requireUser()
	t, _ := models.GetTopicById(help.StrToInt32(this.Ctx.Input.Param(":id")))
	if !this.isOwner(t) {
		beego.NewFlash().Error("没有修改的权限")
		this.Redirect("/")
	}
	nodeId, _ := strconv.Atoi(this.GetString("node_id"))
	node, _ := models.GetNodeById(int32(nodeId))
	t.Node = &node
	t.Title = this.GetString("title")
	t.Body = this.GetString("body")
	v := models.UpdateTopic(&t)
	if v.HasErrors() {
		this.Data["topic"] = &t
		this.Data["nodes"] = models.FindAllNodes()
		this.TplName = "topics/edit.html"
		return
	}
	this.Redirect(fmt.Sprintf("/topics/%v", t.Id))
}

func (this *Topics) Delete() {
	this.requireUser()
	t, _ := models.GetTopicById(help.StrToInt32(this.Ctx.Input.Param(":id")))
	if !this.isOwner(t) {
		beego.NewFlash().Error("没有修改的权限")
		this.Redirect("/")
	}

	err := t.DeleteTopic()
	if err != nil {
		return
	}

	this.Redirect("/topics")
}

func (this *Topics) Watch() {
	this.requireUserForJSON()
	t, _ := models.GetTopicById(help.StrToInt32(this.Ctx.Input.Param(":id")))
	this.currentUser.Watch(t)
	this.successJSON(t.WatchesCount + 1)
}

func (this *Topics) UnWatch() {
	this.requireUserForJSON()
	t, _ := models.GetTopicById(help.StrToInt32(this.Ctx.Input.Param(":id")))
	this.currentUser.UnWatch(t)
	this.successJSON(t.WatchesCount - 1)
}

func (this *Topics) Star() {
	this.requireUserForJSON()
	topicId := help.StrToInt32(this.Ctx.Input.Param(":id"))
	t, _ := models.GetTopicById(topicId)
	this.currentUser.Star(t)
	this.successJSON(t.StarsCount + 1)
}

func (this *Topics) UnStar() {
	this.requireUserForJSON()
	t, _ := models.GetTopicById(help.StrToInt32(this.Ctx.Input.Param(":id")))
	this.currentUser.UnStar(t)
	this.successJSON(t.StarsCount - 1)
}

//加精或者埋贴
func (this *Topics) Rank() {
	this.requireAdmin()
	rankVal := 0
	switch strings.ToLower(this.GetString("v")) {
	case "nopoint":
		rankVal = models.RankNoPoint
	case "awesome":
		rankVal = models.RankAwesome
	default:
		rankVal = models.RankNormal
	}

	t, _ := models.GetTopicById(help.StrToInt32(this.Ctx.Input.Param(":id")))
	err := t.UpdateRank(rankVal)
	if err != nil {
		return
	}
	this.Redirect(fmt.Sprintf("/topics/%v", t.Id))
}

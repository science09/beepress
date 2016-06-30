package controllers

import (
	"beepress/models"

	"github.com/astaxie/beego"
	"golang.org/x/net/websocket"
)

type Home struct {
	BaseController
}

func (this *Home) Index() {
	topics, _ := models.FindTopicPages("popular", 0, 1, 10)
	this.Data["title"] = "Home"
	this.Data["controller_name"] = "Topic"
	this.Data["topics"] = topics
	this.Layout = "layout/layout.html"
	this.TplName = "home/index.html"
}

func (this *Home) Message() {
	if !this.isLogined() {
		beego.Error("not logined")
	}

	//ws := c.Request.Websocket
	//
	//Subscribe(c.currentUser.NotifyChannelId(), func(out interface{}) {
	//	err := websocket.JSON.Send(ws, out)
	//	if err != nil {
	//		fmt.Println("WebSocket send error: ", err)
	//	}
	//})

	ws := &websocket.Conn{}
	models.Subscribe(this.currentUser.NotifyChannelId(), func(out interface{}) {
		err := websocket.JSON.Send(ws, out)
		if err != nil {
			beego.Error("WebSocket send error: ", err)
		}
	})
}

func (this *Home) Search() {
	var topics []*models.Topic
	var page int
	queryStr := this.GetString("q")
	this.Ctx.Input.Bind(&page, "page")
	topics, pageInfo := models.GetSearchPages(queryStr, page, 5)
	this.Data["topics"] = topics
	this.Data["SearchName"] = queryStr
	this.Data["page_info"] = pageInfo
	this.Layout = "layout/layout.html"
	this.TplName = "search/show.html"
}

func (this *Home) About() {
	this.Layout = "layout/layout.html"
	this.TplName = "home/about.html"
}
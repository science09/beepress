package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
)

type Home struct  {
	BaseController
}

func (this *Home) Index() {
	this.Data["title"] = "Home"
	this.Data["controller_name"] = "Topic"
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

}

func (this *Home) Search() {
	url := fmt.Sprintf("https://google.com?q=site:ruby-china.org %v", this.GetString("q"))
	this.Redirect(url)
}
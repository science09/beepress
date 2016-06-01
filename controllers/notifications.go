package controllers

type Notifications struct {
	BaseController
}

func (this *Notifications) Index() {
	this.requireUser()
	var page int
	//this.Params.Bind(&page, "page")
	notes, pageInfo := this.currentUser.NotificationsPage(page, 8)
	this.currentUser.ReadNotifications(notes)
	this.Data["title"] = "社区"
	this.Data["notifications"] = notes
	this.Data["page_info"] = pageInfo
	this.TplName = "notifications/index.html"
}

func (this *Notifications) Clear() {
	this.requireUser()

	this.currentUser.ClearNotifications()
	this.Redirect("/notifications")
}
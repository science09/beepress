package routers

import (
	"beepress/controllers"
	"github.com/astaxie/beego"
)

//GET     /				  Home.Index
//WS      /msg            Home.Message
//GET     /search         Home.Search
//
//GET 	  /signup 			Accounts.New
//POST 	  /signup 			Accounts.Create
//GET 	  /signin 			Accounts.Login
//POST	  /signin			Accounts.LoginCreate
//POST    /signout          Accounts.Logout
//GET     /account/edit     Accounts.Edit
//POST    /account          Accounts.Update
//GET     /account/password     Accounts.Password
//POST    /account/password/update     Accounts.UpdatePassword
//
//GET     /topics           Topics.Index
//GET     /topics/feed           Topics.Feed
//GET     /topics/node/:node_id      Topics.Index("node")
//GET     /topics/popular           Topics.Index("popular")
//GET     /topics/recent           Topics.Index("recent")
//
//GET     /topics/new       Topics.New
//GET     /topics/:id       Topics.Show
//GET     /topics/:id/edit       Topics.Edit
//POST    /topics/:id       Topics.Update
//POST    /topics           Topics.Create
//POST    /topics/:id/delete       Topics.Delete
//POST    /topics/:id/reply          Replies.Create
//POST    /topics/:id/watch       Topics.Watch
//POST    /topics/:id/unwatch       Topics.UnWatch
//POST    /topics/:id/star       Topics.Star
//POST    /topics/:id/unstar       Topics.UnStar
//POST    /topics/:id/rank         Topics.Rank
//
//GET     /replies/:id/edit          Replies.Edit
//POST    /replies/:id          Replies.Update
//POST    /replies/:id/delete          Replies.Delete
//
//GET     /notifications            Notifications.Index
//POST    /notifications/clear      Notifications.Clear
//
//GET     /nodes            Nodes.Index
//POST    /nodes            Nodes.Create
//GET     /nodes/:id/edit            Nodes.Edit
//POST    /nodes/:id            Nodes.Update
//POST    /nodes/:id/delete            Nodes.Delete
//
//GET     /settings                 Settings.Index
//GET     /settings/:key/edit       Settings.Edit
//POST    /settings/:key            Settings.Update
//
//GET     /captcha                  App.Captcha
//
//GET     /:login             Users.Show
//GET     /:login/topics      Users.Topics

func init() {
	beego.Router("/", &controllers.Home{}, "*:Index")
	beego.Router("/search", &controllers.Home{}, "*:Search")
	beego.Router("/signup", &controllers.Accounts{}, "get:New;post:Create")
	beego.Router("/signin", &controllers.Accounts{}, "get:Login;post:LoginCreate")
	beego.Router("/signout", &controllers.Accounts{}, "post:Logout")
	beego.Router("/account/edit", &controllers.Accounts{}, "get:Edit")
	beego.Router("/account", &controllers.Accounts{}, "post:Update")
	beego.Router("/account/password", &controllers.Accounts{}, "get:Password")
	beego.Router("/account/password/update", &controllers.Accounts{}, "post:UpdatePassword")
	beego.Router("/user/?:login", &controllers.Users{}, "get:Show")
	beego.Router("/user/:login/topics", &controllers.Users{}, "get:Topics")
	beego.Router("/captcha", &controllers.BaseController{}, "get:Captcha")
	//
	beego.Router("/topics", &controllers.Topics{}, "get:Index;post:Create")
	beego.Router("/topics/node/?:node_id", &controllers.Topics{}, "get:TopicNode")
	beego.Router("/topics/popular", &controllers.Topics{}, "get:Popular")
	beego.Router("/topics/recent", &controllers.Topics{}, "get:Recent")
	beego.Router("/topics/feed", &controllers.Topics{}, "get:Feed")
	beego.Router("/topics/new", &controllers.Topics{}, "get:New")
	beego.Router("/topics/?:id", &controllers.Topics{}, "get:Show;post:Update")
	beego.Router("/topics/?:id/edit", &controllers.Topics{}, "get:Edit")
	beego.Router("/topics/?:id/delete", &controllers.Topics{}, "post:Delete")
	beego.Router("/topics/?:id/reply", &controllers.Replies{}, "post:Create")
	beego.Router("/topics/?:id/watch", &controllers.Topics{}, "post:Watch")
	beego.Router("/topics/?:id/unwatch", &controllers.Topics{}, "post:UnWatch")
	beego.Router("/topics/?:id/star", &controllers.Topics{}, "post:Star")
	beego.Router("/topics/?:id/unstar", &controllers.Topics{}, "post:UnStar")
	beego.Router("/topics/?:id/rank", &controllers.Topics{}, "post:Rank")
	//
	beego.Router("/replies/?:id/edit", &controllers.Replies{}, "get:Edit")
	beego.Router("/replies/?:id", &controllers.Replies{}, "post:Update")
	beego.Router("/replies/?:id/delete", &controllers.Replies{}, "post:Delete")
	//
	beego.Router("/notifications", &controllers.Notifications{}, "get:Index")
	beego.Router("/notifications/clear", &controllers.Notifications{}, "post:Clear")
	//
	beego.Router("/nodes", &controllers.Nodes{}, "get:Index;post:Create")
	beego.Router("/nodes/?:id/edit", &controllers.Nodes{}, "get:Edit")
	beego.Router("/nodes/?:id", &controllers.Nodes{}, "post:Update")
	beego.Router("/nodes/?:id/delete", &controllers.Nodes{}, "post:Delete")
	//
	beego.Router("/settings", &controllers.Settings{}, "get:Index")
	beego.Router("/settings/?:key/edit", &controllers.Settings{}, "get:Edit")
	beego.Router("/settings/?:key", &controllers.Settings{}, "post:Update")
}

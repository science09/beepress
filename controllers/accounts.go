package controllers

import (
	"regexp"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"github.com/science09/beepress/models"
)

//用户账户相关的控制器
type Accounts struct {
	BaseController
}

var (
	regexRequireUserActions, _ = regexp.Compile("Edit|Update|Password|UpdatePassword")
)

func (this *Accounts) NestPrepare() {
	_, action := this.GetControllerAndAction()
	if regexRequireUserActions.MatchString(action) {
		this.requireUser()
	}
}

func (this *Accounts) New() {
	this.Data["title"] = "注册新用户"
	this.Data["controller_name"] = "Topic"
	this.Layout = "layout/layout.html"
	this.TplName = "accounts/new.html"
}

func (this *Accounts) Create() {
	u := new(models.User)
	flash := beego.NewFlash()

	this.Data["title"] = "注册新用户"
	this.Data["controller_name"] = "Topic"
	this.Layout = "layout/layout.html"
	this.TplName = "accounts/new.html"

	login := this.GetString("login")
	email := this.GetString("email")
	password := this.GetString("password")
	rePassword := this.GetString("password-confirm")

	if !this.validateCaptcha(this.GetString("captcha")) {
		flash.Error("验证码不正确")
		flash.Store(&this.Controller)
		return
	}

	newUser, err := u.Signup(login, email, password, rePassword)
	if err != nil {
		flash.Error(err.Error())
		flash.Store(&this.Controller)
		return
	}
	this.storeUser(&newUser)
	flash.Success("注册成功")
	this.Redirect("/")
}

func (this *Accounts) Login() {
	this.Data["title"] = "登录"
	this.Data["controller_name"] = "Topic"
	this.Layout = "layout/layout.html"
	this.TplName = "accounts/login.html"

}

func (this *Accounts) LoginCreate() {
	u := models.User{}
	newUser := models.User{}
	v := validation.Validation{}
	flash := beego.NewFlash()

	this.Data["title"] = "登录"
	this.Data["controller_name"] = "Topic"
	this.Layout = "layout/layout.html"
	this.TplName = "accounts/login.html"
	captcha := this.GetString("captcha")
	beego.Info(captcha)
	if !this.validateCaptcha(captcha) {
		flash.Error("验证码不正确")
		flash.Store(&this.Controller)
		return
	}

	newUser, v = u.Signin(this.GetString("login"), this.GetString("password"))
	if v.HasErrors() {
		for _, val := range v.Errors {
			flash.Error(val.Message)
		}
		flash.Store(&this.Controller)
		return
	}

	this.storeUser(&newUser)
	flash.Success("登录成功，欢迎再次回来。")
	flash.Store(&this.Controller)
	this.Redirect("/")
}

func (this *Accounts) Logout() {
	flash := beego.NewFlash()
	this.clearUser()
	flash.Success("登出成功")
	flash.Store(&this.Controller)
	this.Redirect("/")
}

func (this *Accounts) Edit() {
	this.Data["title"] = "个人设置"
	this.Data["method_name"] = "edit"
	this.TplName = "accounts/edit.html"
}

func (this *Accounts) Update() {
	flash := beego.NewFlash()
	this.Data["title"] = "个人设置"
	this.Data["method_name"] = "edit"
	this.TplName = "accounts/edit.html"
	this.currentUser.Email = this.GetString("email")
	this.currentUser.GitHub = this.GetString("github")
	this.currentUser.Twitter = this.GetString("twitter")
	this.currentUser.Tagline = this.GetString("tagline")
	this.currentUser.Location = this.GetString("location")
	this.currentUser.Description = this.GetString("description")
	var u models.User
	u = this.currentUser
	_, err := models.UpdateUserProfile(&u)
	if err != nil {
		flash.Error(err.Error())
		flash.Store(&this.Controller)
		return
	}

	flash.Success("个人信息修改成功")
	flash.Store(&this.Controller)
	//this.Redirect("/account/edit")
}

func (this *Accounts) Password() {
	this.Data["title"] = "个人设置"
	this.Data["method_name"] = "password"
	this.TplName = "accounts/password.html"
}

func (this *Accounts) UpdatePassword() {
	this.Data["method_name"] = "password"
	this.TplName = "accounts/password.html"
	flash := beego.NewFlash()
	err := this.currentUser.UpdatePassword(this.GetString("password"), this.GetString("new-password"), this.GetString("confirm-password"))
	if err != nil {
		beego.Debug(err)
		flash.Error(err.Error())
		flash.Store(&this.Controller)
		return
	}
	flash.Success("密码修改成功")
	flash.Store(&this.Controller)
}

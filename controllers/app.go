package controllers

import (
	"beepress/models"
	"bytes"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"github.com/dchest/captcha"
	"reflect"
	"strconv"
	"strings"
)

type NestPreparer interface {
	NestPrepare()
}

type BaseController struct {
	beego.Controller
	currentUser models.User
}

const (
	JSON_CODE_NO_LOGIN = -1
)

func (this *BaseController) Prepare() {

	this.prependCurrentUser()
	this.Data["validation"] = nil
	this.Data["logined"] = this.isLogined()
	this.Data["current_user"] = this.currentUser
	this.Data["app_name"] = "bepress"
	this.Layout = "layout/layout.html"
	this.Data["controller_name"] = "Topic"
	//this.Data["controller_name"] = inflections.Underscore(this.Name)
	//this.Data["method_name"] = inflections.Underscore(this.MethodName)
	//this.Data["route_name"] = fmt.Sprintf("%v#%v", inflections.Underscore(this.Name), inflections.Underscore(this.MethodName))

	if app, ok := this.AppController.(NestPreparer); ok {
		app.NestPrepare()
	}
}

func (this *BaseController) Finish() {
	newParams := make(map[string]string, len(this.Input()))
	for key := range this.Input() {
		newParams[key] = this.Input().Get(key)
	}
	if len(newParams) > 0 {
		this.Data["params"] = newParams
	}
}

func (this *BaseController) prependCurrentUser() {
	beego.Info("prependCurrentUser", this.GetSession("user_id"))
	userId := this.GetSession("user_id")
	if userId == nil {
		return
	}
	uid, _ := strconv.Atoi(this.GetSession("user_id").(string))
	user, _ := models.GetUserById(uid)
	this.currentUser = *user
}

func (this *BaseController) CurrentUser() models.User {
	if this.currentUser.Id > 0 {
		return this.currentUser
	}
	this.prependCurrentUser()
	return this.currentUser
}

func (this *BaseController) storeUser(u *models.User) {
	if u.Id == 0 {
		return
	}
	this.SetSession("user_id", fmt.Sprintf("%v", u.Id))
}

func (this *BaseController) clearUser() {
	this.DelSession("user_id")
}

func (this *BaseController) isLogined() bool {
	return this.currentUser.Id > 0
}

func (this *BaseController) requireUser() {
	if !this.isLogined() {
		beego.Info("你还未登录哦")
		flash := beego.NewFlash()
		flash.Error("你还未登录哦")
		flash.Store(&this.Controller)
		this.Redirect("/signin")
		return
	} else {
		beego.Info("current_user { id: ", this.currentUser.Id,
			", login: ", this.currentUser.Login, " }")
	}
}

func (this *BaseController) requireUserForJSON() {
	if this.isLogined() {
		//c.Finish(c.errorJSON(JSON_CODE_NO_LOGIN, "还未登录"))
	}
}

func (this *BaseController) requireAdmin() {
	this.requireUser()
	if !this.currentUser.IsAdmin() {
		flash := beego.NewFlash()
		flash.Error("此功能需要管理员权限。")
		this.Redirect("/")
		return
	}
}

func (this *BaseController) isOwner(obj interface{}) bool {
	if this.currentUser.IsAdmin() {
		return true
	}
	objType := reflect.TypeOf(obj)
	switch objType.String() {
	case "models.Topic":
		return this.currentUser.Id == obj.(models.Topic).UserId
	case "*models.Topic":
		return this.currentUser.Id == obj.(*models.Topic).UserId
	case "models.User":
		return this.currentUser.Id == obj.(models.User).Id
	case "*models.User":
		return this.currentUser.Id == obj.(*models.User).Id
	case "models.Reply":
		return this.currentUser.Id == obj.(models.Reply).UserId
	case "*models.Reply":
		return this.currentUser.Id == obj.(*models.Reply).UserId
	default:
		panic(fmt.Sprintf("Invalid isOwner type: %v, %v, name: %v", obj, objType, objType.Name()))
	}

	return false
}

func (this *BaseController) Redirect(url string) {
	this.Ctx.Redirect(302, url)
}

type AppResult struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (this *BaseController) errorJSON(code int, msg string) {
	result := AppResult{Code: code, Msg: msg}
	this.Data["json"] = &result
	this.ServeJSON()
}

func (this *BaseController) errorsJSON(code int, errs []*validation.Error) {
	msgs := make([]string, len(errs))
	for i, err := range errs {
		msgs[i] = err.Message
	}
	result := AppResult{Code: code, Msg: strings.Join(msgs, "\n")}
	this.Data["json"] = &result
	this.ServeJSON()
}

func (this *BaseController) successJSON(data interface{}) {
	result := AppResult{Code: 0, Data: data}
	this.Data["json"] = &result
	this.ServeJSON()
}

func (this *BaseController) Captcha( /*id string*/ ) {
	captchaId := captcha.NewLen(4)
	beego.Info("captchaid:", captchaId)
	this.SetSession("captcha_id", captchaId)

	var buffer bytes.Buffer
	captcha.WriteImage(&buffer, captchaId, 200, 80)

	this.Ctx.Output.ContentType("image/png")
	this.Ctx.Output.Status = 200
	this.Ctx.WriteString(buffer.String())
}

func (this *BaseController) validateCaptcha(code string) bool {
	cap := this.GetSession("captcha_id")
	return captcha.VerifyString(cap.(string), code)
}

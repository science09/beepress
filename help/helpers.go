package help

import (
	"beepress/models"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/huacnlee/timeago"
	"github.com/huacnlee/train"
	"html/template"
	"math/rand"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/validation"
)

var (
	shareSites     = []string{"twitter", "weibo", "facebook", "google_plus", "email"}
	shareSiteIcons = map[string]string{
		"twitter":     "twitter",
		"weibo":       "weibo",
		"facebook":    "facebook-square",
		"google_plus": "google-plus-square",
		"email":       "envelope",
	}
)

func Init() {
	beego.AddFuncMap("plus", plus)
	beego.AddFuncMap("join", join)
	beego.AddFuncMap("is_owner", is_owner)
	beego.AddFuncMap("error_messages", error_messages)
	beego.AddFuncMap("timeago", time_ago)
	beego.AddFuncMap("markdown", MarkDown)
	beego.AddFuncMap("user_name_tag", user_name_tag)
	beego.AddFuncMap("user_avatar_tag", user_avatar_tag)
	beego.AddFuncMap("node_name_tag", node_name_tag)
	beego.AddFuncMap("paginate", paginate)
	beego.AddFuncMap("watch_tag", watch_tag)
	beego.AddFuncMap("star_tag", star_tag)
	beego.AddFuncMap("awesome_icon_tag", awesome_icon_tag)
	beego.AddFuncMap("active_class", active_class)
	beego.AddFuncMap("node_list", node_list)
	beego.AddFuncMap("select_tag", select_tag)
	beego.AddFuncMap("total", total)
	beego.AddFuncMap("setting", setting)
	beego.AddFuncMap("random_tip", random_tip)
	beego.AddFuncMap("share_button", share_button)
	beego.AddFuncMap("stylesheet_link_tag", train.StylesheetTag)
	beego.AddFuncMap("javascript_include_tag", train.JavascriptTag)
}

func plus(a, b int) int {
	return a + b
}

func join(args []string, split string) string {
	return strings.Join(args, split)
}

func is_owner(u models.User, obj interface{}) bool {
	if u.IsAdmin() {
		return true
	}

	switch obj.(type) {
	case models.User:
		u1 := obj.(models.User)
		return u1.Id == u.Id
	case models.Topic:
		t := obj.(models.Topic)
		return u.Id == t.UserId
	case models.Reply:
		r := obj.(models.Reply)
		return u.Id == r.UserId
	}

	return false
}

func error_messages(args ...interface{}) interface{} {

	out := ""

	if len(args) == 0 {
		return out
	}

	switch args[0].(type) {
	case string:
		return out
	case validation.Error:
		v := args[0].(validation.Validation)
		var parts []string
		if !v.HasErrors() {
			return out
		}
		parts = append(parts, "<div class=\"alert alert-block alert-warning\" role=\"alert\"><ul>")
		for _, err := range v.Errors {
			parts = append(parts, fmt.Sprintf("<li>%s %s</li>", err.Key, template.HTMLEscaper(err.Message)))
		}
		parts = append(parts, "</ul></div>")
		out = strings.Join(parts, "")

	//case revel.Validation:
	//	v := args[0].(revel.Validation)
	//	var parts []string
	//	if !v.HasErrors() {
	//		return out
	//	}
	//
	//	parts = append(parts, "<div class=\"alert alert-block alert-warning\" role=\"alert\"><ul>")
	//	for _, err := range v.ErrorMap() {
	//		parts = append(parts, fmt.Sprintf("<li>%s %s</li>", err.Key, template.HTMLEscaper(err.Message)))
	//	}
	//
	//	parts = append(parts, "</ul></div>")
	//	out = strings.Join(parts, "")
	default:
		return out
	}

	return template.HTML(out)
}

func time_ago(t time.Time) string {
	return timeago.Chinese.Format(t)
}

func MarkDown(text string) interface{} {
	bytes := []byte(text)
	outBytes := MarkdownGitHub(bytes)
	htmlText := string(outBytes[:])
	return template.HTML(htmlText)
}

func user_name_tag(obj interface{}) interface{} {
	out := "未知用户"
	switch obj.(type) {
	case models.User:
		u := obj.(models.User)
		if u.NewRecord() {
			return out
		}
		out = fmt.Sprintf("<a href='/%v' class='uname'>%v</a>", template.HTMLEscapeString(u.Login), template.HTMLEscapeString(u.Login))
	default:
		login := fmt.Sprintf("%v", obj)
		out = fmt.Sprintf(`<a href="/%v" class="uname">%v</a>`, template.HTMLEscapeString(login), template.HTMLEscapeString(login))

	}

	return template.HTML(out)
}

func user_avatar_tag(obj interface{}, size string) interface{} {
	out := ""
	if obj != nil {
		u := (obj).(models.User)
		if u.NewRecord() {
			return out
		}

		out = fmt.Sprintf("<a href=\"/user/%v\" class=\"uname\"><img src=\"%v\" class=\"media-object avatar-%v\" /></a>", template.HTMLEscapeString(u.Login), u.GavatarURL(size), size)
	}

	return template.HTML(out)
}

func node_name_tag(obj interface{}) interface{} {
	out := ""
	switch obj.(type) {
	case models.Node:
		n := obj.(models.Node)
		if n.NewRecord() {
			return out
		}
		out = fmt.Sprintf("<a href='/topics/node/%v' class='node-name'>%v</a>", n.Id, template.HTMLEscapeString(n.Name))
	}

	return template.HTML(out)
}

func paginate(pageInfo models.Pagination) interface{} {
	if pageInfo.TotalPages < 2 {
		return ""
	}

	linkFlag := "?"

	if strings.ContainsAny(pageInfo.Path, "?") {
		linkFlag = "&"
	}

	html := `<ul class="pager">`
	if pageInfo.Page > 1 {
		html += fmt.Sprintf(`<li class="previous"><a href="%s%spage=%d"><i class="fa fa-arrow-left" aria-hidden="true"></i> 上一页</a></li>`, pageInfo.Path, linkFlag, pageInfo.Page-1)
	} else {
		html += fmt.Sprintf(`<li class="previous disabled"><a href="#""><i class="fa fa-arrow-left" aria-hidden="true"></i> 上一页</a></li>`)
	}

	html += fmt.Sprintf(`<li class="info"><samp>%d</samp> / <samp>%d</samp></li>`, pageInfo.Page, pageInfo.TotalPages)

	if pageInfo.Page < pageInfo.TotalPages {
		html += fmt.Sprintf(`<li class="next"><a href="%s%spage=%d">下一页 <i class="fa fa-arrow-right" aria-hidden="true"></i></a></li>`, pageInfo.Path, linkFlag, pageInfo.Page+1)
	} else {
		html += fmt.Sprintf(`<li class="next disabled"><a href="#">下一页 <i class="fa fa-arrow-right" aria-hidden="true"></i></a></li>`)
	}
	html += "</ul>"

	return template.HTML(html)
}

func watch_tag(t models.Topic, u models.User) interface{} {
	out := ""
	if t.NewRecord() {
		return out
	}
	out = fmt.Sprintf(`<a href="#" data-id="%v" class="watch" title="关注此话题，当有新回帖的时候会收到通知"><i class="fa fa-eye"></i> 关注</a>`, t.Id)

	if u.NewRecord() {
		return template.HTML(out)
	}

	if u.IsWatched(t) {
		out = fmt.Sprintf(`<a href="#" data-id="%v" class="watch followed" title="点击取消关注"><i class="fa fa-eye"></i> 已关注</a>`, t.Id)
	}

	return template.HTML(out)
}

func star_tag(t models.Topic, u models.User) interface{} {
	out := ""
	if t.NewRecord() {
		return out
	}
	label := fmt.Sprintf("%v 人收藏", t.StarsCount)
	out = fmt.Sprintf(`<a href="#" data-id="%v" data-count="%v" class="star"><i class="fa fa-star-o"></i> %v</a>`, t.Id, t.StarsCount, label)

	if u.NewRecord() {
		return template.HTML(out)
	}

	if u.IsStared(t) {
		out = fmt.Sprintf(`<a href="#" data-id="%v" data-count="%v" class="star followed"><i class="fa fa-star"></i> %v</a>`, t.Id, t.StarsCount, label)
	}

	return template.HTML(out)
}

func awesome_icon_tag(t models.Topic) interface{} {
	out := ""
	if !t.IsAwesome() {
		return out
	}

	out = `<i class="fa fa-diamond awesome" title="精华帖标记"></i>`
	return template.HTML(out)
}

func active_class(a string, b string) string {
	if strings.EqualFold(a, b) {
		return " active "
	} else {
		return ""
	}
}

//节点列表模板函数
func node_list() interface{} {
	groups := models.FindAllNodeGroups()
	outs := []string{}
	subs := []string{}
	outs = append(outs, `<div class="row node-list">`)
	for _, group := range groups {
		beego.Info("ddddd", group.Name)
		subs = []string{
			`<div class="node media clearfix">`,
			fmt.Sprintf(`<label class="media-left">%v</label>`, group.Name),
			`<div class="nodes media-body">`,
		}
		for _, node := range group.Nodes {
			subs = append(subs, fmt.Sprintf(`<span class="name"><a href="/topics/node/%v">%v</a></span>`, node.Id, node.Name))
		}
		subs = append(subs, "</div></div>")

		outs = append(outs, strings.Join(subs, ""))
	}
	outs = append(outs, "</div>")
	return template.HTML(strings.Join(outs, ""))
}

func select_tag(objs interface{}, nameKey, valueKey, formName string, defaultValue interface{}) interface{} {
	objsVal := reflect.ValueOf(objs)
	if objsVal.Kind() != reflect.Slice {
		fmt.Println("Give a bad params, objs need to be a Slice")
		return ""
	}

	outs := []string{}

	subs := []string{}
	var nameField reflect.Value
	var valueField reflect.Value

	defaultName := "请选择"

	for i := 0; i < objsVal.Len(); i++ {
		val := objsVal.Index(i)
		nameField = val.FieldByName(nameKey)
		valueField = val.FieldByName(valueKey)
		subs = append(subs, fmt.Sprintf(`
               <li data-id="%v"><a href="#">%v</a></li>
            `, valueField.Int(), nameField.String()))

		// check current name
		if strings.EqualFold(fmt.Sprintf("%v", valueField.Int()), fmt.Sprintf("%v", defaultValue)) {
			defaultName = nameField.String()
		}
	}

	outs = append(outs, `<div class="input-group-btn md-dropdown">`)
	outs = append(outs, fmt.Sprintf(`
        <button class="btn btn-default dropdown-toggle" type="button" data-toggle="dropdown" aria-expanded="false">
            <span data-bind="label">%v</span> <span class="caret"></span>
        </button>
        <input type="hidden" data-bind="value" value="%v" name="%v" />`,
		defaultName, defaultValue, formName))

	outs = append(outs, `<ul class="dropdown-menu" role="menu">`)
	outs = append(outs, strings.Join(subs, ""))
	outs = append(outs, `</ul>`)
	outs = append(outs, `</div>`)

	return template.HTML(strings.Join(outs, ""))
}

func total(key string) interface{} {
	switch key {
	case "users":
		return models.UsersCountCached()
	case "topics":
		return models.TopicsCountCached()
	case "replies":
		return models.RepliesCountCached()
	}

	return nil
}

func setting(key string) interface{} {
	return template.HTML(models.GetSetting(key))
}

func random_tip() interface{} {
	tipText := models.GetSetting("tips")
	tips := strings.Split(tipText, "\n")
	return template.HTML(tips[rand.Intn(len(tips))])
}

func share_button(title, url string) interface{} {
	results := []string{}
	results = append(results, fmt.Sprintf(`<div class="social-share-button" data-via="gochina" data-title="%v" data-url="%v">`, title, url))
	for _, siteName := range shareSites {
		link := fmt.Sprintf(`<a rel="nofollow" data-site="%v" href="#"><i class="fa fa-%v"></i></a>`, siteName, shareSiteIcons[siteName])
		results = append(results, link)
	}
	results = append(results, "</div>")
	results = append(results, `<a href="#" class="share-button"><i class="fa fa-share-square-o"></i> 转发</a>`)
	return template.HTML(strings.Join(results, ""))
}

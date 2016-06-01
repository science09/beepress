package models

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

var ADMIN_LOGINS = []string{"admin"}

type User struct {
	Id          int32     `orm:"pk;auto"`
	Login       string    `orm:"size(100)" valid:"Required;MaxSize(100)"`
	Password    string    `orm:"size(100)" valid:"Required;MaxSize(100)"`
	Email       string    `orm:"size(100)" valid:"Email;MaxSize(100)"`
	Avatar      string    `orm:"size(100)"`
	GitHub      string    `orm:"size(100)"`
	Twitter     string    `orm:"size(100)"`
	HomePage    string    `orm:"size(100)"`
	Tagline     string    `orm:"size(100)"`
	Description string    `orm:"size(250)"`
	Location    string    `orm:"size(200)"`
	Topics      []*Topic  `orm:"-"`
	Replies     []*Reply  `orm:"-"`
	CreatedAt   time.Time `orm:"type(datetime);auto_now_add"`
	UpdatedAt   time.Time `orm:"type(datetime);auto_now"`
	DeletedAt   time.Time `orm:"type(datetime);null"`
}

func (u *User) TableName() string {
	return TableName("user")
}

func (u User) BeforeCreate() error {
	u.Login = strings.ToLower(u.Login)
	return nil
}

//生成图片路径
func (u User) GavatarURL(size string) string {
	//emailMD5 := u.EncodePassword(u.Email)
	//return fmt.Sprintf("https://ruby-china.org/avatar/%v?s=%v", emailMD5, size)
	return "/static/img/avatar/user.png"
}

func (u User) NotifyChannelId() string {
	return fmt.Sprintf("notify/%v", u.Id)
}

func (u User) SameAs(obj interface{}) bool {
	return obj.(User).Id == u.Id
}

func (u User) IsAdmin() bool {
	for _, str := range ADMIN_LOGINS {
		if u.Login == str {
			return true
		}
	}
	return false
}

func (u *User) NewRecord() bool {
	return u.Id <= 0
}

func (u User) UnReadNotificationsCount() (count int32) {
	n := new(Notification)
	count64, err := orm.NewOrm().QueryTable(n.TableName()).Filter("user_id", u.Id).Filter("read", 0).Count()
	if err != nil {
		count64 = 0
		beego.Error(err)
	}
	count = int32(count64)
	//db.Model(&Notification{}).Where("`user_id` = ? and `read` = 0", u.Id).Count(&count)
	return
}

func PushNotifyInfoToUser(userId int32, note Notification, isNew bool) {
	u := User{}
	u.Id = userId
	unreadCount := u.UnReadNotificationsCount()

	// Set read, update client unread_count
	if !isNew {
		go PushMessage(u.NotifyChannelId(), &NotifyInfo{UnreadCount: unreadCount, IsNew: false})
		return
	}

	actor := User{}
	if note.Id > 0 {
		//db.First(&actor, note.ActorId)
	}
	info := NotifyInfo{
		UnreadCount: unreadCount,
		IsNew:       true,
		Title:       note.NotifyableTitle(),
		Avatar:      actor.GavatarURL("256x256"),
		Path:        note.NotifyableURL(),
	}

	beego.Info("[Push] Notify:", info)

	go PushMessage(u.NotifyChannelId(), &info)
}

func (u User) EncodePassword(raw string) (md5Digest string) {
	data := []byte(raw)
	result := md5.Sum(data)
	md5Digest = hex.EncodeToString(result[:])
	return
}

func (u User) Signup(login string, email string, password string, passwordConfirm string) (user User, v validation.Validation) {
	u.Login = strings.ToLower(strings.Trim(login, " "))
	u.Email = strings.ToLower(strings.Trim(email, " "))
	v = validation.Validation{}
	v.Required(email, "Email").Message("不能为空")
	v.MinSize(login, 5, "用户名").Message("最少要 5 个字符")
	v.MinSize(password, 6, "密码").Message("最少要 6 个字符")
	v.Email(email, "Email").Message("格式不正确")

	if password != passwordConfirm {
		v.Error("密码与确认密码不一致")
	}

	var existCount int64
	existCount, _ = orm.NewOrm().QueryTable(TableName("user")).Filter("login", login).Count()
	fmt.Println("login name as: ", login, " have ", existCount)
	if existCount > 0 {
		v.SetError("user", "帐号已经被注册")
	}

	if v.HasErrors() {
		return u, v
	}

	u.Password = u.EncodePassword(password)
	_, err := orm.NewOrm().Insert(&u)
	if err != nil {
		v.Error(fmt.Sprintf("服务器异常, %v", err))
	}
	beego.Info("created user: ", u)
	return u, v
}

func (u User) Signin(login string, password string) (user User, v validation.Validation) {
	login = strings.TrimSpace(login)
	if len(password) == 0 {
		v.Error("还未输入密码")
	}
	err := orm.NewOrm().QueryTable(u).Filter("login", login).Filter("password", u.EncodePassword(password)).One(&user)
	if err != nil {
		beego.Error(err.Error())
		v.Error("帐号密码不正确")
	}
	return user, v
}

func UpdateUserProfile(u User) (user User, v validation.Validation) {
	v.Required(u.Email, "email").Message("格式不正确")
	//v.Email(u.Email).Key("Email").Message("格式不正确")
	if v.HasErrors() {
		return u, v
	}
	willUpdateUser := User{
		Email:       u.Email,
		Location:    u.Location,
		Description: u.Description,
		GitHub:      u.GitHub,
		Twitter:     u.Twitter,
		Tagline:     u.Tagline,
	}
	//err := db.First(&u, u.Id).Updates(willUpdateUser).Error
	_, err := orm.NewOrm().Update(&willUpdateUser)
	if err != nil {
		v.Error(err.Error())
	}
	return u, v
}

func (u User) UpdatePassword(oldPassword, newPassword, confirmPassword string) (v validation.Validation) {
	user := User{}
	//v.Required(oldPassword).Key("旧密码").Message("不能为空")
	v.Required(oldPassword, "旧密码").Message("不能为空")
	//db.First(&user, "id = ? and password = ?", u.Id, u.EncodePassword(oldPassword))
	orm.NewOrm().Read(&user)
	if user.NewRecord() {
		v.Error("旧密码不正确")
	}
	//v.MinSize(newPassword, 6).Key("新密码").Message("最少要 6 个子符")
	v.MinSize(newPassword, 6, "新密码").Message("最少要 6 个子符")
	if newPassword != confirmPassword {
		v.Error("新密码与确认新密码输入的内容不一致")
	}
	if v.HasErrors() {
		return v
	}

	//err := db.Model(u).Update("password", u.EncodePassword(newPassword)).Error
	_, err := orm.NewOrm().Update(&user, "password")
	if err != nil {
		v.Error(err.Error())
	}
	return v
}

func FindUserByLogin(login string) (u User, err error) {
	err = orm.NewOrm().QueryTable(TableName("user")).Filter("login", strings.ToLower(login)).One(&u)
	return
}

func GetUserById(id int) (u *User, err error) {
	user := User{}
	err = orm.NewOrm().QueryTable(TableName("user")).Filter("id", id).One(&user)
	u = &user
	return
}

func UsersCountCached() (count int) {
	if !Cache.IsExist("users/total") {
		if count, err := orm.NewOrm().QueryTable(TableName("user")).Count(); err == nil {
			go Cache.Put("users/total", int(count), 30*time.Minute)
		}
	} else {
		count = (Cache.Get("users/total")).(int)
	}

	return
}

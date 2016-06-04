package models

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

var (
	filePath     = "/static/avatar/"
	suffix       = ".svg"
	ADMIN_LOGINS = []string{"admin"}
	bgColor      = []string{"#aa4325", "#e84e40", "#C41411", "#ad1457", "#673ab7", "#5677FC",
		"#7e57C2", "#7986CB", "#03A9F4", "#00BCD4", "#009688", "#795548", "#ff6e40", "#607d8B"}
)

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
	emailMD5 := u.EncodePassword(u.Login)
	file := emailMD5 + ".svg"
	return fmt.Sprintf("%s%v?s=%v", filePath, file, size)
}

func GetRand() int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Int()
}

func GenerateAvatar(name string) error {
	strByte := []byte(name)
	u := User{Login: name}
	fileName := u.EncodePassword(name) + suffix
	bg := bgColor[GetRand()%len(bgColor)]
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	file := dir + filePath + fileName
	f, err := os.Create(file)
	defer f.Close()
	if err != nil {
		return err
	}
	svg := fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="240" height="240">
  <rect x="0" y="0" width="240" height="240" rx="6" ry="6" fill="%s" />
  <text fill="white" x="120" y="120" font-size="160" font-weight="bold" text-anchor="middle" style="dominant-baseline: central;">%s</text>
</svg>`, bg, strings.ToUpper(string(strByte[0])))
	f.WriteString(svg)

	return err
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

func (u User) Signup(login string, email string, password string, passwordConfirm string) (user User, err error) {
	u.Login = strings.ToLower(strings.Trim(login, " "))
	u.Email = strings.ToLower(strings.Trim(email, " "))
	v := validation.Validation{}
	v.MinSize(login, 5, "用户名").Message("用户名最少要 5 个字符")
	v.Required(email, "Email").Message("Email不能为空")
	v.MinSize(password, 6, "密码").Message("密码最少要 6 个字符")
	v.Email(email, "Email").Message("格式不正确")

	if password != passwordConfirm {
		v.Error("密码与确认密码不一致")
	}
	if v.HasErrors() {
		for _, val := range v.Errors {
			err = errors.New(val.Message)
			return u, err
		}
	}

	var existCount int64
	existCount, _ = orm.NewOrm().QueryTable(TableName("user")).Filter("login", login).Count()
	fmt.Println("login name as: ", login, " have ", existCount)
	if existCount > 0 {
		err := errors.New("帐号已经被注册")
		return u, err
	}

	GenerateAvatar(u.Login)
	u.Password = u.EncodePassword(password)
	_, err = orm.NewOrm().Insert(&u)
	if err != nil {
		msg := fmt.Sprintf("服务器异常, %v", err)
		err = errors.New(msg)
	}
	beego.Info("created user: ", u)
	return u, err
}

func (u User) Signin(login string, password string) (user User, v validation.Validation) {
	login = strings.TrimSpace(login)
	if len(password) == 0 {
		v.Error("还未输入密码")
	}
	err := orm.NewOrm().QueryTable(u).Filter("login", login).Filter("password", u.EncodePassword(password)).One(&user)
	if err != nil {
		v.Error("帐号或密码不正确")
	}
	return user, v
}

func UpdateUserProfile(u *User) (user User, err error) {
	v := validation.Validation{}
	v.Required(u.Login, "login").Message("用户名不能为空")
	v.Required(u.Email, "email").Message("格式不正确")
	if v.HasErrors() {
		for _, val := range v.Errors {
			err = errors.New(val.Message)
		}
		return *u, err
	}
	//willUpdateUser := User{
	//	Email:       u.Email,
	//	Location:    u.Location,
	//	Description: u.Description,
	//	GitHub:      u.GitHub,
	//	Twitter:     u.Twitter,
	//	Tagline:     u.Tagline,
	//}
	_, err = orm.NewOrm().Update(u)
	if err != nil {
		err = errors.New("服务器更新失败")
	}
	return *u, err
}

func (u User) UpdatePassword(oldPassword, newPassword, confirmPassword string) error {
	v := validation.Validation{}
	v.Required(oldPassword, "旧密码").Message("旧密码不能为空")
	orm.NewOrm().Read(&u)
	if u.NewRecord() {
		v.Error("旧密码不正确")
	}
	v.MinSize(newPassword, 6, "新密码").Message("新密码最少要 6 个子符")
	if newPassword != confirmPassword {
		v.Error("新密码与确认新密码输入的内容不一致")
	}
	if newPassword == oldPassword {
		v.Error("新密码不能和旧密码一致")
	}
	if v.HasErrors() {
		for _, err := range v.Errors {
			return errors.New(err.Message)
		}
	}
	u.Password = u.EncodePassword(newPassword)
	_, err := orm.NewOrm().Update(&u, "password")
	if err != nil {
		err = errors.New("更新密码失败")
	}
	return err
}

func FindUserByLogin(login string) (u *User, err error) {
	user := &User{}
	err = orm.NewOrm().QueryTable(TableName("user")).Filter("login", strings.ToLower(login)).One(user)
	u = user
	return
}

func GetUserById(id int) (u *User, err error) {
	user := &User{}
	err = orm.NewOrm().QueryTable(TableName("user")).Filter("id", id).One(user)
	u = user
	return
}

func GetUserCount() int {
	count, err := orm.NewOrm().QueryTable(TableName("user")).Count()
	if err != nil {
		return 0
	}
	return int(count)
}

func UsersCountCached() (count int) {
	if !Cache.IsExist("users/total") {
		count = GetUserCount()
		go Cache.Put("users/total", count, 30*time.Minute)
	} else {
		count = (Cache.Get("users/total")).(int)
	}

	return
}

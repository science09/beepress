package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"time"
)

type Notification struct {
	Id             int32
	NotifyType     string
	Read           bool `orm:"default(false)"`
	UserId         int32
	User           *User `orm:"-"`
	ActorId        int32
	Actor          *User `orm:"-"`
	NotifyableType string
	NotifyableId   int32
	CreatedAt      time.Time `orm:"auto_now_add;type(datetime)"`
	UpdatedAt      time.Time `orm:"auto_now;type(datetime)"`
}

type NotifyInfo struct {
	UnreadCount int32  `json:"unread_count"`
	IsNew       bool   `json:"is_new"`
	Title       string `json:"title"`
	Avatar      string `json:"avatar"`
	Path        string `json:"path"`
}

func (n *Notification) TableName() string {
	return TableName("notification")
}

func (n *Notification) Topic() (t Topic) {
	if n.NotifyableType == "Topic" {
		t, _ = GetTopicById(n.NotifyableId)
	}
	return
}

func (n *Notification) Reply() (r Reply, err error) {
	if n.NotifyableType == "Reply" {
		r, err = GetReplyById(n.NotifyableId)
	}
	return
}

func (n *Notification) NotifyableTitle() string {
	switch n.NotifyableType {
	case "Topic":
		return n.Topic().Title
	case "Reply":
		//t := Topic{}
		//db.First(&t, n.Reply().TopicId)
		r, _ := n.Reply()
		t, _ := GetTopicById(r.TopicId)
		return t.Title
	default:
		return ""
	}
}

func (n *Notification) NotifyableURL() string {
	switch n.NotifyableType {
	case "Topic":
		return fmt.Sprintf("/topics/%v", n.NotifyableId)
	case "Reply":
		r, _ := n.Reply()
		return fmt.Sprintf("/topics/%v", r.TopicId)
	default:
		return ""
	}
}

func createNotification(notifyType string, userId int32, actorId int32, notifyableType string, notifyableId int32) error {
	note := Notification{
		NotifyType:     notifyType,
		UserId:         userId,
		ActorId:        actorId,
		NotifyableType: notifyableType,
		NotifyableId:   notifyableId,
	}

	//exitCount := 0
	//db.Model(Notification{}).Where(
	//	"user_id = ? and actor_id = ? and notifyable_type = ? and notifyable_id = ?",
	//	userId, actorId, notifyableType, notifyableId).Count(&exitCount)
	//if exitCount > 0 {
	//	return nil
	//}
	//

	//err := db.Save(&note).Error
	_, err := orm.NewOrm().Insert(&note)
	go PushNotifyInfoToUser(userId, note, true)

	return err
}

func (r *Reply) NotifyReply() error {
	if r.NewRecord() {
		return nil
	}

	//t := Topic{}
	//
	//err := db.Find(&t, r.TopicId).Error
	t, err := GetTopicById(r.TopicId)
	if err != nil {
		return nil
	}

	if t.UserId != r.UserId {
		// 跳过回复人
		go createNotification("Reply", t.UserId, r.UserId, "Reply", r.Id)
	}

	followerIds := t.FollowerIds()
	for _, followerId := range followerIds {
		if followerId == r.UserId || followerId == t.UserId {
			// 跳过回复人和发帖人
			continue
		}
		go createNotification("Reply", followerId, r.UserId, "Reply", r.Id)
	}

	return nil
}

func NotifyMention(userId, actorId int32, notifyableType string, notifyableId int32) error {
	return createNotification("Mention", userId, actorId, notifyableType, notifyableId)
}

func (u User) NotificationsPage(page, perPage int) (notes []Notification, pageInfo Pagination) {
	pageInfo = Pagination{}
	//pageInfo.Query = db.Model(&Notification{}).Preload("Actor")
	//pageInfo.Query = pageInfo.Query.Where("user_id = ?", u.Id).Order("id desc")
	//
	//pageInfo.Path = "/notifications"
	//pageInfo.PerPage = perPage
	//pageInfo.Paginate(page).All(&notes)
	return
}

func (u User) ReadNotifications(notes []Notification) error {
	ids := []int32{}
	for _, note := range notes {
		ids = append(ids, note.Id)
	}
	if len(ids) > 0 {
		//err := db.Model(Notification{}).Where("user_id = ? and id in (?)", u.Id, ids).Update("read", true).Error
		_, err := orm.NewOrm().QueryTable(TableName("notification")).Filter("user_id", u.Id).Filter("id__in", ids).Update(orm.Params{"read": true})
		go PushNotifyInfoToUser(u.Id, Notification{}, false)
		return err
	}

	return nil
}

func (u User) ClearNotifications() error {
	_, err := orm.NewOrm().QueryTable(TableName("notification")).Filter("user_id", u.Id).Delete()
	return err
}

func (n *Notification) IsTopic() bool {
	return n.NotifyType == "Topic"
}

func (n *Notification) IsReply() bool {
	return n.NotifyType == "Reply"
}

func (n *Notification) IsMention() bool {
	return n.NotifyType == "Mention"
}

func (n *Notification) IsNotifyableReply() bool {
	return n.NotifyableType == "Reply"
}

func (n *Notification) IsNotifyableTopic() bool {
	return n.NotifyableType == "Topic"
}

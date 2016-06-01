package models

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"time"
)

type Topic struct {
	Id                 int32 `orm:"pk;auto"`
	UserId             int32
	User               User `orm:"-"`
	NodeId             int32
	Node               *Node `orm:"-"`
	Title              string
	Body               string    `orm:"type(text)"`
	Replies            []*Reply  `orm:"-"`
	RepliesCount       int32     `orm:"default(0)"`
	LastActiveMark     int64     `orm:"default(0)"`
	LastRepliedAt      time.Time `orm:"type(datetime);null"`
	LastReplyId        int32
	LastReplyUserId    int32
	LastReplyUser      *User `orm:"-"`
	LastReplyUserLogin string
	StarsCount         int32     `orm:"default(0)"`
	WatchesCount       int32     `orm:"default(0)"`
	Rank               int32     `orm:"default(0)"`
	CreatedAt          time.Time `orm:"type(datetime);auto_now_add"`
	UpdatedAt          time.Time `orm:"type(datetime);auto_now"`
	DeletedAt          time.Time `orm:"type(datetime);null"`
}

const (
	RankNoPoint = -1
	RankNormal  = 0
	RankAwesome = 1
)

func (t *Topic) TableName() string {
	return TableName("topic")
}

func GetTopicCount() (count int32) {
	total, _ := orm.NewOrm().QueryTable(TableName("topic")).Count()
	count = int32(total)
	return
}

func GetTopicById(id int32) (topic Topic, err error) {
	err = orm.NewOrm().QueryTable(TableName("topic")).Filter("id", id).One(&topic)
	u, _ := GetUserById(int(topic.UserId))
	topic.User = *u
	return
}

func GetTopicByUserId(user_id int32) (topic []Topic, err error) {
	_, err = orm.NewOrm().QueryTable(TableName("topic")).Filter("user_id", user_id).All(&topic)
	return
}

func FindTopicPages(channel string, nodeId, page, perPage int) (topics []Topic, pageInfo Pagination) {
	pageInfo = Pagination{}
	//pageInfo.Query = db.Model(&Topic{}).Preload("User").Preload("Node")
	//
	qs := orm.NewOrm().QueryTable(TableName("topic"))
	switch channel {
	case "recent":
		//pageInfo.Query = pageInfo.Query.Order("id desc")
		qs.OrderBy("-id").All(&topics)
	case "popular":
		cond := orm.NewCondition().And("rank", 1).Or("stars_count__gte", 5)
		qs.SetCond(cond).All(&topics)

		//pageInfo.Query = pageInfo.Query.Where("rank = 1 or stars_count >= 5")
		//pageInfo.Query = pageInfo.Query.Order("last_active_mark desc, id desc")
	case "node":
		qs.Filter("node_id", nodeId).OrderBy("-last_active_mark", "-id").All(&topics)
		//pageInfo.Query = pageInfo.Query.Where("node_id = ?", nodeId)
		//pageInfo.Query = pageInfo.Query.Order("last_active_mark desc, id desc")
	default:
		qs.OrderBy("-last_active_mark", "-id").All(&topics)
		for _, v := range topics {
			v.User.Login = "admin"
		}
		//pageInfo.Query = pageInfo.Query.Where("rank >= 0").Order("last_active_mark desc, id desc")
	}
	//pageInfo.Path = "/topics"
	//pageInfo.PerPage = perPage
	//pageInfo.Paginate(page).Find(&topics)
	return
}

func CreateTopic(t *Topic) error {
	//先验证topic的有效性

	(*t).LastActiveMark = time.Now().Unix()
	_, err := orm.NewOrm().Insert(t)
	if err != nil {
		beego.Debug("err:", err)
		err = errors.New("服务器异常创建失败")
	}
	return err
}

func UpdateTopic(t *Topic) validation.Validation {
	v := validation.Validation{}
	if v.HasErrors() {
		return v
	}

	_, err := orm.NewOrm().Insert(t)
	if err != nil {
		v.Error("服务器异常更新失败")
	}
	return v
}

func (t *Topic) UpdateLastReply(reply *Reply) (err error) {
	if reply == nil {
		return errors.New("Reply is nil")
	}
	//db.First(&reply.User, reply.UserId)
	//err = db.Exec(`UPDATE topics SET updated_at = ?, last_active_mark = ?, last_replied_at = ?,
	//	last_reply_id = ?, last_reply_user_login = ?, last_reply_user_id = ? WHERE id = ?`,
	//	time.Now(),
	//	time.Now().Unix(),
	//	time.Now(),
	//	reply.Id,
	//	reply.User.Login,
	//	reply.UserId,
	//	reply.TopicId).Error

	return err
}

func (t *Topic) NewRecord() bool {
	return t.Id <= 0
}

func (t Topic) UpdateRank(rank int) error {
	if t.NewRecord() {
		return errors.New("Give a empty record.")
	}
	_, err := orm.NewOrm().QueryTable(&t).Update(orm.Params{"rank": rank})
	return err
}

func (t Topic) IsAwesome() bool {
	return t.Rank == RankAwesome
}

func (t Topic) IsNormal() bool {
	return t.Rank == RankNormal
}

func (t Topic) IsNoPoint() bool {
	return t.Rank == RankNoPoint
}

func (t Topic) URL() string {
	if t.NewRecord() {
		return ""
	}
	return fmt.Sprintf("%v/topics/%v", "https://127.0.0.1:3000", t.Id)
}

func (t Topic) FollowerIds() (ids []int32) {
	//db.Model(Followable{}).Where("follow_type = 'Watch' and topic_id = ?", t.Id).Pluck("user_id", &ids)
	return
}

func (t *Topic) DeleteTopic() error {
	_, err := orm.NewOrm().Delete(t)
	return err
}

func TopicsCountCached() (count int32) {
	if !Cache.IsExist("topics/total") {
		count = GetTopicCount()
		go Cache.Put("topics/total", count, 30*time.Minute)
	} else {
		count = Cache.Get("topics/total").(int32)
	}

	return
}

package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

type Topic struct {
	Id                 int32     `orm:"pk;auto"`
	User               *User     `orm:"rel(fk)"`
	Node               *Node     `orm:"rel(fk)"`
	Title              string    ``
	Body               string    `orm:"type(text)"`
	Replies            []*Reply  `orm:"reverse(many)"`
	RepliesCount       int32     `orm:"default(0)"`
	LastActiveMark     int64     `orm:"default(0)"`
	LastRepliedAt      time.Time `orm:"type(datetime);null"`
	LastReplyId        int32     ``
	LastReplyUser      *User     `orm:"-"`
	LastReplyUserLogin string    ``
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

func GetTopicCount() (count int) {
	total, _ := orm.NewOrm().QueryTable(TableName("topic")).Count()
	count = int(total)
	return
}

func GetTopicById(id int32) (topic Topic, err error) {
	err = orm.NewOrm().QueryTable(TableName("topic")).Filter("id", id).RelatedSel().One(&topic)
	if err != nil {
		beego.Error(err, id)
	}
	return
}

func GetTopicByUserId(user_id int32) (topic []Topic, err error) {
	_, err = orm.NewOrm().QueryTable(TableName("topic")).Filter("user_id", user_id).All(&topic)
	return
}

//获取最近的10篇文章
func GetRecentTopics(user_id int32) (topic []Topic, err error) {
	_, err = orm.NewOrm().QueryTable(TableName("topic")).Filter("user_id", user_id).Limit(10).All(&topic)
	return
}

func FindTopicPages(channel string, nodeId, page, perPage int) (topics []*Topic, pageInfo Pagination) {
	o := orm.NewOrm()
	qs := o.QueryTable(TableName("topic"))
	switch channel {
	case "recent":
		pageInfo.Query = qs.OrderBy("-id").RelatedSel()
	case "popular":
		cond := orm.NewCondition().And("rank", 1).Or("stars_count__gte", 5)
		pageInfo.Query = qs.SetCond(cond).RelatedSel().OrderBy("-last_active_mark", "-id")
		//pageInfo.Query = pageInfo.Query.Where("rank = 1 or stars_count >= 5")
		//pageInfo.Query = pageInfo.Query.Order("last_active_mark desc, id desc")
	case "node":
		//qs.Filter("node_id", nodeId).RelatedSel().OrderBy("-last_active_mark", "-id").All(&topics)
		pageInfo.Query = qs.Filter("node_id", nodeId).RelatedSel().OrderBy("-last_active_mark", "-id")
		//pageInfo.Query = pageInfo.Query.Where("node_id = ?", nodeId)
		//pageInfo.Query = pageInfo.Query.Order("last_active_mark desc, id desc")
	default:
		//qs.RelatedSel().OrderBy("-last_active_mark", "-id").All(&topics)
		//pageInfo.Query = pageInfo.Query.Where("rank >= 0").Order("last_active_mark desc, id desc")
		pageInfo.Query = qs.RelatedSel().OrderBy("-last_active_mark", "-id")
	}

	pageInfo.Path = "/topics"
	pageInfo.PerPage = perPage
	pageInfo.Paginate(page).All(&topics)

	return
}

func GetSearchPages(search string, page int, perPage int) (topics []*Topic, pageInfo Pagination) {
	o := orm.NewOrm()
	qs := o.QueryTable(&Topic{})
	cond := orm.NewCondition().And("title__icontains", search).Or("body__icontains", search)
	pageInfo.Query = qs.SetCond(cond).RelatedSel()
	path := fmt.Sprintf("/search?q=%v", search)
	pageInfo.Path = path
	pageInfo.PerPage = perPage
	pageInfo.Paginate(page).All(&topics)
	return
}

func (t *Topic) validate() validation.Validation {
	v := validation.Validation{}
	v.Required(t.Title, "标题").Message("不能为空")
	v.Required(t.Node.Id, "").Message("请选择节点名称")
	v.Required(t.Body, "内容").Message("不能为空")
	return v
}

func CreateTopic(t *Topic) error {
	//先验证topic的有效性
	if v := t.validate(); v.HasErrors() {
		for _, e := range v.Errors {
			return errors.New(e.Message)
		}
	}
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
	_, err := orm.NewOrm().Update(t)
	if err != nil {
		v.Error("服务器异常更新失败")
	}
	return v
}

func (t *Topic) UpdateLastReply(reply *Reply) (err error) {
	if reply == nil {
		return errors.New("Reply is nil")
	}

	o := orm.NewOrm()
	o.QueryTable(&User{}).Filter("id", reply.User.Id).One(reply.User)
	_, err = o.QueryTable(&Topic{}).Filter("id", reply.Topic.Id).Update(orm.Params{"updated_at": time.Now(),
		"last_active_mark": time.Now().Unix(), "last_replied_at": time.Now(), "last_reply_id": reply.Id,
		"last_reply_user_login": reply.User.Login}) //"last_reply_user_id": reply.User.Id

	//db.First(&reply.User, reply.UserId)
	//err = db.Exec(`UPDATE topics SET updated_at = ?, last_active_mark = ?, last_replied_at = ?,
	//	last_reply_id = ?, last_reply_user_login = ?, last_reply_user_id = ? WHERE id = ?`,
	//	time.Now(),d
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

func (t *Topic) UpdateRank(rank int) error {
	if t.NewRecord() {
		return errors.New("Give a empty record.")
	}
	_, err := orm.NewOrm().QueryTable(&t).Filter("id", t.Id).Update(orm.Params{"rank": rank})
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
	return fmt.Sprintf("%v/topics/%v", "https://127.0.0.1:8080", t.Id)
}

func (t *Topic) FollowerIds() (ids []int32) {
	follows := []Followable{}
	orm.NewOrm().QueryTable(&Followable{}).Filter("follow_type", "watch").Filter("topic_id", t.Id).All(&follows, "user_id")
	ids = make([]int32, len(follows))
	for key, val := range follows {
		ids[key] = int32(val.UserId)
	}

	return
}

func (t *Topic) DeleteTopic() error {
	_, err := orm.NewOrm().Delete(t)
	return err
}

func TopicsCountCached() (count int) {
	if !Cache.IsExist("topics/total") {
		count = GetTopicCount()
		go Cache.Put("topics/total", count, 30*time.Minute)
	} else {
		count = Cache.Get("topics/total").(int)
	}

	return
}

package models

import (
	"errors"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"time"
)

type Reply struct {
	Id        int32     `orm:"pk;auto"`
	UserId    int32     ``
	User      User      `orm:"-"`
	TopicId   int32     ``
	Topic     *Topic    `orm:"-"`
	Body      string    `orm:"type:text;"`
	IsDeleted bool      ``
	CreatedAt time.Time `orm:"auto_now_add;type(datetime)"`
	UpdatedAt time.Time `orm:"type(datetime);auto_now"`
}

func (r *Reply) TableName() string {
	return TableName("reply")
}

func (r *Reply) NewRecord() bool {
	return r.Id <= 0
}

func (r *Reply) validate() validation.Validation {
	v := validation.Validation{}
	switch r.NewRecord() {
	case false:
	default:
		v.Required(r.TopicId, "topic_id").Message("不能为空")
		v.Min(int(r.TopicId), 1, "topic_id").Message("不能为空")
		v.Required(r.UserId, "user_id").Message("不能为空")
		v.Min(int(r.UserId), 1, "user_id").Message("不能为空")
		v.MinSize(r.Body, 1, "内容").Message("不能为空")
		v.MaxSize(r.Body, 10000, "内容").Message("最多不允许超过 10000 个子符")
	}
	return v
}

func CreateReply(r *Reply) (err error) {
	//v := r.validate()
	//if v.HasErrors() {
	//	return v
	//}

	//需要先验证reply的正确性

	_, err = orm.NewOrm().Insert(r)
	if err != nil {
		err = errors.New("服务器异常创建失败")
	}

	return
}

func RepliesCountCached() (count int) {
	if !Cache.IsExist("replies/total") {
		count, _ := orm.NewOrm().QueryTable(TableName("reply")).Count()
		go Cache.Put("replies/total", int(count), 30*time.Minute)
	} else {
		count = Cache.Get("replies/total").(int)
	}

	return
}

func GetReplyById(id int32) (reply Reply, err error) {
	err = orm.NewOrm().QueryTable(TableName("reply")).Filter("id", id).One(&reply)
	return
}

func GetReplyByTopicId(id int32) (reply []Reply, err error) {
	_, err = orm.NewOrm().QueryTable(TableName("reply")).Filter("topic_id", id).OrderBy("id").All(&reply)
	return
}

func UpdateReply(reply Reply) error {
	_, err := orm.NewOrm().Update(&reply)
	return err
}

func (r *Reply) Del() (err error) {
	_, err = orm.NewOrm().QueryTable(TableName("reply")).Filter("Id", r.Id).Update(orm.Params{"IsDeleted": true})
	if err != nil {
		err = errors.New("删除评论失败!")
	}
	return
}

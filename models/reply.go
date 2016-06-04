package models

import (
	"errors"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

type Reply struct {
	Id        int32     `orm:"pk;auto"`
	User      *User     `orm:"rel(fk)"`
	Topic     *Topic    `orm:"rel(fk)"`
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
		v.Required(r.Topic.Id, "topic_id").Message("不能为空")
		v.Min(int(r.Topic.Id), 1, "topic_id").Message("不能为空")
		v.Required(r.User.Id, "user_id").Message("不能为空")
		v.Min(int(r.User.Id), 1, "user_id").Message("不能为空")
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

func GetReplyCount() int {
	count, err := orm.NewOrm().QueryTable(TableName("reply")).Count()
	if err != nil {
		return 0
	}
	return int(count)
}

func RepliesCountCached() (count int) {
	if !Cache.IsExist("replies/total") {
		count = GetReplyCount()
		go Cache.Put("replies/total", count, 30*time.Minute)
	} else {
		count = Cache.Get("replies/total").(int)
	}

	return
}

func GetReplyById(id int32) (reply Reply, err error) {
	err = orm.NewOrm().QueryTable(TableName("reply")).Filter("id", id).One(&reply)
	return
}

//获取文章的所有评论
func GetReplyByTopicId(id int32) (reply []Reply, err error) {
	_, err = orm.NewOrm().QueryTable(TableName("reply")).Filter("topic_id", id).RelatedSel().OrderBy("id").All(&reply)
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

package models

import (
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type Followable struct {
	Id         int32     `orm:"pk;auto"`
	FollowType string    `orm:"size(20)"`
	TopicId    int32     ``
	Topic      *Topic    `orm:"-"`
	UserId     int32     ``
	User       *User     `orm:"-"`
	CreatedAt  time.Time `orm:"type(datetime);auto_now_add"`
	UpdatedAt  time.Time `orm:"type(datetime);auto_now"`
}

func (f *Followable) TableName() string {
	return TableName("followable")
}

func (u User) isFollowed(ftype string, t Topic) bool {
	followId, _ := orm.NewOrm().QueryTable(&Followable{}).Filter("follow_type", ftype).
		Filter("topic_id", t.Id).Filter("user_id", u.Id).Count()
	if followId > 0 {
		return true
	} else {
		return false
	}
}

func (u User) follow(ftype string, t Topic) bool {
	if t.NewRecord() || u.NewRecord() {
		return false
	}

	if u.isFollowed(ftype, t) {
		return false
	}
	follow := Followable{FollowType: ftype, TopicId: t.Id, UserId: u.Id}
	if _, err := orm.NewOrm().Insert(&follow); err != nil {
		return false
	}
	t.updateFollowCounter(ftype)
	return true
}

func (u User) unFollow(ftype string, t Topic) bool {
	if t.NewRecord() || u.NewRecord() {
		return false
	}

	if !u.isFollowed(ftype, t) {
		return false
	}
	_, err := orm.NewOrm().QueryTable(&Followable{}).Filter("follow_type", ftype).
		Filter("topic_id", t.Id).Filter("user_id", u.Id).Delete()
	if err != nil {
		return false
	}
	t.updateFollowCounter(ftype)
	return true
}

func (t Topic) updateFollowCounter(ftype string) {
	o := orm.NewOrm()
	count, _ := o.QueryTable(&Followable{}).Filter("follow_type", ftype).Filter("topic_id", t.Id).Count()
	counterCacheKey := "watches_count"
	if ftype == "Star" {
		counterCacheKey = "stars_count"
	}
	_, err := o.QueryTable(&t).Filter("id", t.Id).Update(orm.Params{counterCacheKey: count})
	if err != nil {
		beego.Error("WARNING: updateFollowCounter execute failed: ", err)
	}
}

func (u User) IsWatched(t Topic) bool {
	return u.isFollowed("Watch", t)
}

func (u User) Watch(t Topic) bool {
	return u.follow("Watch", t)
}

func (u User) UnWatch(t Topic) bool {
	return u.unFollow("Watch", t)
}

func (u User) IsStared(t Topic) bool {
	return u.isFollowed("Star", t)
}

func (u User) Star(t Topic) bool {
	return u.follow("Star", t)
}

func (u User) UnStar(t Topic) bool {
	return u.unFollow("Star", t)
}

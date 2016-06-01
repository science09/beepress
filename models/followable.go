package models

import (
	"time"
	//"fmt"
)

type Followable struct {
	Id         int32  `orm:"pk;auto"`
	FollowType string `orm:"size(20)"`
	TopicId    int32
	Topic      *Topic `orm:"-"`
	UserId     int32
	User       *User     `orm:"-"`
	CreatedAt  time.Time `orm:"type(datetime);auto_now_add"`
	UpdatedAt  time.Time `orm:"type(datetime);auto_now"`
}

func (f *Followable) TableName() string {
	return TableName("followable")
}

func (u User) isFollowed(ftype string, t Topic) bool {
	var existCount int
	//DB.Model(&Followable{}).Where("follow_type = ? and topic_id = ? and user_id = ?", ftype, t.Id, u.Id).Count(&existCount)
	if existCount > 0 {
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

	//follow := Followable{FollowType: ftype, TopicId: t.Id, UserId: u.Id}
	//if DB.Save(&follow).Error != nil {
	//	return false
	//}
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

	//DB.Where("follow_type = ? and topic_id = ? and user_id = ?", ftype, t.Id, u.Id).Delete(&Followable{})
	t.updateFollowCounter(ftype)
	return true
}

func (t Topic) updateFollowCounter(ftype string) {
	//var count int
	//DB.Model(&Followable{}).Where("follow_type = ? and topic_id = ?", ftype, t.Id).Count(&count)

	//counterCacheKey := "watches_count"
	//if ftype == "Star" {
	//	counterCacheKey = "stars_count"
	//}

	//err := DB.Model(t).Update(counterCacheKey, count).Error
	//if err != nil {
	//	fmt.Println("WARNING: updateFollowCounter execute failed: ", err)
	//}
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

package models

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"

)

type Setting struct {
	Id  int32 `orm:"pk;auto"`
	Key string
	Val string
}

func (s *Setting) TableName() string {
	return TableName("setting")
}

func (s *Setting) AfterSave() {
	s.RewriteCache()
}

func settingCacheKey(key string) string {
	return fmt.Sprintf("setting/%v/v1", key)
}

func (s *Setting) RewriteCache() {
	Cache.Put(settingCacheKey(s.Key), s.Val, 7*24*time.Hour)
}

func FindSettingByKey(key string) (s Setting) {
	s.Key = key
	//DB.Where("`key` = ?", key).First(&s)
	orm.NewOrm().QueryTable(s.TableName()).Filter("key", key).One(&s)
	return s
}

func GetSetting(key string) (out string) {
	out = ""
	if Cache.IsExist(settingCacheKey(key)) {
		out = Cache.Get(settingCacheKey(key)).(string)
	}else {
		s := FindSettingByKey(key)
		if s.Id <= 0 {
			//保存数据库
		}
		out = s.Val
		s.RewriteCache()
	}

	//if err := Cache.Get(settingCacheKey(key), &out); err != nil {
	//	s := FindSettingByKey(key)
	//	if s.Id <= 0 {
	//		//		db.Save(&s)
	//	}
	//
	//	out = s.Val
	//	s.RewriteCache()
	//}

	return
}

func UpdateSetting(setting Setting) error {
	_, err := orm.NewOrm().Update(&setting)
	return err
}

package models

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

var Cache cache.Cache

type BaseModel struct {
	Id        int32
	DeletedAt *time.Time
}

func (m BaseModel) NewRecord() bool {
	return m.Id <= 0
}

func (m BaseModel) Destroy() error {
	//err := db.Delete(&m).Error
	err := errors.New("Not implemented!")
	return err
}

func (m BaseModel) IsDeleted() bool {
	return m.DeletedAt != nil
}

func Init() {
	dbHost := beego.AppConfig.String("db_host")
	dbPort := beego.AppConfig.String("db_port")
	dbName := beego.AppConfig.String("db_name")
	dbUser := beego.AppConfig.String("db_user")
	dbPass := beego.AppConfig.String("db_pass")
	timezone := beego.AppConfig.String("db_timezone")

	if dbPort == "" {
		dbPort = "3306"
	}
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
		dbUser, dbPass, dbHost, dbPort, dbName)
	if timezone != "" {
		connStr = connStr + "&loc=" + url.QueryEscape(timezone)
	}
	if beego.AppConfig.String("runmode") == "dev" {
		beego.Info("ormDebug:%v", orm.Debug)
	}
	orm.Debug = true
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", connStr)
	orm.RegisterModel(new(User), new(Topic), new(Setting), new(Reply), new(Followable), new(Node), new(NodeGroup), new(Notification))
	orm.RunSyncdb("default", false, true)

	var err error
	Cache, err = cache.NewCache("memory", `{"interval":60}`)
	if err != nil {
		beego.Error("cache init failed!")
	}

	InitNodeGroup()

	//db.LogMode(false)
	//logger = Logger{log.New(os.Stdout, "  ", 0)}
	//db.SetLogger(logger)
	//db.AutoMigrate(&User{}, &Topic{}, &Reply{}, &Node{}, &NodeGroup{}, &Followable{}, &Notification{}, &Setting{})
	//db.Model(NodeGroup{}).AddIndex("index_on_sort", "sort")
	//db.Model(Node{}).AddIndex("index_on_group_and_sort", "node_group_id", "sort")
	//db.Model(User{}).AddUniqueIndex("index_on_login", "login")
	//db.Model(Topic{}).AddIndex("index_on_user_id", "user_id")
	//db.Model(Topic{}).AddIndex("index_on_last_active_mark_deleted_at", "last_active_mark", "deleted_at")
	//db.Model(Topic{}).AddIndex("index_on_deleted_at", "deleted_at")
	//db.Model(Topic{}).AddIndex("index_on_rank", "rank")
	//db.Model(User{}).AddIndex("index_on_deleted_at", "deleted_at")
	//db.Model(Reply{}).AddIndex("index_on_deleted_at", "deleted_at")
	//db.Model(Notification{}).AddIndex("index_on_user_id", "user_id")
	//db.Model(Notification{}).AddIndex("index_on_notifyable", "notifyable_type", "notifyable_id")
	//db.Model(Setting{}).AddUniqueIndex("index_on_key", "key")
	//db.LogMode(true)
	//
	//initPubsub()
}

func TableName(name string) string {
	return beego.AppConfig.String("db_prefix") + name
}

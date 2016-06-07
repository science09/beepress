package models

import (
	"errors"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

//NodeGroupId int       ``
type Node struct {
	Id        int32      `orm:"pk;auto"`
	Name      string     `orm:"unique"`
	Summary   string     `orm:"type(text)"`
	NodeGroup *NodeGroup `orm:"rel(fk)"`
	Sort      int        `orm:"default(0)"`
	CreatedAt time.Time  `orm:"type(datetime);auto_now_add"`
	UpdatedAt time.Time  `orm:"type(datetime);auto_now"`
	DeletedAt time.Time  `orm:"type(datetime);null"`
}

func (n *Node) TableName() string {
	return TableName("node")
}

type NodeGroup struct {
	Id    int32   `orm:"pk;auto"`
	Name  string  `orm:"unique"`
	Sort  int     `orm:"default(0)"`
	Nodes []*Node `orm:"reverse(many)"`
}

func (n *NodeGroup) TableName() string {
	return TableName("node_groups")
}

//先初始化两个节点组
func InitNodeGroup() {
	o := orm.NewOrm()
	ng := &NodeGroup{}
	ng1 := &NodeGroup{Name: "Languages"}
	ng2 := &NodeGroup{Name: "Stack China"}
	if exist := o.QueryTable(ng.TableName()).Filter("name", ng1.Name).Exist(); !exist {
		orm.NewOrm().Insert(ng1)
	}
	if exist := o.QueryTable(ng.TableName()).Filter("name", ng2.Name).Exist(); !exist {
		orm.NewOrm().Insert(ng2)
	}
}

func (n *Node) validate() {
	//验证表单

}

func (n *Node) NewRecord() bool {
	return n.Id <= 0
}

func CreateNode(n *Node) (err error) {
	v := validation.Validation{}
	v.Valid(n)
	if v.HasErrors() {
		err = errors.New("验证出错")
		return
	}
	_, err = orm.NewOrm().Insert(n)
	if err != nil {
		err = errors.New("服务器异常创建失败")
	}
	return
}

func UpdateNode(n *Node) (err error) {
	_, err = orm.NewOrm().Update(n)
	if err != nil {
		err = errors.New("服务器异常更新失败")
	}
	return
}

//这一步需要优化
func FindAllNodeGroups() (groups []NodeGroup) {
	o := orm.NewOrm()
	_, err := o.QueryTable(TableName("node_groups")).All(&groups)
	if err != nil {
		beego.Error(err)
	}
	for key, val := range groups {
		o.LoadRelated(&val, "Nodes")
		groups[key].Nodes = val.Nodes
	}

	return
}

func FindAllNodes() (nodes []*Node) {
	orm.NewOrm().QueryTable(TableName("node")).OrderBy("name").All(&nodes)
	return
}

func FindNodeBySort(limit int) (nodes []*Node) {
	orm.NewOrm().QueryTable(TableName("node")).OrderBy("-sort name").Limit(limit).All(&nodes)
	return
}

func GetNodeById(id int32) (node Node, err error) {
	err = orm.NewOrm().QueryTable(TableName("node")).Filter("id", id).One(&node)
	return
}

func DeleteNode(n *Node) (err error) {
	_, err = orm.NewOrm().Delete(n)
	return
}

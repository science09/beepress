package controllers

import (
	"github.com/astaxie/beego"
	"github.com/science09/beepress/models"
	"strconv"
)

type Nodes struct {
	BaseController
}

func (this *Nodes) NestPrepare() {
	this.requireAdmin()
}

func (this *Nodes) loadNodeGroups() {
	groups := models.FindAllNodeGroups()
	this.Data["groups"] = groups
	beego.Error(groups)
}

func (this *Nodes) Index() {
	this.loadNodeGroups()
	nodes := models.FindAllNodes()
	this.Data["title"] = "节点管理"
	this.Data["nodes"] = nodes
	this.TplName = "nodes/index.html"
}

func (this *Nodes) Create() {
	flash := beego.NewFlash()
	nodeGroupId, _ := strconv.Atoi(this.GetString("node_group_id"))
	nodeGroup := &models.NodeGroup{Id: int32(nodeGroupId)}
	n := models.Node{
		Name:      this.GetString("name"),
		NodeGroup: nodeGroup,
	}

	err := models.CreateNode(&n)
	if err != nil {
		beego.Info("err", err.Error())
		this.loadNodeGroups()
		this.Data["node"] = n
		flash.Error(err.Error())
		flash.Store(&this.Controller)
		this.TplName = "nodes/index.html"
		return
	}
	flash.Success("节点创建成功")
	flash.Store(&this.Controller)
	this.Redirect("/nodes")
}

func (this *Nodes) Edit() {
	this.loadNodeGroups()

	strId := this.Ctx.Input.Param(":id")
	beego.Info("id:", strId)
	id, _ := strconv.ParseInt(strId, 10, 0)
	node, err := models.GetNodeById(int32(id))
	if err != nil {
		this.Abort("404")
	}
	this.Data["title"] = "修改节点"
	this.Data["node"] = node
	this.TplName = "nodes/edit.html"
}

func (this *Nodes) Update() {
	flash := beego.NewFlash()
	id, _ := strconv.Atoi(this.Ctx.Input.Param(":id"))
	node, err := models.GetNodeById(int32(id))
	if err != nil {
		beego.Info(err.Error())
		this.TplName = "nodes/edit.html"
		return
	}
	node.Name = this.GetString("name")
	node.Summary = this.GetString("summary")
	nodeGroupId, _ := strconv.Atoi(this.GetString("node_group_id"))
	node.NodeGroup.Id = int32(nodeGroupId)
	beego.Info("name, Summary, GroupId", node.Name, node.Summary, node.NodeGroup.Id)
	err = models.UpdateNode(&node)
	this.Data["node"] = node
	if err != nil {
		this.loadNodeGroups()
		this.TplName = "nodes/edit.html"
		return
	}
	flash.Success("节点更新成功")
	flash.Store(&this.Controller)
	this.Redirect("/nodes")
}

func (this *Nodes) Delete() {
	id, _ := strconv.Atoi(this.Ctx.Input.Param(":id"))
	node, err := models.GetNodeById(int32(id))
	if err != nil {
		this.TplName = "nodes/edit.html"
		return
	}
	models.DeleteNode(&node)
	this.Redirect("/nodes")
}

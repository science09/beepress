package main

import (
	"github.com/science09/beepress/help"
	"github.com/science09/beepress/models"
	_ "github.com/science09/beepress/routers"

	"github.com/astaxie/beego"
	"github.com/huacnlee/train"
)

func main() {
	models.Init()
	help.Init()
	train.Config.AssetsPath = "assets"
	train.Config.SASS.DebugInfo = true
	train.Config.SASS.LineNumbers = true
	train.Config.Verbose = true
	train.Config.BundleAssets = false
	train.ConfigureHttpHandler(nil)
	beego.SetStaticPath("/assets", "assets")
	beego.Run()
}

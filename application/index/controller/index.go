package controller

import (
	"github.com/kataras/iris/v12"
	"gocms/application/model"
	"gocms/application/service"
)

type IndexController struct {
	BaseController
}

func (ic *IndexController) GetIndex(ctx iris.Context) {
	ctx.ViewData("nav", "index_index")

	categoryM := model.Category{}
	service.InitDB().Where(map[string]interface{}{"classmodule": "photo", "classtitle": "banner"}).Find(&categoryM)
	banner_list := []model.Photo{}
	service.InitDB().Where(map[string]interface{}{"tid": categoryM.Classid}).Find(&banner_list)
	ctx.ViewData("banner_list", banner_list)

	article_list := []model.Article{}
	service.InitDB().Where("attrtj like '%1%'").Find(&article_list)
	ctx.ViewData("article_list", article_list)

	ctx.View("index/index.html")
}

package controller

import "github.com/kataras/iris/v12"

type ArticleController struct {
	BaseController
}

func (ac *ArticleController) GetIndex(ctx iris.Context) {
	ctx.ViewData("nav", "article_index")
	ctx.ViewData("index", "hello word")
	ctx.View("article/index.html")
}

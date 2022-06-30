package controller

import "github.com/kataras/iris/v12"

type AboutController struct {
	BaseController
}

func (ac *AboutController) GetIndex(ctx iris.Context) {
	ctx.ViewData("nav", "about_index")

	ctx.ViewData("index", "hello word")
	ctx.View("about/index.html")
}

func (ac *AboutController) GetJoin(ctx iris.Context) {
	ctx.ViewData("nav", "about_join")
	ctx.ViewData("index", "hello word")
	ctx.View("about/join.html")
}

func (ac *AboutController) GetJoininfo(ctx iris.Context) {
	ctx.ViewData("nav", "about_join")
	ctx.ViewData("index", "hello word")
	ctx.View("about/join_info.html")
}

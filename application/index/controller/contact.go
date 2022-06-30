package controller

import "github.com/kataras/iris/v12"

type ContactController struct {
	BaseController
}

func (cc *ContactController) GetIndex(ctx iris.Context) {
	ctx.ViewData("nav", "contact_index")
	ctx.ViewData("index", "hello word")
	ctx.View("contact/index.html")
}

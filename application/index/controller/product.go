package controller

import "github.com/kataras/iris/v12"

type ProductController struct {
	BaseController
}

func (pc *ProductController) GetIndex(ctx iris.Context) {
	ctx.ViewData("nav", "product_index")
	ctx.ViewData("index", "hello word")
	ctx.View("product/index.html")
}

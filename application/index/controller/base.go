package controller

import "github.com/kataras/iris/v12"

type BaseController struct {
}

func IndexAuth(ctx iris.Context) {
	ctx.Next()
}

package controller

import (
	"github.com/kataras/iris/v12"
	"strconv"
)

type MessageController struct {
	BaseController
}

func (mc *MessageController) GetIndex(ctx iris.Context) {

}

/**
未读消息个数
*/
func (mc *MessageController) PostNoreadnum(ctx iris.Context) {
	ctx.WriteString("0")
}

/**
消息列表
*/
func (mc *MessageController) GetList(ctx iris.Context) {
	var list []map[string]interface{}
	list = append(list, map[string]interface{}{"id": 1, "text": "..."})
	ctx.JSON(map[string]interface{}{"num": 0, "list": list})
}

/**
读消息详情
*/
func (mc *MessageController) PostInfo(ctx iris.Context) {
	id, _ := strconv.Atoi(ctx.PostValue("id"))
	text := "..."
	ctx.JSON(map[string]interface{}{"id": id, "text": text})
}

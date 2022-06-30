package controller

import (
	"github.com/kataras/iris/v12"
	"gocms/application/model"
	"strconv"
)

type ModelController struct {
	BaseController
}

func (mc *ModelController) PostList(ctx iris.Context) {
	page, _ := strconv.Atoi(ctx.PostValue("page"))
	rows, _ := strconv.Atoi(ctx.PostValue("rows"))
	limit_start := (page - 1) * rows

	ename := ctx.PostValue("ename")
	list := []model.Model{}
	var count int64
	if ename == "" {
		model.InitDB().Offset(limit_start).Limit(rows).Find(&list)
		model.InitDB().Model(model.Model{}).Count(&count)
	} else {
		model.InitDB().Where("ename like ?", "%"+ename+"%").Offset(limit_start).Limit(rows).Find(&list)
		model.InitDB().Model(model.Model{}).Where("ename like ?", "%"+ename+"%").Count(&count)
	}
	var pages int
	if len(list)%rows == 0 {
		pages = len(list) / rows
	} else {
		pages = len(list)/rows + 1
	}
	ctx.JSON(map[string]interface{}{"pages": pages, "total": count, "rows": list})
}

func (mc *ModelController) PostAdd(ctx iris.Context) {
	modelM := model.Model{}
	modelM.Ename = ctx.PostValue("ename")
	if modelM.Ename == "" {
		mc.TopjuiError(ctx, "模型名必填")
	} else {
		var count int64
		model.InitDB().Model(modelM).Where(map[string]interface{}{"ename": modelM.Ename}).Count(&count)
		if count > 0 {
			mc.TopjuiError(ctx, "模型名已经存在")
		}
	}
	modelM.Cname = ctx.PostValue("cname")
	modelM.Sort, _ = strconv.Atoi(ctx.PostValue("sort"))
	modelM.Status, _ = strconv.Atoi(ctx.PostValue("status"))
	res := model.InitDB().Model(modelM).Create(&modelM)
	if res.Error == nil {
		mc.TopjuiSucess(ctx, "新增成功")
	} else {
		mc.TopjuiError(ctx, "新增失败")
	}
}

func (mc *ModelController) GetEdit(ctx iris.Context) {
	id, _ := strconv.Atoi(ctx.URLParam("id"))
	if id == 0 {
		mc.TopjuiError(ctx, "Id参数丢失")
	}
	modelM := model.Model{}
	model.InitDB().Model(modelM).Where(map[string]interface{}{"id": id}).Find(&modelM)
	ctx.ViewData("info", modelM)
	mc.BaseController.GetEdit(ctx)
}

func (mc *ModelController) PostEdit(ctx iris.Context) {
	id, _ := strconv.Atoi(ctx.PostValue("id"))
	if id == 0 {
		mc.TopjuiError(ctx, "Id参数丢失")
	}
	modelM := model.Model{}
	modelM.Ename = ctx.PostValue("ename")
	modelM.Cname = ctx.PostValue("cname")
	modelM.Sort, _ = strconv.Atoi(ctx.PostValue("sort"))
	modelM.Status, _ = strconv.Atoi(ctx.PostValue("status"))
	res := model.InitDB().Model(modelM).Where(map[string]interface{}{"id": id}).Updates(&modelM)
	if res.Error == nil {
		mc.TopjuiSucess(ctx, "编辑成功")
	} else {
		mc.TopjuiError(ctx, "编辑失败")
	}
}

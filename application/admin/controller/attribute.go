package controller

import (
	"encoding/json"
	"github.com/kataras/iris/v12"
	"gocms/application/model"
	"gocms/application/service"
	"gorm.io/gorm"
	"strconv"
)

type AttributeController struct {
	BaseController
}

func (ac *AttributeController) PostList(ctx iris.Context) {
	page, _ := strconv.Atoi(ctx.PostValue("page"))
	rows, _ := strconv.Atoi(ctx.PostValue("rows"))
	limit_start := (page - 1) * rows

	ename := ctx.PostValue("ename")
	list := []model.Attribute{}
	var count int64
	if ename == "" {
		service.InitDB().Select([]string{"ANY_VALUE(id)", "model_ename", "ANY_VALUE(attrname)", "ANY_VALUE(attrvalue)", "ANY_VALUE(sort)", "ANY_VALUE(status)"}).Group("model_ename").Offset(limit_start).Limit(rows).Find(&list)
		service.InitDB().Model(model.Attribute{}).Select([]string{"ANY_VALUE(id)", "model_ename", "ANY_VALUE(attrname)", "ANY_VALUE(attrvalue)", "ANY_VALUE(sort)", "ANY_VALUE(status)"}).Group("model_ename").Count(&count)
	} else {
		service.InitDB().Where(map[string]interface{}{"model_ename": ename}).Offset(limit_start).Limit(rows).Find(&list)
		service.InitDB().Model(model.Attribute{}).Where(map[string]interface{}{"model_ename": ename}).Count(&count)
	}
	var pages int
	if len(list)%rows == 0 {
		pages = len(list) / rows
	} else {
		pages = len(list)/rows + 1
	}
	ctx.JSON(map[string]interface{}{"pages": pages, "total": count, "rows": list})
}

func (ac *AttributeController) GetAdd(ctx iris.Context) {
	modelM := []model.Model{}
	service.InitDB().Select([]string{"ename", "cname"}).Where(map[string]interface{}{"status": 1}).Find(&modelM)
	var model_list []map[string]interface{}
	for key, vo := range modelM {
		vo_map := make(map[string]interface{})
		vo_map["text"] = vo.Cname
		vo_map["value"] = vo.Ename
		if key == 0 {
			vo_map["selected"] = true
		}
		model_list = append(model_list, vo_map)
	}
	model_json_byte, _ := json.Marshal(model_list)
	ctx.ViewData("model_list", string(model_json_byte))
	ctx.View(RequestController + "/add.html")
}

func (ac *AttributeController) PostAdd(ctx iris.Context) {
	model_ename := ctx.URLParam("ename")
	if model_ename == "" {
		model_ename = ctx.PostValue("model_ename")
	}
	attrname := ctx.PostValue("attrname")
	attrvalue := ctx.PostValue("attrvalue")
	if model_ename == "" || attrname == "" || attrvalue == "" {
		ac.TopjuiError(ctx, "参数错误")
	}
	attributeM := model.Attribute{}

	if service.InitDB().Where(map[string]interface{}{"model_ename": model_ename, "attrvalue": attrvalue}).Find(&attributeM); attributeM.Id != 0 {
		ac.TopjuiError(ctx, "该属性值已存在")
	}

	attributeM.ModelEname = model_ename
	attributeM.Attrname = attrname
	attributeM.Attrvalue = attrvalue
	attributeM.Status = 1
	attributeM.Sort = 99
	res := service.InitDB().Create(&attributeM)
	if res.Error == nil {
		ac.TopjuiSucess(ctx, "添加成功")
	} else {
		ac.TopjuiError(ctx, "添加失败")
	}
}

func (ac *AttributeController) PostEdit(ctx iris.Context) {
	id, _ := strconv.Atoi(ctx.PostValue("id"))
	attrname := ctx.PostValue("attrname")
	attrvalue := ctx.PostValue("attrvalue")
	if id == 0 || attrname == "" || attrvalue == "" {
		ac.TopjuiError(ctx, "参数错误")
	}
	attributeM := model.Attribute{}
	attributeM.Id = id
	attributeM.Attrname = attrname
	attributeM.Attrvalue = attrvalue
	res := service.InitDB().Updates(&attributeM)
	if res.Error == nil {
		ac.TopjuiSucess(ctx, "编辑成功")
	} else {
		ac.TopjuiError(ctx, "编辑失败")
	}
}

func (ac *AttributeController) PostForever_del(ctx iris.Context) {
	id, _ := strconv.Atoi(ctx.PostValue("id"))
	if id == 0 {
		ac.TopjuiError(ctx, "参数错误")
	}
	model_ename := ctx.URLParam("ename")
	attributeM := model.Attribute{}
	var res *gorm.DB
	if model_ename == "" {
		res = service.InitDB().Where(map[string]interface{}{"id": id}).Delete(&attributeM)
	} else {
		res = service.InitDB().Where(map[string]interface{}{"model_ename": model_ename}).Delete(&attributeM)
	}
	if res.Error == nil {
		ac.TopjuiSucess(ctx, "删除成功")
	} else {
		ac.TopjuiError(ctx, "删除失败")
	}
}

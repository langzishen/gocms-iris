package controller

import (
	"github.com/kataras/iris/v12"
	"gocms/application/model"
	"strconv"
)

type PhotoController struct {
	BaseController
}

func (pc *PhotoController) PostList(ctx iris.Context) {
	page, _ := strconv.Atoi(ctx.PostValue("page"))
	rows, _ := strconv.Atoi(ctx.PostValue("rows"))
	limit_start := (page - 1) * rows

	search := ctx.PostValue("search")
	search_str := ""
	if search != "" { //topjui prompt 有BUG
		search_str = "`title` like \"%" + search + "%\""
	}

	type photoM struct {
		model.Photo
		TTitle  string `json:"t_title"`
		Creater string `json:"creater"`
	}

	list := []photoM{}
	if search != "" {
		model.InitDB().Where(search_str).Offset(limit_start).Limit(rows).Find(&list)
	} else {
		model.InitDB().Offset(limit_start).Limit(rows).Find(&list)
	}
	var count int64
	if search != "" {
		model.InitDB().Model(model.Photo{}).Where(search_str).Count(&count)
	} else {
		model.InitDB().Model(model.Photo{}).Count(&count)
	}
	var pages int
	if len(list)%rows == 0 {
		pages = len(list) / rows
	} else {
		pages = len(list)/rows + 1
	}

	//var list2 []articleM
	for i, vo := range list {
		if vo.Tid != "" {
			categoryM := model.Category{}
			model.InitDB().Model(categoryM).Where(map[string]interface{}{"classid": vo.Tid}).Find(&categoryM)
			//vo.TTitle = categoryM.Classtitle
			//使用(&list[i])不用引入第三方变量
			(&list[i]).TTitle = categoryM.Classtitle
			userM := model.User{}
			model.InitDB().Model(userM).Where(map[string]interface{}{"id": vo.CreaterId}).Find(&userM)
			(&list[i]).Creater = userM.Nickname
		}
		//list2 = append(list2, vo)
	}
	//fmt.Printf("%+v", list2)
	ctx.JSON(map[string]interface{}{"pages": pages, "total": count, "rows": list})
}

func (pc *PhotoController) PostZree(ctx iris.Context) {
	classpid, _ := strconv.Atoi(ctx.URLParam("classpid"))
	category_lst := new(model.Category).ZreeList("1", classpid, map[string]interface{}{"classmodule": "photo"})
	ctx.JSON(category_lst)
}

func (pc *PhotoController) PostAdd(ctx iris.Context) {
	photoM := model.Photo{}
	photoM.Tid = ctx.PostValue("tid")
	photoM.Title = ctx.PostValue("title")
	photoM.Img = ctx.PostValue("img")
	photoM.CreaterId = uint(AuthId)
	photoM.Status, _ = strconv.Atoi(ctx.PostValue("status"))
	res := model.InitDB().Create(&photoM)
	if res.Error != nil {
		pc.TopjuiError(ctx, res.Error.Error())
	} else {
		pc.TopjuiSucess(ctx, "新增成功")
	}
}

func (pc *PhotoController) GetEdit(ctx iris.Context) {
	id, _ := strconv.Atoi(ctx.URLParam("id"))
	if id == 0 {
		pc.TopjuiError(ctx, "参数Id丢失")
	}
	photoM := model.Photo{}
	model.InitDB().Where(map[string]interface{}{"id": id}).Find(&photoM)

	ctx.ViewData("info", photoM)
	pc.BaseController.GetEdit(ctx)
}

func (pc *PhotoController) PostEdit(ctx iris.Context) {
	id, _ := strconv.Atoi(ctx.PostValue("id"))
	if id == 0 {
		pc.TopjuiError(ctx, "参数Id丢失")
	}
	photoM := model.Photo{}
	photoM.Tid = ctx.PostValue("tid")
	photoM.Title = ctx.PostValue("title")
	photoM.Img = ctx.PostValue("img")
	photoM.Status, _ = strconv.Atoi(ctx.PostValue("status"))
	res := model.InitDB().Where(map[string]interface{}{"id": id}).Updates(&photoM)
	if res.Error != nil {
		pc.TopjuiError(ctx, res.Error.Error())
	} else {
		pc.TopjuiSucess(ctx, "保存成功")
	}
}

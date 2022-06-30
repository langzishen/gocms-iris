package controller

import (
	"encoding/json"
	"github.com/kataras/iris/v12"
	"gocms/application/model"
	"strconv"
	"strings"
)

type ArticleController struct {
	BaseController
}

func (ac *ArticleController) PostList(ctx iris.Context) {
	page, _ := strconv.Atoi(ctx.PostValue("page"))
	rows, _ := strconv.Atoi(ctx.PostValue("rows"))
	limit_start := (page - 1) * rows

	search := ctx.PostValue("search")
	search_str := ""
	if search != "" { //topjui prompt 有BUG
		search_str = "`title` like \"%" + search + "%\""
	}

	type articleM struct {
		model.Article
		TTitle  string `json:"t_title"`
		Creater string `json:"creater"`
	}

	list := []articleM{}
	if search != "" {
		model.InitDB().Where(search_str).Offset(limit_start).Limit(rows).Find(&list)
	} else {

		model.InitDB().Offset(limit_start).Limit(rows).Find(&list)
	}
	var count int64
	if search != "" {
		model.InitDB().Model(model.Article{}).Where(search_str).Count(&count)
	} else {
		model.InitDB().Model(model.Article{}).Count(&count)
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

func (ac *ArticleController) GetAdd(ctx iris.Context) {
	attrlist := new(model.Attribute).GetList(map[string]interface{}{"model_ename": "article"})
	ctx.ViewData("attrlist", attrlist)
	ac.BaseController.GetAdd(ctx)
}

func (ac *ArticleController) PostZree(ctx iris.Context) {
	classpid, _ := strconv.Atoi(ctx.URLParam("classpid"))
	category_lst := new(model.Category).ZreeList("1", classpid, map[string]interface{}{"classmodule": "article"})
	ctx.JSON(category_lst)
}

func (ac *ArticleController) PostAdd(ctx iris.Context) {
	articleM := model.Article{}
	articleM.Tid = ctx.PostValue("tid")
	articleM.Title = ctx.PostValue("title")
	articleM.Keywords = ctx.PostValue("keywords")
	articleM.Description = ctx.PostValue("description")
	articleM.Img = ctx.PostValue("img")
	articleM.Content = ctx.PostValue("content")
	articleM.CreaterId = AuthId
	articleM.Sort, _ = strconv.Atoi(ctx.PostValue("sort"))
	articleM.Apv = 0
	articleM.Status, _ = strconv.Atoi(ctx.PostValue("status"))
	Attrtj_str, _ := ctx.PostValues("attrtj")
	articleM.Attrtj = strings.Join(Attrtj_str, ",")
	IsLoginShow, _ := strconv.Atoi(ctx.PostValue("is_login_show"))
	articleM.IsLoginShow = uint(IsLoginShow)
	res := model.InitDB().Create(&articleM)
	if res.Error != nil {
		ac.TopjuiError(ctx, res.Error.Error())
	} else {
		ac.TopjuiSucess(ctx, "新增成功")
	}
}

func (ac *ArticleController) GetEdit(ctx iris.Context) {
	id, _ := strconv.Atoi(ctx.URLParam("id"))
	if id == 0 {
		ac.TopjuiError(ctx, "参数Id丢失")
	}
	articleM := model.Article{}
	model.InitDB().Where(map[string]interface{}{"id": id}).Find(&articleM)

	attrtj := strings.Split(articleM.Attrtj, ",")
	attrlist := new(model.Attribute).GetList(map[string]interface{}{"model_ename": "article"})
	var attrlist2 []map[string]interface{}
	if len(attrlist) > 0 {
		for _, vo := range attrlist {
			map_byte, _ := json.Marshal(vo)
			map_vo := map[string]interface{}{}
			json.Unmarshal(map_byte, &map_vo)
			map_vo["is_check"] = IsInArr(attrtj, vo.Attrvalue)
			attrlist2 = append(attrlist2, map_vo)
		}
	}
	ctx.ViewData("attrlist", attrlist2)

	ctx.ViewData("info", articleM)
	ac.BaseController.GetEdit(ctx)
}

func (ac *ArticleController) PostEdit(ctx iris.Context) {
	id, _ := strconv.Atoi(ctx.PostValue("id"))
	if id == 0 {
		ac.TopjuiError(ctx, "参数Id丢失")
	}
	articleM := model.Article{}
	articleM.Tid = ctx.PostValue("tid")
	articleM.Title = ctx.PostValue("title")
	articleM.Keywords = ctx.PostValue("keywords")
	articleM.Description = ctx.PostValue("description")
	articleM.Img = ctx.PostValue("img")
	articleM.Content = ctx.PostValue("content")
	articleM.Sort, _ = strconv.Atoi(ctx.PostValue("sort"))
	articleM.Status, _ = strconv.Atoi(ctx.PostValue("status"))
	Attrtj_str, _ := ctx.PostValues("attrtj")
	articleM.Attrtj = strings.Join(Attrtj_str, ",")
	IsLoginShow, _ := strconv.Atoi(ctx.PostValue("is_login_show"))
	articleM.IsLoginShow = uint(IsLoginShow)
	res := model.InitDB().Where(map[string]interface{}{"id": id}).Updates(&articleM)
	if res.Error != nil {
		ac.TopjuiError(ctx, res.Error.Error())
	} else {
		ac.TopjuiSucess(ctx, "修改成功")
	}
}

func (ac *ArticleController) GetInfo(ctx iris.Context) {
	id, _ := strconv.Atoi(ctx.URLParam("id"))
	title := ctx.URLParam("title")
	info := map[string]interface{}{"id": id, "title": title}
	ctx.ViewData("info", info)
	ctx.View(RequestController + "/info.html")
}

func IsInArr(arr []string, str string) (is bool) {
	is = false
	for _, vo := range arr {
		if vo == str {
			is = true
			return is
		}
	}
	return is
}

func IntIsInArr(arr []int, i int) (is bool) {
	is = false
	for _, vo := range arr {
		if vo == i {
			is = true
			return is
		}
	}
	return is
}

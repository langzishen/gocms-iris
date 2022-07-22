package controller

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"go/importer"
	"gocms/application/app_session"
	"gocms/application/model"
	"gocms/application/service"
	"gocms/config"
)

type IndexController struct {
	BaseController
}

/**
后台首页
*/
func (ic *IndexController) GetIndex(ctx iris.Context) {
	ctx.ViewData("topmenujson", ic.TopMenu())
	ctx.ViewData("authId", AuthId)
	ctx.ViewData("authNickName", AuthNickname)
	conf := config.InitConfig()
	ctx.ViewData("app_icp", conf.SiteIcpNum)
	ctx.ViewData("app_name", conf.AppName)
	ctx.ViewData("site_keywords", conf.SiteKeywords)
	ctx.ViewData("site_description", conf.SiteDescription)
	ctx.ViewData("app_version", conf.AppVersion)

	ctx.View("index/index.html")
}

/**
后台主标签页
*/
func (ic *IndexController) GetHome(ctx iris.Context) {
	ctx.ViewData("message_count", "123,456")
	ctx.ViewData("message_today_count", "12,345")
	ctx.ViewData("pv_count", "3,456,789")
	ctx.ViewData("pv_day_count", "1,000")
	ctx.ViewData("user_count", "1,234")
	ctx.ViewData("user_day_count", "123")
	ctx.ViewData("article_count", "456,789")
	ctx.ViewData("article_day_count", "9,876")

	new_article_list := []map[string]interface{}{
		{"id": 1, "title": "最新文章1"},
		{"id": 2, "title": "最新文章2"},
		{"id": 3, "title": "最新文章3"},
		{"id": 4, "title": "最新文章4"},
		{"id": 5, "title": "最新文章5"},
	}
	hot_article_list := []map[string]interface{}{
		{"id": 1, "title": "热门文章1"},
		{"id": 2, "title": "热门文章2"},
		{"id": 3, "title": "热门文章3"},
		{"id": 4, "title": "热门文章4"},
		{"id": 5, "title": "热门文章5"},
	}
	notice_article_list := []map[string]interface{}{
		{"id": 1, "title": "通知公告1"},
		{"id": 2, "title": "通知公告2"},
		{"id": 3, "title": "通知公告3"},
		{"id": 4, "title": "通知公告4"},
		{"id": 5, "title": "通知公告5"},
	}

	ctx.ViewData("new_article_list", new_article_list)
	ctx.ViewData("hot_article_list", hot_article_list)
	ctx.ViewData("notice_article_list", notice_article_list)
	ctx.View("index/home.html")
}

func (ic *IndexController) GetIndex2(ctx iris.Context) {
	pkg2, err2 := importer.Default().Import("github.com/kataras/iris")
	fmt.Println(pkg2, err2)
	pkg, err := importer.Default().Import("fmt")
	fmt.Println(pkg, err)
}

func (ic *IndexController) GetIndex3(ctx iris.Context) {
	list := []model.Article{}
	where := make(map[string]map[string]interface{})
	where["eq"] = map[string]interface{}{"id": 1}
	where["gte"] = map[string]interface{}{"status": 0}
	where["like"] = map[string]interface{}{"title": "%文章%"}
	list_json, total_count := new(service.Base).List(&list, []string{"id", "title"}, where, "", 0, 10)
	fmt.Println(total_count)
	fmt.Println("--------------------------------------")
	fmt.Println(list)
	fmt.Println("++++++++++++++++++++++++++++++++++++")
	fmt.Println(list_json)
}

func (ic *IndexController) GetLogout(ctx iris.Context) {
	app_session.Sess.Start(ctx).Delete("authId")
	/*
		app_session.Sess.Start(ctx).Delete("email")
		app_session.Sess.Start(ctx).Delete("loginUserName")

		app_session.Sess.Start(ctx).Delete("type_id")
		appconf := config.InitConfig()
		app_session.Sess.Start(ctx).Delete(appconf.Rbac.AdminAuthKey)**/
	app_session.Sess.Start(ctx).Destroy()
	ctx.Redirect("/" + RequestApp + "/public/login")
}

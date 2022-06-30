package controller

import (
	"github.com/kataras/iris/v12"
	"gocms/application/admin/extend/rbac"
	"gocms/application/app_session"
	"gocms/application/extend/upload"
	"gocms/application/model"
	"gocms/config"
	"strconv"
	"strings"
)

var (
	RequestApp        string //访问的当前app
	RequestController string //访问的当前控制器
	RequestAction     string //访问的当前操作
	AuthId            int    //登录者Id
	AuthNickname      string //登陆者昵称
	AuthTypeId        int    //登录者类型，9是超级管理员
	IsAdministrator   bool   //是否是超级管理员
	SearchTids        []int  //数据授权的tids
)

type BaseController struct {
}

/**
初始化
*/
func initialize(ctx iris.Context) {
	ctx.ViewData("RequestApp", RequestApp)
	ctx.ViewData("RequestController", RequestController)
	ctx.ViewData("RequestAction", RequestAction)

	Rbac(ctx)
	SearchTids = DataAccess(ctx)
}

/**
权限认证
*/
func Rbac(ctx iris.Context) {
	conf := config.InitConfig()
	no_auth_c := strings.Split(conf.Rbac.NoAuthController, ",")
	no_auth_a := strings.Split(conf.Rbac.NoAuthAction, ",")
	if !IsInArr(no_auth_c, RequestController) {
		if !IsInArr(no_auth_a, RequestController+"-"+RequestAction) {
			if !new(rbac.RBAC).AccessDecision(RequestApp, RequestController, RequestAction, ctx) {
				TopjuiCtxEnd(ctx)
			}
		}
	}
}

/**
增强型RBAC-数据授权
	请求中带`tid`参数时为操作，将判断该tid是否有数据权限，没有tid时，返回tid的[]int 用于list查询
*/
func DataAccess(ctx iris.Context) (tids []int) {
	conf := config.InitConfig()

	if !IsAdministrator { //不是超级管理员
		if conf.DataAccess.IsOpen == 1 { //开启了数据授权
			data_access_controllers := strings.Split(conf.DataAccess.Controllers, ",")
			if IsInArr(data_access_controllers, RequestController) {

				role_userM := []model.RoleUser{}
				role_id_s := []int{}
				model.InitDB().Where(map[string]interface{}{"user_id": AuthId}).Find(&role_userM)
				if len(role_userM) > 0 {
					for _, role_user_v := range role_userM {
						role_id_s = append(role_id_s, int(role_user_v.RoleId))
					}
				}

				data_access_map := map[string]interface{}{
					"role_id": role_id_s,
					"model":   RequestController,
				}
				if ctx.URLParam("node_id") == "" {
					data_access_map["node_id"] = ctx.URLParam("node_id")
				}
				if ctx.PostValue("node_id") == "" {
					data_access_map["node_id"] = ctx.PostValue("node_id")
				}
				data_accessM := []model.DataAccess{}
				model.InitDB().Where(data_access_map).Find(&data_accessM)
				if len(data_accessM) > 0 {
					for _, data_access_v := range data_accessM {
						tids = append(tids, int(data_access_v.Tid))
					}
				}

				request_tid := 0
				if ctx.URLParam("tid") != "" {
					request_tid, _ = strconv.Atoi(ctx.URLParam("tid"))
				}
				if ctx.PostValue("tid") != "" {
					request_tid, _ = strconv.Atoi(ctx.PostValue("tid"))
				}
				if request_tid != 0 {
					if !IntIsInArr(tids, request_tid) {
						TopjuiCtxEnd(ctx)
					}
				}
				return tids
			}
		}
	}
	return []int{}
}

func TopjuiCtxEnd(ctx iris.Context) {
	if ctx.IsAjax() {
		ctx.Header("Content-Type", "application/json; charset=utf-8") //输出JSON头信息
		ctx.JSON(map[string]interface{}{"statusCode": 500, "status": 0, "message": "权限不足!", "closeCurrent": true})
		ctx.StopExecution()
	} else {
		//ctx.View("public/no_access.html")
		body := `<div class="topjui-fluid">
						<br>
    					<div class="topjui-row" >
       	 					无访问权限
    					</div>
					</div>`
		ctx.WriteString(body)
		ctx.StopExecution()
	}
}

/**
输出头部大菜单
*/
func (bc *BaseController) TopMenu() []map[string]interface{} {
	appConf := config.InitConfig()
	menu_map_content := map[string]interface{}{"iconCls": appConf.Admin.Menu.Content.Icons, "text": appConf.Admin.Menu.Content.Title, "codeSetId": "menu", "resourceType": "menu", "levelId": 1, "id": "content", "pid": 0, "state": "closed", "url": "", "textColour": ""}
	menu_map_framework := map[string]interface{}{"iconCls": appConf.Admin.Menu.Framework.Icons, "text": appConf.Admin.Menu.Framework.Title, "codeSetId": "menu", "resourceType": "menu", "levelId": 1, "id": "framework", "pid": 0, "state": "closed", "url": "", "textColour": ""}
	menu_map_groupuser := map[string]interface{}{"iconCls": appConf.Admin.Menu.Groupuser.Icons, "text": appConf.Admin.Menu.Groupuser.Title, "codeSetId": "menu", "resourceType": "menu", "levelId": 1, "id": "groupuser", "pid": 0, "state": "closed", "url": "", "textColour": ""}
	menu_map_developers := map[string]interface{}{"iconCls": appConf.Admin.Menu.Developers.Icons, "text": appConf.Admin.Menu.Developers.Title, "codeSetId": "menu", "resourceType": "menu", "levelId": 1, "id": "developers", "pid": 0, "state": "closed", "url": "", "textColour": ""}
	menu_map_system := map[string]interface{}{"iconCls": appConf.Admin.Menu.System.Icons, "text": appConf.Admin.Menu.System.Title, "codeSetId": "menu", "resourceType": "menu", "levelId": 1, "id": "system", "pid": 0, "state": "closed", "url": "", "textColour": ""}
	menu_map_plugin := map[string]interface{}{"iconCls": appConf.Admin.Menu.Plugin.Icons, "text": appConf.Admin.Menu.Plugin.Title, "codeSetId": "menu", "resourceType": "menu", "levelId": 1, "id": "plugin", "pid": 0, "state": "closed", "url": "", "textColour": ""}
	var top_menu []map[string]interface{}
	top_menu = append(top_menu, menu_map_content, menu_map_framework, menu_map_groupuser, menu_map_developers, menu_map_system, menu_map_plugin)
	return top_menu
}

/**
输出左侧某大菜单下的全部控制器和操作
*/
func (bc *BaseController) LeftMenu(ctx iris.Context) []map[string]interface{} {
	top_menu_item := ctx.URLParam("menu_id")
	node_list := []*model.Node{}
	var left_menu []map[string]interface{}
	model.InitDB().Where("level=2 AND status=1 and group_id=?", top_menu_item).Order("sort").Find(&node_list)
	for _, node_vo := range node_list {
		node_list2 := []*model.Node{}
		model.InitDB().Where("level=3 AND status=1 AND left_menu_action=1 and pid=?", node_vo.Id).Order("sort").Find(&node_list2)
		i := 0
		for _, node_vo2 := range node_list2 {
			if new(rbac.RBAC).AccessDecision(RequestApp, node_vo.Name, node_vo2.Name, ctx) {
				i++
				left_menu = append(left_menu, map[string]interface{}{"iconCls": node_vo2.Icon, "text": node_vo2.Title, "codeSetId": top_menu_item, "resourceType": "menu", "levelId": 3, "id": int(node_vo2.Id), "pid": int(node_vo2.Pid), "state": "open", "url": "/" + RequestApp + "/" + node_vo.Name + "/" + node_vo2.Name, "textColour": ""})
			}
		}
		if i > 0 { //如果控制器下所有的操作都没有权限就不显示这个控制器
			left_menu = append(left_menu, map[string]interface{}{"iconCls": node_vo.Icon, "text": node_vo.Title, "codeSetId": top_menu_item, "resourceType": "menu", "levelId": 2, "id": int(node_vo.Id), "pid": int(node_vo.Pid), "state": "closed", "url": "", "textColour": ""})
		}
	}
	return left_menu
}

/**
输出左侧单个大菜单下的菜单
*/
func (bc *BaseController) GetMenu(ctx iris.Context) {
	ctx.JSON(bc.LeftMenu(ctx))
}

/**
全局获取列表页的视图
*/
func (bc *BaseController) GetIndex(ctx iris.Context) {
	ctx.View(RequestController + "/index.html")
}

/**
搜索左侧菜单栏
*/
func (bc *BaseController) PostSearchleftmenu(ctx iris.Context) {
	search := ctx.PostValue("text")
	if search == "" {
		return
	}
	node_list := []*model.Node{}
	var left_menu []map[string]interface{}
	model.InitDB().Where("level=2 AND status=1").Order("sort").Find(&node_list)
	for _, node_vo := range node_list {
		node_list2 := []*model.Node{}
		model.InitDB().Debug().Where("level=3 AND status=1 AND left_menu_action=1 and pid=? AND title like '%"+search+"%'", node_vo.Id).Order("sort").Find(&node_list2)
		for _, node_vo2 := range node_list2 {
			if new(rbac.RBAC).AccessDecision(RequestApp, node_vo.Name, node_vo2.Name, ctx) {
				left_menu = append(left_menu, map[string]interface{}{"text": node_vo2.Title, "url": "/" + RequestApp + "/" + node_vo.Name + "/" + node_vo2.Name})
			}
		}
	}
	ctx.JSON(left_menu)
}

/**
全局添加视图
*/
func (bc *BaseController) GetAdd(ctx iris.Context) {
	ctx.View(RequestController + "/add.html")
}

/**
全局修改视图
*/
func (bc *BaseController) GetEdit(ctx iris.Context) {
	id_s := ctx.URLParam("id")
	if id_s == "" {
		bc.TopjuiError(ctx, "参数Id丢失")
	}

	ctx.View(RequestController + "/edit.html")
}

/**
通用删除byID
*/
func (bc *BaseController) PostDel(ctx iris.Context) {
	bc.AjaxDone(ctx, "del")
}

func (bc *BaseController) PostResume(ctx iris.Context) {
	bc.AjaxDone(ctx, "resume")
}

func (bc *BaseController) PostForbid(ctx iris.Context) {
	bc.AjaxDone(ctx, "forbid")
}

func (bc *BaseController) PostRecycle(ctx iris.Context) {
	bc.AjaxDone(ctx, "recycle")
}

/**
彻底删除
*/
func (bc *BaseController) PostForever_del(ctx iris.Context) {
	id := ctx.PostValue("id")
	if id == "" {
		bc.TopjuiError(ctx, "删除操作参数Id必须")
	}
	var id_arr []int
	for _, id_str := range strings.Split(id, ",") {
		id_i, _ := strconv.Atoi(id_str)
		id_arr = append(id_arr, id_i)
	}
	res := model.InitDB().Exec("DELETE FROM `"+RequestController+"` WHERE id IN ?", id_arr)
	if res.Error == nil {
		bc.TopjuiSucess(ctx, "删除成功")
	} else {
		bc.TopjuiError(ctx, "删除失败")
	}
}

/**
审核、禁用、删除、恢复
*/
func (bc *BaseController) AjaxDone(ctx iris.Context, do string) {
	id := ctx.PostValue("id")
	if id == "" {
		bc.TopjuiError(ctx, "删除操作参数Id必须")
	}
	var status int
	var say string
	switch do {
	case "del": //删除      1->-1
		status = -1
		say = "删除"
	case "resume": //审核   0->1
		status = 1
		say = "审核"
	case "forbid": //禁用   1->0
		status = 0
		say = "禁用"
	case "recycle": //恢复   -1->0
		status = 1
		say = "恢复"
	}
	var id_arr []int
	for _, id_str := range strings.Split(id, ",") {
		id_i, _ := strconv.Atoi(id_str)
		id_arr = append(id_arr, id_i)
	}

	res := model.InitDB().Exec("UPDATE `"+RequestController+"` SET status="+strconv.Itoa(status)+" WHERE id IN ?", id_arr)
	if res.Error == nil {
		bc.TopjuiSucess(ctx, say+"成功")
	} else {
		bc.TopjuiError(ctx, say+"失败")
	}
}

func (bc *BaseController) GetMenu_zree(ctx iris.Context) {
	node_list := []*model.Node{}
	model.InitDB().Where("level=3 AND status=1 and pid=?", ctx.URLParam("pid")).Order("sort").Find(&node_list)
	nodeZree_list := []map[string]interface{}{}
	for _, v := range node_list {
		node := model.Node{}
		model.InitDB().Where("id=?", v.Pid).First(&node)
		nodeZree_list = append(nodeZree_list, map[string]interface{}{
			"id":      v.Id,
			"pid":     v.Pid,
			"text":    v.Title,
			"state":   "close",
			"iconCls": "",
			"url":     "/" + RequestApp + "/" + strings.ToLower(node.Name) + "/" + v.Name})
	}
	ctx.JSON(nodeZree_list)
}

/**
通用上传
*/
func (b *BaseController) PostAjax_upload(ctx iris.Context) {
	var upload_file_name string
	if ctx.PostValue("file") != "" {
		upload_file_name = ctx.PostValue("file")
	} else {
		upload_file_name = "file"
	}
	img, _, err := upload.Upload(ctx, upload_file_name) //topjui的单文件上传的文件名默认为file，与input的name无关
	if err != nil {
		ctx.JSON(map[string]interface{}{"statusCode": 500, "title": "操作提示", "message": err, "filePath": ""})
	} else {
		ctx.JSON(map[string]interface{}{"statusCode": 200, "title": "操作提示", "message": "", "filePath": img})
	}
	ctx.StopExecution()
}

func (b *BaseController) TopjuiSucess(ctx iris.Context, message string) {
	var navTabTid, forwardUrl string
	if ctx.Method() == "GET" {
		navTabTid = ctx.URLParam("navTabId")
		forwardUrl = ctx.URLParam("forwardUrl")
	} else {
		navTabTid = ctx.PostValue("navTabId")
		forwardUrl = ctx.PostValue("forwardUrl")
	}
	b.TopjuiReJSON(ctx, 1, message, navTabTid, forwardUrl, true)
	ctx.StopExecution() //结束ctx向下执行
}

func (b *BaseController) TopjuiError(ctx iris.Context, message string) {
	var navTabTid, forwardUrl string
	if ctx.Method() == "GET" {
		navTabTid = ctx.URLParam("navTabId")
		forwardUrl = ctx.URLParam("forwardUrl")
	} else {
		navTabTid = ctx.PostValue("navTabId")
		forwardUrl = ctx.PostValue("forwardUrl")
	}
	b.TopjuiReJSON(ctx, 0, message, navTabTid, forwardUrl, false)
	ctx.StopExecution() //结束ctx向下执行
}

/**
 *  code返回的状态码
 */
func (b *BaseController) TopjuiReJSON(ctx iris.Context, code int, message string, navTabTid string, forwardUrl string, callbackType bool) {
	res := make(map[string]interface{})
	res["status"] = code
	if code == 1 {
		res["statusCode"] = 200
	} else {
		if code == 0 {
			res["statusCode"] = 300
		} else {
			res["statusCode"] = code
		}
	}
	res["message"] = message
	res["tabid"] = navTabTid
	res["forwardUrl"] = forwardUrl
	res["closeCurrent"] = callbackType
	ctx.Header("Content-Type", "application/json; charset=utf-8") //输出JSON头信息
	ctx.StopWithJSON(200, res)
}

/**
后台登录权限验证中间件
*/
func AdminAuth(ctx iris.Context) {
	app_conf := config.InitConfig()
	app_session.Sess.UseDatabase(app_session.Sessdb) //开启session的本地储存，防止长期服务时session丢失
	session := app_session.Sess.Start(ctx)

	pathArr := strings.Split(ctx.Path(), "/")
	RequestApp = pathArr[1]        //app
	RequestController = pathArr[2] //controller
	RequestAction = pathArr[3]     //action

	if authId, err := session.GetInt(app_conf.Rbac.UserAuthKey); err == nil && authId != -1 {
		AuthId = authId
		AuthTypeId, _ = session.GetInt("type_id")
		if AuthTypeId == 9 { //超级管理员
			AuthNickname = session.GetString("loginUserName") + "(超级管理员)"
			IsAdministrator = true
		} else {
			AuthNickname = session.GetString("loginUserName")
			IsAdministrator = false
		}
		initialize(ctx)
		ctx.Next()
	} else { //没有session中的authId
		if (RequestController == "public" && RequestAction == "login") || (RequestController == "public" && RequestAction == "loginin") || (RequestController == "public" && RequestAction == "captcha") { //登录操作不做权限验证
			initialize(ctx)
			ctx.Next() //继续向下运行
		} else {
			//跳转到登录
			//ctx.Writef("no login")
			ctx.Redirect("/" + RequestApp + "/public/login")
		}
	}
}

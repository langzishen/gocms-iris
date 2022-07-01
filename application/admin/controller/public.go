package controller

import (
	"crypto/md5"
	"fmt"
	"github.com/kataras/iris/v12"
	"gocms/application/admin/extend/rbac"
	"gocms/application/app_session"
	"gocms/application/extend/captcha"
	"gocms/application/model"
	"gocms/config"
	"strings"
)

type PublicController struct {
	BaseController
}

func (_ *PublicController) GetLogin(ctx iris.Context) {
	app_conf := config.InitConfig()
	if AuthId, _ = app_session.Sess.Start(ctx).GetInt(app_conf.Rbac.UserAuthKey); AuthId != -1 { //有session中的authId跳转到登录页面
		ctx.Redirect("/" + RequestApp + "/index/index")
	}
	ctx.View("public/login.html")
}

func (pc *PublicController) PostLoginin(ctx iris.Context) {
	if ctx.IsAjax() {
		data := make(map[string]interface{})
		ctx.ReadJSON(&data)
		if "" == data["account"] {
			ctx.JSON(map[string]interface{}{"code": 403, "msg": "账号必填!", "data": "", "referer": "/" + RequestApp + "/index/index"})
		}
		if "" == data["password"] {
			ctx.JSON(map[string]interface{}{"code": 403, "msg": "账号必填!", "data": "", "referer": "/" + RequestApp + "/index/index"})
		}
		if "" == data["captcha"].(string) {
			ctx.JSON(map[string]interface{}{"code": 403, "msg": "请输入验证码!", "data": "", "referer": "/" + RequestApp + "/index/index"})
		}
		sess := app_session.Sess.Start(ctx)
		code := sess.GetString("admin_login_captcha")
		if code != "" {
			if strings.ToUpper(code) != strings.ToUpper(data["captcha"].(string)) {
				//ctx.JSON(map[string]interface{}{"code": 403, "msg": "验证码不正确!", "data": "", "referer": "/" + RequestApp + "/public/login"})
				//ctx.StopExecution()
				//测试时不做验证码验证
				//return
			}
		}

		where_map := make(map[string]interface{})
		where_map["account"] = data["account"]
		if isOk, userInfo := new(rbac.RBAC).Authenticate(where_map); !isOk {
			ctx.JSON(map[string]interface{}{"code": 403, "msg": "帐号不存在或已禁用！", "data": "", "referer": ""})
		} else {
			if fmt.Sprintf("%x", md5.Sum([]byte(data["password"].(string)))) != userInfo.Password {
				ctx.JSON(map[string]interface{}{"code": 403, "msg": "用户名和密码不正确！", "data": "", "referer": ""})
			} else {
				pc.saveLoginSession(ctx, userInfo)              // 设置后台的SESSION
				pc.setLoginLog("boss", userInfo.Id)             // 保存登录信息
				new(rbac.RBAC).SaveAccessList(ctx, userInfo.Id) // 缓存访问权限
				ctx.JSON(map[string]interface{}{"code": 200, "msg": "", "data": "", "referer": "/" + RequestApp + "/index/index"})
			}
		}
	} else {
		ctx.Redirect("/" + RequestApp + "/index/index")
	}
}

func (pc *PublicController) GetCaptcha(ctx iris.Context) {
	sess := app_session.Sess.Start(ctx)
	code := sess.GetString("admin_login_captcha")
	if code != "" {
		sess.Delete("admin_login_captcha")
	}

	code = captcha.NewCaptcha(ctx, 120, 40)
	sess.Set("admin_login_captcha", code)
}

/**
 * 设置后台的SESSION
 */
func (pc *PublicController) saveLoginSession(ctx iris.Context, userInfo model.User) {
	session := app_session.Sess.Start(ctx)
	appconf := config.InitConfig()
	session.Set(appconf.Rbac.UserAuthKey, int(userInfo.Id))
	new(rbac.RBAC).SaveAccessList(ctx, userInfo.Id) //设置RBAC权限
	session.Set("email", userInfo.Email)
	session.Set("loginUserName", userInfo.Nickname)
	session.Set("type_id", userInfo.TypeId)
	if userInfo.TypeId == 9 { //9是超级管理员
		session.Set(appconf.Rbac.AdminAuthKey, true)
	}
}

/**
 * 保存后台的登录信息
 */
func (pc *PublicController) setLoginLog(appName string, authId uint) {
	new(model.LogLogin).AddLogLogin(appName, authId)
}

package controller

import (
	"crypto/md5"
	"fmt"
	"github.com/kataras/iris/v12"
	"gocms/application/model"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type UserController struct {
	BaseController
}

/**
列表
*/
func (uc *UserController) PostList(ctx iris.Context) {
	page, _ := strconv.Atoi(ctx.PostValue("page"))
	rows, _ := strconv.Atoi(ctx.PostValue("rows"))
	limit_start := (page - 1) * rows

	search := ctx.PostValue("search_value")
	search_str := "`id`>1" //排除admin超级管理员
	if search != "" {      //topjui prompt 有BUG
		search_str += " and `" + ctx.PostValue("search_key") + "` like \"%" + search + "%\""
	}

	type user struct {
		model.User
		LastLoginTime int    `json:"last_login_time"`
		LastLoginIp   string `json:"last_login_ip"`
		LoginCount    int64  `json:"login_count"`
	}

	list := []*user{}
	model.InitDB().Where(search_str).Offset(limit_start).Limit(rows).Find(&list)
	var count int64
	model.InitDB().Model(model.User{}).Where(search_str).Count(&count)
	var pages int
	if len(list)%rows == 0 {
		pages = len(list) / rows
	} else {
		pages = len(list)/rows + 1
	}

	for _, vo := range list {
		vo.LastLoginTime, vo.LastLoginIp, vo.LoginCount = new(model.LogLogin).LastLogLogin("boss", vo.Id)
	}
	ctx.JSON(map[string]interface{}{"pages": pages, "total": count, "rows": list})
}

func (uc *UserController) PostAdd(ctx iris.Context) {
	userM := model.User{}
	userM.ObjectId = RandString(8)
	userM.Account = ctx.PostValue("account")
	userM.Nickname = ctx.PostValue("nickname")
	userM.Password = fmt.Sprintf("%x", md5.Sum([]byte(ctx.PostValue("password"))))
	userM.Phone = ctx.PostValue("phone")
	userM.Email = ctx.PostValue("email")
	userM.CreateTime = uint(time.Now().Unix())
	userM.UpdateTime = uint(time.Now().Unix())
	userM.Status, _ = strconv.Atoi(ctx.PostValue("status"))
	userM.TypeId = 1
	res := model.InitDB().Create(&userM)
	if res.Error == nil {
		uc.TopjuiSucess(ctx, "新增成功")
	} else {
		uc.TopjuiError(ctx, "新增失败")
	}
}

func (uc *UserController) GetEdit(ctx iris.Context) {
	id, _ := strconv.Atoi(ctx.URLParam("id"))
	if id == 0 {
		uc.TopjuiError(ctx, "参数Id必须")
	}
	userM := model.User{}
	model.InitDB().Find(&userM, id)
	ctx.ViewData("info", userM)
	ctx.View(RequestController + "/edit.html")
}

func (uc *UserController) PostEdit(ctx iris.Context) {
	id, _ := strconv.Atoi(ctx.PostValue("id"))
	if id == 0 {
		uc.TopjuiError(ctx, "参数Id必须")
	}
	userM := model.User{}
	userM.Id = uint(id)
	userM.Account = ctx.PostValue("account")
	userM.Nickname = ctx.PostValue("nickname")
	userM.Phone = ctx.PostValue("phone")
	userM.Email = ctx.PostValue("email")
	userM.UpdateTime = uint(time.Now().Unix())
	userM.Status, _ = strconv.Atoi(ctx.PostValue("status"))
	res := model.InitDB().Save(&userM)
	if res.Error == nil {
		uc.TopjuiSucess(ctx, "保存成功")
	} else {
		uc.TopjuiError(ctx, "保存失败")
	}
}

func (uc *UserController) GetResetpassword(ctx iris.Context) {
	id, _ := strconv.Atoi(ctx.URLParam("id"))
	if id == 0 {
		uc.TopjuiError(ctx, "参数Id必须")
	}
	ctx.ViewData("id", id)
	ctx.View(RequestController + "/resetpassword.html")
}

func (uc *UserController) PostResetpassword(ctx iris.Context) {
	id, _ := strconv.Atoi(ctx.PostValue("id"))
	if id == 0 {
		uc.TopjuiError(ctx, "参数Id必须")
	}
	userM := model.User{}
	userM.Id = uint(id)
	userM.Password = fmt.Sprintf("%x", md5.Sum([]byte(ctx.PostValue("password"))))
	res := model.InitDB().Save(&userM)
	if res.Error == nil {
		uc.TopjuiSucess(ctx, "重置成功")
	} else {
		uc.TopjuiError(ctx, "重置失败")
	}
}

func (uc *UserController) GetPersonal(ctx iris.Context) {
	userM := model.User{}
	model.InitDB().Where(map[string]interface{}{"id": AuthId}).Find(&userM)
	ctx.ViewData("info", userM)
	role_user_list := []model.RoleUser{}
	model.InitDB().Where(map[string]interface{}{"id": AuthId}).Find(&role_user_list)
	role_name := []string{}
	if AuthTypeId == 9 {
		role_name = append(role_name, "超级管理员")
	}
	if len(role_user_list) > 0 {
		role_id_arr := []int{}
		for _, role_user_vo := range role_user_list {
			role_id_arr = append(role_id_arr, int(role_user_vo.RoleId))
		}
		roleList := []model.Role{}
		model.InitDB().Where(map[string]interface{}{"id": role_id_arr}).Find(&roleList)
		for _, roleVo := range roleList {
			role_name = append(role_name, roleVo.Name)
		}
	}
	ctx.ViewData("role_name", strings.Join(role_name, ","))
	ctx.View(RequestController + "/personal.html")
}

func (uc *UserController) PostDopersonal(ctx iris.Context) {
	userM := model.User{}
	userM.Nickname = ctx.PostValue("nickname")
	sex, _ := strconv.Atoi(ctx.PostValue("sex"))
	userM.Sex = uint(sex)
	age, _ := strconv.Atoi(ctx.PostValue("age"))
	userM.Age = uint(age)
	userM.Logo = ctx.PostValue("logo")
	userM.Email = ctx.PostValue("email")
	userM.Phone = ctx.PostValue("phone")
	res := model.InitDB().Where(map[string]interface{}{"id": AuthId}).Updates(&userM)
	if res.Error != nil {
		uc.TopjuiError(ctx, "设置失败："+res.Error.Error())
	} else {
		uc.TopjuiSucess(ctx, "设置成功")
	}
}

func (uc *UserController) GetModifypassword(ctx iris.Context) {
	ctx.View("user/modifypassword.html")
}

func (uc *UserController) PostDomodifypassword(ctx iris.Context) {
	originalPassword := ctx.PostValue("originalPassword")
	password := ctx.PostValue("password")
	password2 := ctx.PostValue("password2")
	if password != password2 {
		uc.TopjuiError(ctx, "新密码两次不一致")
	}
	userM := model.User{}
	model.InitDB().Where(map[string]interface{}{"id": AuthId}).Find(&userM)
	if userM.Password != fmt.Sprintf("%x", md5.Sum([]byte(originalPassword))) {
		uc.TopjuiError(ctx, "原密码不正确,请联系系统管理员重置")
	}
	password = fmt.Sprintf("%x", md5.Sum([]byte(password)))
	res := model.InitDB().Model(userM).Where(map[string]interface{}{"id": AuthId}).Update("password", password)
	if res.Error != nil {
		uc.TopjuiError(ctx, "设置失败："+res.Error.Error())
	} else {
		uc.TopjuiSucess(ctx, "设置成功")
	}
}

func RandString(n int) string {
	allstring := "qwertyuiopasdfghjklzxcvbnm1233456789"
	ret := ""
	for i := 0; i < n; i++ {
		tmp := rand.Intn(len(allstring))
		ret += allstring[tmp : tmp+1]
	}
	return ret
}

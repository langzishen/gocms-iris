package controller

import (
	"encoding/json"
	"github.com/kataras/iris/v12"
	"gocms/application/model"
	"gocms/application/service"
	"gocms/config"
	"strconv"
	"strings"
	"time"
)

type RoleController struct {
	BaseController
}

func (rc *RoleController) PostList(ctx iris.Context) {
	page, _ := strconv.Atoi(ctx.PostValue("page"))
	rows, _ := strconv.Atoi(ctx.PostValue("rows"))
	limit_start := (page - 1) * rows

	search := ctx.PostValue("search")
	search_str := ""
	if search != "" { //topjui prompt 有BUG
		search_str = "`name` like \"%" + search + "%\""
	}

	type role_d struct {
		model.Role
		PName string `json:"p_name"`
	}
	list := []role_d{}
	if search != "" {
		service.InitDB().Debug().Where(search_str).Offset(limit_start).Limit(rows).Find(&list)
	} else {
		service.InitDB().Debug().Offset(limit_start).Limit(rows).Find(&list)
	}
	var count int64
	if search != "" {
		service.InitDB().Model(model.Role{}).Where(search_str).Count(&count)
	} else {
		service.InitDB().Model(model.Role{}).Count(&count)
	}
	var pages int
	if len(list)%rows == 0 {
		pages = len(list) / rows
	} else {
		pages = len(list)/rows + 1
	}

	var list2 []role_d
	for _, vo := range list {
		roleM := model.Role{}
		res := service.InitDB().Where(map[string]interface{}{"id": vo.Pid}).Last(&roleM)
		if res.Error != nil {
			vo.PName = "顶级分组"
		} else {
			vo.PName = roleM.Name
		}
		list2 = append(list2, vo)
	}

	ctx.JSON(map[string]interface{}{"pages": pages, "total": count, "rows": list2})
}

func (rc *RoleController) PostZree(ctx iris.Context) {
	pid_s := ctx.URLParam("pid")
	var pid int
	if pid_s == "" { //没有传pid默认顶级
		pid = 0
	} else {
		pid, _ = strconv.Atoi(pid_s)
	}
	roleM := []model.Role{}
	service.InitDB().Where(map[string]interface{}{"pid": pid}).Find(&roleM)
	var tree_list []map[string]interface{}
	if pid == 0 {
		tree_list = append(tree_list, map[string]interface{}{"id": 0, "pid": 0, "text": "无上级组"})
	}
	for _, vo := range roleM {
		roleM2 := []model.Role{}
		service.InitDB().Where(map[string]interface{}{"pid": vo.Id}).Find(&roleM2)
		if len(roleM2) > 0 { //map["state"]  state不可少，否则视图的 expandUrl 参数不起作用
			tree_list = append(tree_list, map[string]interface{}{"id": int(vo.Id), "state": "closed", "pid": vo.Pid, "text": vo.Name})
		} else { //没有state的树分支不能展开也不能向下访问expandUrl
			tree_list = append(tree_list, map[string]interface{}{"id": int(vo.Id), "pid": vo.Pid, "text": vo.Name})
		}
	}
	ctx.JSON(tree_list)
}

func (rc *RoleController) PostAdd(ctx iris.Context) {
	roleM := model.Role{}
	roleM.Pid, _ = strconv.Atoi(ctx.PostValue("pid"))
	roleM.Name = ctx.PostValue("name")
	roleM.Status, _ = strconv.Atoi(ctx.PostValue("status"))
	roleM.Ename = ""
	roleM.Remark = ctx.PostValue("remark")
	roleM.CreateTime = uint(time.Now().Unix())
	roleM.UpdateTime = uint(time.Now().Unix())
	res := service.InitDB().Create(&roleM)
	if res.Error == nil {
		rc.TopjuiSucess(ctx, "新增成功")
	} else {
		rc.TopjuiError(ctx, "新增失败")
	}
}

func (rc *RoleController) GetEdit(ctx iris.Context) {
	id_s := ctx.URLParam("id")
	if id_s == "" {
		rc.TopjuiError(ctx, "参数Id丢失")
	}
	id, _ := strconv.Atoi(id_s)
	roleM := model.Role{}
	service.InitDB().Where(map[string]interface{}{"id": id}).Last(&roleM)
	ctx.ViewData("info", &roleM)
	ctx.View(RequestController + "/edit.html")
}

func (rc *RoleController) PostEdit(ctx iris.Context) {
	id_s := ctx.PostValue("id")
	if id_s == "" {
		rc.TopjuiError(ctx, "参数Id必须")
	}
	roleM := model.Role{}
	id_i, _ := strconv.Atoi(id_s)
	roleM.Id = uint(id_i)
	roleM.Name = ctx.PostValue("name")
	roleM.Ename = ""
	roleM.Pid, _ = strconv.Atoi(ctx.PostValue("pid"))
	roleM.Status, _ = strconv.Atoi(ctx.PostValue("status"))
	roleM.Remark = ctx.PostValue("remark")
	roleM.UpdateTime = uint(time.Now().Unix())

	res := service.InitDB().Model(roleM).Save(&roleM)
	if res.Error == nil {
		rc.TopjuiSucess(ctx, "更新成功")
	} else {
		rc.TopjuiError(ctx, "更新失败")
	}
}

func (rc *RoleController) GetAccess(ctx iris.Context) {
	id, _ := strconv.Atoi(ctx.URLParam("id"))
	ctx.ViewData("id", id)
	ctx.View(RequestController + "/access.html")
}

func (rc *RoleController) PostAccesstree(ctx iris.Context) {
	role_id, _ := strconv.Atoi(ctx.URLParam("id"))
	zree := new(service.Node).AccessZree(0, role_id)
	ctx.JSON(zree)
}

func (rc *RoleController) GetDoaccess(ctx iris.Context) {
	node_id_json := ctx.URLParam("node_ids_json")
	role_id, _ := strconv.Atoi(ctx.URLParam("role_id"))
	var node_json []map[string]int
	json.Unmarshal([]byte(node_id_json), &node_json)
	node_list := unique_slice_map(node_json) //[]map[string]int{}去重
	accessM := model.Access{}
	service.InitDB().Debug().Where(map[string]interface{}{"role_id": role_id}).Delete(&accessM)
	var accessMList = []model.Access{}
	var pid_arr []int
	for _, node_vo := range node_list {
		if node_vo["level"] == 1 || node_vo["level"] == 2 {
			continue
		}
		accessMList = append(accessMList, model.Access{RoleId: uint(role_id), NodeId: uint(node_vo["node_id"]), Level: node_vo["level"], Pid: node_vo["pid"]})
		pid_arr = append(pid_arr, node_vo["pid"])
	}
	nodeM := []model.Node{}
	service.InitDB().Debug().Where(map[string]interface{}{"id": pid_arr}).Find(&nodeM)
	if len(nodeM) > 0 {
		for _, node_vo2 := range nodeM {
			accessMList = append(accessMList, model.Access{RoleId: uint(role_id), NodeId: node_vo2.Id, Level: int(node_vo2.Level), Pid: int(node_vo2.Pid)})
		}
		accessMList = append(accessMList, model.Access{RoleId: uint(role_id), NodeId: 1, Level: 1, Pid: 0})
	}
	res := service.InitDB().Create(&accessMList)
	if res.Error == nil {
		rc.TopjuiSucess(ctx, "保存成功")
	} else {
		rc.TopjuiError(ctx, "保存失败")
	}
}

func (rc *RoleController) GetRoleuser(ctx iris.Context) {
	role_id, _ := strconv.Atoi(ctx.URLParam("id"))
	ctx.ViewData("role_id", role_id)
	type roleuser struct {
		model.User
		IsChecked bool
	}
	userM := []*roleuser{}
	service.InitDB().Where("id>?", 1).Find(&userM)
	for _, user_vo := range userM {
		roleuserM := model.RoleUser{}
		res := service.InitDB().Where(map[string]interface{}{"role_id": role_id, "user_id": user_vo.Id}).Find(&roleuserM)
		if res.RowsAffected > 0 {
			user_vo.IsChecked = true
		}
	}
	ctx.ViewData("role_user_list", userM)
	ctx.View(RequestController + "/roleuser.html")
}

func (rc *RoleController) PostRoleuser(ctx iris.Context) {
	role_id, _ := strconv.Atoi(ctx.PostValue("role_id"))
	user_id_arr, _ := ctx.PostValues("user_id[]")
	if role_id != 0 {
		roleuserM := model.RoleUser{}
		service.InitDB().Where(map[string]interface{}{"role_id": role_id}).Delete(&roleuserM)
	}
	if role_id != 0 && len(user_id_arr) > 0 {
		for _, user_id_str := range user_id_arr {
			//user_id ,_:=strconv.Atoi(user_id_str)
			roleuserM2 := model.RoleUser{}
			roleuserM2.RoleId = uint(role_id)
			roleuserM2.UserId = user_id_str
			service.InitDB().Create(&roleuserM2)
		}
		rc.TopjuiSucess(ctx, "设置成功")
	}
	rc.TopjuiSucess(ctx, "设置成功")
}

func (rc *RoleController) GetDataaccess(ctx iris.Context) {
	role_id, _ := strconv.Atoi(ctx.URLParam("role_id"))
	ctx.ViewData("role_id", role_id)
	ctx.View(RequestController + "/data_access.html")
}

func (rc *RoleController) PostData_access_controllers(ctx iris.Context) {
	conf := config.InitConfig()
	controllers := strings.Split(conf.DataAccess.Controllers, ",")
	nodeM := []model.Node{}
	service.InitDB().Where(map[string]interface{}{"level": 2, "status": 1, "name": controllers}).Find(&nodeM)
	ctx.JSON(map[string]interface{}{"rows": nodeM})
}

func (rc *RoleController) PostData_access_actions(ctx iris.Context) {
	pid, _ := strconv.Atoi(ctx.PostValue("id"))
	if pid == 0 {
		rc.TopjuiError(ctx, "参数id丢失")
	}
	nodeM := []model.Node{}
	service.InitDB().Where(map[string]interface{}{"level": 3, "status": 1, "pid": pid}).Find(&nodeM)
	ctx.JSON(map[string]interface{}{"rows": nodeM})
}

func (rc *RoleController) PostData_access_tids(ctx iris.Context) {
	role_id, _ := strconv.Atoi(ctx.URLParam("role_id"))
	id, _ := strconv.Atoi(ctx.PostValue("id"))
	nodeM := model.Node{}
	service.InitDB().Select([]string{"id,pid"}).Where(map[string]interface{}{"id": id}).Find(&nodeM)
	nodeM2 := model.Node{}
	service.InitDB().Select([]string{"id,pid,name"}).Where(map[string]interface{}{"id": nodeM.Pid}).Find(&nodeM2)

	classmodule := nodeM2.Name

	type category2 struct {
		model.Category
		Checked  int    `json:"checked"`
		CheckedS int    `json:"checked_s"`
		NodeId   int    `json:"node_id"`
		Model    string `json:"model"`
	}

	categoryM := []category2{}
	service.InitDB().Select([]string{"classid,classtitle,classmodule"}).Where(map[string]interface{}{"classmodule": classmodule}).Find(&categoryM)

	for key, vo := range categoryM {
		(&categoryM[key]).NodeId = id
		(&categoryM[key]).Model = classmodule
		var count int64
		service.InitDB().Model(model.DataAccess{}).Where(map[string]interface{}{"role_id": role_id, "node_id": id, "tid": vo.Classid}).Count(&count)
		if count > 0 {
			(&categoryM[key]).Checked = 1
			(&categoryM[key]).CheckedS = 1
		}
	}
	ctx.JSON(map[string]interface{}{"rows": categoryM})
}

func (rc *RoleController) PostDoadddataaccess(ctx iris.Context) {
	role_id, _ := strconv.Atoi(ctx.PostValue("role_id"))
	node_id, _ := strconv.Atoi(ctx.PostValue("node_id"))
	tid, _ := strconv.Atoi(ctx.PostValue("tid"))
	model_str := ctx.PostValue("model")

	if role_id == 0 || node_id == 0 || tid == 0 || model_str == "" {
		rc.TopjuiError(ctx, "参数丢失")
		return
	}
	data_accessM := model.DataAccess{}
	res := service.InitDB().Where(map[string]interface{}{"role_id": role_id, "node_id": node_id, "tid": tid}).Find(&data_accessM)
	if res.RowsAffected > 0 {
		service.InitDB().Where(map[string]interface{}{"role_id": role_id, "node_id": node_id, "tid": tid}).Delete(&data_accessM)
	}
	data_accessM.RoleId = uint(role_id)
	data_accessM.NodeId = uint(node_id)
	data_accessM.Tid = uint(tid)
	data_accessM.Model = model_str
	res2 := service.InitDB().Create(&data_accessM)
	if res2.RowsAffected > 0 {
		rc.TopjuiSucess(ctx, "数据授权成功")
	} else {
		rc.TopjuiError(ctx, "数据授权失败："+res2.Error.Error())
	}
}

func (rc *RoleController) PostDodeldataaccess(ctx iris.Context) {
	role_id, _ := strconv.Atoi(ctx.PostValue("role_id"))
	node_id, _ := strconv.Atoi(ctx.PostValue("node_id"))
	tid, _ := strconv.Atoi(ctx.PostValue("tid"))

	if role_id == 0 || node_id == 0 || tid == 0 {
		rc.TopjuiError(ctx, "参数丢失")
		return
	}
	data_accessM := model.DataAccess{}
	res := service.InitDB().Where(map[string]interface{}{"role_id": role_id, "node_id": node_id, "tid": tid}).Delete(&data_accessM)
	if res.Error != nil {
		rc.TopjuiError(ctx, "数据授权删除失败："+res.Error.Error())
	} else {
		rc.TopjuiSucess(ctx, "数据授权删除成功")
	}
}

/**
去重[]map[string]int
*/
func unique_slice_map(slice_map []map[string]int) []map[string]int {
	s_len := len(slice_map) - 1
	for ; s_len > 0; s_len-- {
		for j := s_len - 1; j >= 0; j-- {
			if cmpMap(slice_map[s_len], slice_map[j]) {
				slice_map = append(slice_map[:s_len], slice_map[s_len+1:]...)
				break
			}
		}
	}
	return slice_map
}

/**
比较两个map[string]int是否相同
*/
func cmpMap(m1, m2 map[string]int) bool {
	for k1, v1 := range m1 {
		if v2, ok := m2[k1]; ok {
			if v1 != v2 {
				return false
			}
		} else {
			return false
		}
	}
	for k2, v2 := range m2 {
		if v1, ok := m1[k2]; ok {
			if v1 != v2 {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

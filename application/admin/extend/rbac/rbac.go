package rbac

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/v12"
	"gocms/application/app_session"
	"gocms/application/model"
	"gocms/config"
	"strconv"
	"strings"
)

/*
 *数据库表:
 * -- --------------------------------------------------------
 * CREATE TABLE IF NOT EXISTS `access` (
 * `role_id` smallint(6) unsigned NOT NULL,
 * `node_id` smallint(6) unsigned NOT NULL,
 * `level` tinyint(1) NOT NULL,
 * `module` varchar(50) DEFAULT NULL,
 * KEY `groupId` (`role_id`),
 * KEY `nodeId` (`node_id`)
 * ) ENGINE=MyISAM DEFAULT CHARSET=utf8;
 *
 * CREATE TABLE IF NOT EXISTS `node` (
 * `id` smallint(6) unsigned NOT NULL AUTO_INCREMENT,
 * `name` varchar(20) NOT NULL,
 * `title` varchar(50) DEFAULT NULL,
 * `status` tinyint(1) DEFAULT '0',
 * `remark` varchar(255) DEFAULT NULL,
 * `sort` smallint(6) unsigned DEFAULT NULL,
 * `pid` smallint(6) unsigned NOT NULL,
 * `level` tinyint(1) unsigned NOT NULL,
 * PRIMARY KEY (`id`),
 * KEY `level` (`level`),
 * KEY `pid` (`pid`),
 * KEY `status` (`status`),
 * KEY `name` (`name`)
 * ) ENGINE=MyISAM DEFAULT CHARSET=utf8;
 *
 * CREATE TABLE IF NOT EXISTS `role` (
 * `id` smallint(6) unsigned NOT NULL AUTO_INCREMENT,
 * `name` varchar(20) NOT NULL,
 * `pid` smallint(6) DEFAULT NULL,
 * `status` tinyint(1) unsigned DEFAULT NULL,
 * `remark` varchar(255) DEFAULT NULL,
 * PRIMARY KEY (`id`),
 * KEY `pid` (`pid`),
 * KEY `status` (`status`)
 * ) ENGINE=MyISAM DEFAULT CHARSET=utf8 ;
 *
 * CREATE TABLE IF NOT EXISTS `role_user` (
 * `role_id` mediumint(9) unsigned DEFAULT NULL,
 * `user_id` char(32) DEFAULT NULL,
 * KEY `group_id` (`role_id`),
 * KEY `user_id` (`user_id`)
 * ) ENGINE=MyISAM DEFAULT CHARSET=utf8;
 *
 * CREATE TABLE IF NOT EXISTS `user` (
 * `id` mediumint(9) unsigned DEFAULT NULL,
 * `name` char(32) DEFAULT NULL,
 *.......
 * ) ENGINE=InnoDB DEFAULT CHARSET=utf8;
 */
type RBAC struct {
	NodeTable             string
	RoleTable             string
	RoleUserTable         string
	UserTable             string
	AccessTable           string
	UserAuthType          int
	UserAuthKey           string
	AdminAuthKey          string
	RequireAuthController string
	NoAuthController      string
	RequireAuthAction     string
	NoAuthAction          string
}

func initRBAC() *RBAC {
	app_conf := config.InitConfig()
	return &RBAC{
		NodeTable:             app_conf.Rbac.TableNode,
		RoleTable:             app_conf.Rbac.TableRole,
		RoleUserTable:         app_conf.Rbac.TableUser,
		UserTable:             app_conf.Rbac.TableUser,
		AccessTable:           app_conf.Rbac.TableAccess,
		AdminAuthKey:          app_conf.Rbac.AdminAuthKey,
		UserAuthType:          app_conf.Rbac.UserAuthType,
		UserAuthKey:           app_conf.Rbac.UserAuthKey,
		RequireAuthController: app_conf.Rbac.RequireAuthController,
		NoAuthController:      app_conf.Rbac.NoAuthController,
		RequireAuthAction:     app_conf.Rbac.RequireAuthAction,
		NoAuthAction:          app_conf.Rbac.NoAuthAction}
}

func (r *RBAC) Init(NodeTable, RoleTable, RoleUserTable, UserTable, AccessTable string) *RBAC {
	r.NodeTable = NodeTable
	r.RoleTable = RoleTable
	r.RoleUserTable = RoleUserTable
	r.UserTable = UserTable
	r.AccessTable = AccessTable
	return r
}

/**
*认证账号密码
 */
func (r *RBAC) Authenticate(whereMap map[string]interface{}) (bool, model.User) {
	var user model.User
	whereMap["status"] = 1
	model.InitDB().Where(whereMap).Last(&user)
	if user.Id > 0 {
		return true, user
	}
	return false, user
}

/**
用于检测用户权限的方法,并保存到Session中
*/
func (r *RBAC) SaveAccessList(ctx iris.Context, authId uint) {
	r = initRBAC()
	session := app_session.Sess.Start(ctx)
	if is_admin, _ := session.GetBoolean(r.AdminAuthKey); r.UserAuthType != 2 && !is_admin {
		session.Set("_ACCESS_LIST", r.GetAccessMap(authId))
	}
}

// 检查当前操作是否需要认证
func (r *RBAC) CheckAccess(request_controller, request_action string) bool {
	r = initRBAC()
	controller, action := map[string][]string{}, map[string][]string{} //初始化map,nil map不能赋值,不能var nap[string][]srring
	if "all" != r.RequireAuthController {
		controller["yes"] = strings.Split(strings.ToUpper(r.RequireAuthController), ",") //需要认证的控制器切片
	} else {
		controller["no"] = strings.Split(strings.ToUpper(r.NoAuthController), ",") //不需要认证的控制器切片
	}
	// 检查当前模块是否需要认证
	if (len(controller["no"]) != 0 && !isInSlice(strings.ToUpper(request_controller), controller["no"])) || (len(controller["yes"]) != 0 && isInSlice(strings.ToUpper(request_controller), controller["yes"])) {
		if "all" != r.RequireAuthAction {
			action["yes"] = strings.Split(strings.ToUpper(r.RequireAuthAction), ",") //需要认证的action切片
		} else {
			action["no"] = strings.Split(strings.ToUpper(r.RequireAuthAction), ",") //不需要认证的action切片
		}
		// 检查当前操作是否需要认证
		if (len(action["no"]) != 0 && !isInSlice(strings.ToUpper(request_action), action["no"])) || (len(action["yes"]) != 0 && isInSlice(strings.ToUpper(request_action), action["yes"])) {
			return true
		}
	}
	return false
}

/**
判断切片中是否包含字符串
*/
func isInSlice(s string, slice []string) bool {
	for _, si := range slice {
		if si == s {
			return true
		}
	}
	return false
}

// 权限认证的过滤器方法
func (r *RBAC) AccessDecision(request_app, request_controller, request_action string, ctx iris.Context) bool {
	r = initRBAC()
	// 检查是否需要认证
	if r.CheckAccess(request_controller, request_action) {
		accessGuid := fmt.Sprintf("%x", md5.Sum([]byte(request_app+request_controller+request_action))) // 存在认证识别号，则进行进一步的访问决策
		session := app_session.Sess.Start(ctx)
		if isok, _ := session.GetBoolean(r.AdminAuthKey); !isok {
			var accessMap map[string]map[string]map[string]string
			if r.UserAuthType == 1 { //登录验证
				accessMap = make(map[string]map[string]map[string]string)
				accessInterface := session.Get("_ACCESS_LIST")
				accessInterface_b, _ := json.Marshal(accessInterface)
				json.Unmarshal(accessInterface_b, &accessMap)
			} else {
				//r.UserAuthType == 2
				// 加强验证和即时验证模式 更加安全 后台权限修改可以即时生效
				// 通过数据库进行访问检查
				authId, _ := session.GetInt(r.UserAuthKey)
				accessMap = r.GetAccessMap(uint(authId))
			}
			if accessMap[strings.ToUpper(request_app)][strings.ToUpper(request_controller)][strings.ToUpper(request_action)] != "" {
				session.Set(accessGuid, true)
			} else {
				session.Set(accessGuid, false)
				return false
			}
		} else { //超级管理员无需认证
			return true
		}
	}
	return true
}

/**
 * +----------------------------------------------------------
 * 取得当前认证号的所有权限列表
 * +----------------------------------------------------------
 */
func (r *RBAC) GetAccessMap(authId uint) map[string]map[string]map[string]string {
	app_conf := config.InitConfig()
	r2 := r.Init(app_conf.Rbac.TableNode, app_conf.Rbac.TableRole, app_conf.Rbac.TableRoleUser, app_conf.Rbac.TableUser, app_conf.Rbac.TableAccess)
	sql := "select node.id,node.name from " + r2.RoleTable + " as role," + r2.RoleUserTable + " as user," + r2.AccessTable + " as access ," + r2.NodeTable + " as node " + "where user.user_id='" + strconv.Itoa(int(authId)) + "' and user.role_id=role.id and ( access.role_id=role.id  or (access.role_id=role.pid and role.pid!=0 ) ) and role.status=1 and access.node_id=node.id and node.level=1 and node.status=1"
	apps := []*nodeAccess{}
	model.InitDB().Raw(sql).Scan(&apps)
	access := make(map[string]map[string]map[string]string)
	for _, app := range apps {
		appId := app.Id
		appName := app.Name
		sql = "select node.id,node.name from " + r2.RoleTable + " as role," + r2.RoleUserTable + " as user," + r2.AccessTable + " as access ," + r2.NodeTable + " as node " + "where user.user_id='" + strconv.Itoa(int(authId)) + "' and user.role_id=role.id and ( access.role_id=role.id  or (access.role_id=role.pid and role.pid!=0 ) ) and role.status=1 and access.node_id=node.id and node.level=2 and node.pid=" + strconv.Itoa(int(appId)) + " and node.status=1"
		modules := []*nodeAccess{}
		model.InitDB().Raw(sql).Scan(&modules)
		moduleMap := make(map[string]map[string]string)
		for _, module := range modules {
			moduleId := module.Id
			moduleName := module.Name
			sql = "select node.id,node.name from " + r2.RoleTable + " as role," + r2.RoleUserTable + " as user," + r2.AccessTable + " as access ," + r2.NodeTable + " as node " + "where user.user_id='" + strconv.Itoa(int(authId)) + "' and user.role_id=role.id and ( access.role_id=role.id  or (access.role_id=role.pid and role.pid!=0 ) ) and role.status=1 and access.node_id=node.id and node.level=3 and node.pid=" + strconv.Itoa(int(moduleId)) + " and node.status=1"
			actions := []*nodeAccess{}
			model.InitDB().Raw(sql).Scan(&actions)
			actionMap := make(map[string]string)
			for _, action := range actions {
				actionId := action.Id
				actionName := action.Name
				actionMap[strings.ToUpper(actionName)] = strconv.Itoa(int(actionId))
			}
			moduleMap[strings.ToUpper(moduleName)] = actionMap
		}
		access[strings.ToUpper(appName)] = moduleMap
	}
	return access
}

//数据库sql查询后的类型结构体
type nodeAccess struct {
	Id   uint
	Name string
}

/**
* 读取模块所属的记录访问权限
 */
func (r *RBAC) GetModuleAccessList(authId uint, module string) {
	app_conf := config.InitConfig()
	r2 := r.Init(app_conf.Rbac.TableNode, app_conf.Rbac.TableRole, app_conf.Rbac.TableRoleUser, app_conf.Rbac.TableUser, app_conf.Rbac.TableAccess)
	sql := "select access.node_id from " + r2.RoleTable + " as role," + r2.RoleUserTable + " as user," + r2.AccessTable + " as access where user.user_id='" + strconv.Itoa(int(authId)) + "' and user.role_id=role.id and ( access.role_id=role.id  or (access.role_id=role.pid and role.pid!=0 ) ) and role.status=1 and  access.module='" + module + "'"
	aa := []interface{}{}
	model.InitDB().Raw(sql).Scan(&aa)
}

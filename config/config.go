package config

import (
	_ "embed"
	"encoding/json"
)

//go:embed config.json
var Config_json string

//go:embed config_demo.json
var Config_demo_json string

type AppConfig struct {
	AppProtocol     string         `json:"App_protocol"`     //协议 http/https
	AppUrl          string         `json:"app_url"`          //域名
	AppUrlPort      int            `json:"app_url_port"`     //端口号
	AppName         string         `json:"app_name"`         //网站名称
	SiteKeywords    string         `json:"site_keywords"`    //网站关键字
	SiteDescription string         `json:"site_description"` //网站描述
	AppAuthorEmail  string         `json:"app_author_email"` //作者邮箱
	OfflineMessage  string         `json:"offline_message"`  //网站关闭时提示问题
	SiteIcpNum      string         `json:"site_icp_num"`     //icp备案号
	AppVersion      string         `json:"app_version"`      //版本
	Db              DbConfig       `json:"db"`
	Rbac            RbacConfig     `json:"rbac"`
	DataAccess      DataAccessConf `json:"data_access"`
	Admin           Admin          `json:"admin"`
	Upload          UploadConf     `json:"upload"`
}

type DbConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type UploadConf struct {
	FileUploadType string         `json:"file_upload_type"` // 可能值:Local、AliyunOSS、Bcs、Ftp、Qiniu、Sae、Upyun
	UploadTypeConf UploadTypeConf `json:"upload_type_config"`
}

type UploadTypeConf struct {
	MaxSize         int64  `json:"max_size"`          // 文件大小限制 限制在100M
	RootPath        string `json:"root_path"`         // 保存根路径
	Regionid        string `json:"region_id"`         // Regionid
	AccessKeyId     string `json:"access_key_id"`     // AccessKeyId
	AccessKeySecret string `json:"access_key_secret"` // AccessKeySecret
	Bucket          string `json:"bucket"`            // Bucket
	Endpoint        string `json:"endpoint"`          // Endpoint设置,如果是阿里云,杭州节点,杭州内网节点，如果ECS服务器也是杭州节点，请使用该内网节点"'oss-cn-hangzhou-internal.aliyuncs.com"，节省流量和较快的上传速度
	RequestUri      string `json:"request_uri"`       // 文件内网访问的url，内网访问不走流量，速度更快。
	Requesturl      string `json:"request_url"`       // 文件访问的url
}

type Admin struct {
	Menu Menu `json:"menu"` //后台大菜单
}

type Menu struct {
	Content    Content `json:"content"`
	Framework  Content `json:"framework"`
	Groupuser  Content `json:"groupuser"`
	Developers Content `json:"developers"`
	System     Content `json:"system"`
	Plugin     Content `json:"plugin"`
}

type Content struct {
	Icons string `json:"icons"`
	Title string `json:"title"`
}

type RbacConfig struct {
	TableNode             string `json:"table_node"`              //节点表的表名
	TableAccess           string `json:"table_access"`            //权限表的表名
	TableRole             string `json:"table_role"`              //角色表的表名
	TableRoleUser         string `json:"table_role_user"`         //角色和用户的关系表
	TableUser             string `json:"table_user"`              //用户表的表名
	UserAuthType          int    `json:"user_auth_type"`          //验证类型:1 登录认证 2 实时认证
	AdminAuthKey          string `json:"admin_auth_key"`          //超级管理员
	UserAuthKey           string `json:"user_auth_key"`           //用户认证SESSION标记
	RequireAuthController string `json:"require_auth_controller"` //需要认证的控制器    多个用","连接,all代表全部需要
	NoAuthController      string `json:"no_auth_controller"`      //不需要认证的控制    多个用","连接
	RequireAuthAction     string `json:"require_auth_action"`     //需要认证的action   多个用","连接,all代表全部需要
	NoAuthAction          string `json:"no_auth_action"`          //不需要认证的action   多个用","连接
}

type DataAccessConf struct {
	IsOpen      int    `json:"is_open"`
	Controllers string `json:"controllers"`
	TableField  string `json:"table_field"`
}

//初始化服务器配置
/**
func InitConfig() *AppConfig {
	file, err := os.Open("./config/config.json")
	if err != nil {
		panic(err.Error())
	}
	decoder := json.NewDecoder(file)
	conf := AppConfig{}
	err = decoder.Decode(&conf)
	if err != nil {
		panic(err.Error())
	}
	return &conf
}
*/
func InitConfig() *AppConfig {
	conf := AppConfig{}
	err := json.Unmarshal([]byte(Config_json), &conf)
	if err != nil {
		panic(err.Error())
	}
	return &conf
}

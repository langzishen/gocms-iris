package model

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

//分类表
type Category struct {
	Classid          int    `gorm:"column:classid;type:int(10);primary_key;AUTO_INCREMENT" json:"classid"`        // 栏目id
	Classpid         int    `gorm:"column:classpid;type:int(10);default:0;NOT NULL" json:"classpid"`              // 栏目父id
	Classpids        string `gorm:"column:classpids;type:varchar(100)" json:"classpids"`                          // 栏目父ids
	Classchild       int    `gorm:"column:classchild;type:tinyint(1);default:0;NOT NULL" json:"classchild"`       // 是否有下级
	Classchildids    string `gorm:"column:classchildids;type:varchar(2000)" json:"classchildids"`                 // 栏目下级ids
	Classarrchildids string `gorm:"column:classarrchildids;type:mediumtext" json:"classarrchildids"`              // 栏目下级对象
	Classtitle       string `gorm:"column:classtitle;type:varchar(100);NOT NULL" json:"classtitle"`               // 栏目标题
	Classsort        int    `gorm:"column:classsort;type:smallint(5);default:0;NOT NULL" json:"classsort"`        // 排序
	Classstatus      int    `gorm:"column:classstatus;type:tinyint(1);default:1;NOT NULL" json:"classstatus"`     // 状态
	Classkeywords    string `gorm:"column:classkeywords;type:varchar(150)" json:"classkeywords"`                  // 关键字
	Classdescription string `gorm:"column:classdescription;type:varchar(200)" json:"classdescription"`            // 描述
	Classmodule      string `gorm:"column:classmodule;type:varchar(30)" json:"classmodule"`                       // 所属模型
	Classrewrite     string `gorm:"column:classrewrite;type:varchar(100)" json:"classrewrite"`                    // URL重写值
	Classtemplate    string `gorm:"column:classtemplate;type:varchar(50)" json:"classtemplate"`                   // 栏目模版
	Newstemplate     string `gorm:"column:newstemplate;type:varchar(50)" json:"newstemplate"`                     // 文章模版
	Classimg         string `gorm:"column:classimg;type:varchar(100)" json:"classimg"`                            // 栏目预览图
	Classapv         int    `gorm:"column:classapv;type:int(10);default:0;NOT NULL" json:"classapv"`              // 栏目浏览量
	Classdomain      string `gorm:"column:classdomain;type:varchar(100)" json:"classdomain"`                      // 栏目二级域名
	Classouturl      string `gorm:"column:classouturl;type:varchar(150)" json:"classouturl"`                      // 栏目外部网址
	Classmenushow    int    `gorm:"column:classmenushow;type:tinyint(1);default:1;NOT NULL" json:"classmenushow"` // 前台菜单中显示
}

func (m *Category) TableName() string {
	return "category"
}

/**
获取结构列表
classstatus:可选值：'*','p1,p2,p3...',多个用,链接
arg:不定参
*/
func (m *Category) ZreeList(classstatus string, classpid int, arg ...map[string]interface{}) (zree_category []map[string]interface{}) {
	type zree struct {
		*Category
		State string `json:"state"`
	}
	where_map := make(map[string]interface{})
	if len(arg) > 0 {
		for _, arg_vo := range arg {
			for arg_vo_key, arg_vo_value := range arg_vo {
				where_map[arg_vo_key] = arg_vo_value
			}
		}
	}
	if classstatus != "*" {
		classstatus_arr := strings.Split(classstatus, ",")
		var classstatus_s []int
		for _, classstatus_ss := range classstatus_arr {
			classstatus_i, _ := strconv.Atoi(classstatus_ss)
			classstatus_s = append(classstatus_s, classstatus_i)
		}
	}
	where_map["classpid"] = classpid
	var zree_list []*zree
	InitDB().Where(where_map).Find(&zree_list)
	for _, zree_vo := range zree_list {
		/*(
		if zree_vo.Classchild == 1 {
			zree_vo.State = "closed"
		} else {
			zree_vo.State = "open"
		}*/
		var is_child int64
		InitDB().Model(Category{}).Where(map[string]interface{}{"classpid": zree_vo.Classid}).Count(&is_child)
		if is_child > 0 {
			zree_vo.State = "closed"
		} else {
			zree_vo.State = "open"
		}
	}
	zree_json_byte, _ := json.Marshal(zree_list)
	json.Unmarshal(zree_json_byte, &zree_category)
	return zree_category
}

/**
根据classid返回classpids
*/
func (m *Category) get_classpids(classid int) (classpids string) {
	categoryM := Category{}
	InitDB().Where(map[string]interface{}{"classid": classid}).Find(&categoryM)
	if categoryM.Classpids == "" {
		classpids = fmt.Sprintf("0,%d", classid)
	} else {
		pids_arr := strings.Split(categoryM.Classpids, ",")
		pids_arr = pids_arr[0 : len(pids_arr)-1]
		classpids = fmt.Sprintf("%s,%d", strings.Join(pids_arr, ","), classid)
	}
	return classpids
}

/**
根据classid返回父类的classpids
*/
func (m *Category) get_parent_classpids(classid int) (classpids string) {
	categoryM := Category{}
	InitDB().Where(map[string]interface{}{"classid": classid}).Find(&categoryM)
	classpids = m.get_classpids(categoryM.Classpid)
	return classpids
}

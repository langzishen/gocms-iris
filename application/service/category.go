package service

import (
	"encoding/json"
	"fmt"
	"gocms/application/model"
	"strconv"
	"strings"
)

type Category struct {
	Base
}

/**
获取结构列表
classstatus:可选值：'*','p1,p2,p3...',多个用,链接
arg:不定参
*/
func (m *Category) ZreeList(classstatus string, classpid int, arg ...map[string]interface{}) (zree_category []map[string]interface{}) {
	type zree struct {
		*model.Category
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
	categoryM := model.Category{}
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
	categoryM := model.Category{}
	InitDB().Where(map[string]interface{}{"classid": classid}).Find(&categoryM)
	classpids = m.get_classpids(categoryM.Classpid)
	return classpids
}

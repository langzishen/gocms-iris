package service

import "gocms/application/model"

type Model struct {
	Base
}

/**
获取模型列表
	默认：状态：1
*/
func (m *Model) GetModelList(arg ...map[string]interface{}) []*model.Model {
	list := []*model.Model{}
	where_map := make(map[string]interface{})
	if len(arg) > 0 {
		for _, arg_vo := range arg {
			for arg_vo_key, arg_vo_val := range arg_vo {
				where_map[arg_vo_key] = arg_vo_val
			}
		}
	} else {
		where_map["status"] = 1
	}

	InitDB().Where(where_map).Find(&list)
	return list
}

package service

import "gocms/application/model"

type Attribute struct {
	Base
}

func (m *Attribute) GetList(arg ...map[string]interface{}) (list []*model.Attribute) {
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

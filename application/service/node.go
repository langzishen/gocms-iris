package service

import "gocms/application/model"

type Node struct {
	Base
}

func (m *Node) Zree() {

}

func (m *Node) AccessZree(pid, role_id int) []map[string]interface{} {
	return m.acceesszree(pid, role_id)
}

func (m *Node) acceesszree(pid, role_id int) []map[string]interface{} {
	nodeM := []model.Node{}
	InitDB().Select([]string{"id", "pid", "name", "icon", "title"}).Where(map[string]interface{}{"status": 1, "pid": pid}).Order("sort asc,id desc").Find(&nodeM)
	var zree []map[string]interface{}
	if len(nodeM) > 0 {
		for _, node_vo := range nodeM {
			zree_zree := m.acceesszree(int(node_vo.Id), role_id)
			if len(zree_zree) > 0 {
				zree = append(zree, map[string]interface{}{"iconCls": node_vo.Icon, "name": node_vo.Name, "text": node_vo.Title + "[" + node_vo.Name + "]", "status": node_vo.Status, "id": int(node_vo.Id), "pid": int(node_vo.Pid), "state": "open", "checkbox_checked": true, "children": zree_zree})
			} else {
				accessM := Access{}
				res := InitDB().Where(map[string]interface{}{"role_id": role_id, "node_id": node_vo.Id}).Find(&accessM)
				if res.RowsAffected > 0 {
					zree = append(zree, map[string]interface{}{"iconCls": node_vo.Icon, "name": node_vo.Name, "text": node_vo.Title + "[" + node_vo.Name + "]", "status": node_vo.Status, "id": int(node_vo.Id), "pid": int(node_vo.Pid), "checked": true, "checkbox_checked": true})
				} else {
					zree = append(zree, map[string]interface{}{"iconCls": node_vo.Icon, "name": node_vo.Name, "text": node_vo.Title + "[" + node_vo.Name + "]", "status": node_vo.Status, "id": int(node_vo.Id), "pid": int(node_vo.Pid), "checkbox_checked": true})
				}
			}
		}
		return zree
	} else {
		return []map[string]interface{}{}
	}
}

package model

//节点表---------RBAC
type Node struct {
	Id             uint   `gorm:"column:id;type:int(10) unsigned;primary_key;AUTO_INCREMENT" json:"id"`
	Name           string `gorm:"column:name;type:varchar(20);NOT NULL" json:"name"`
	Title          string `gorm:"column:title;type:varchar(50);NOT NULL" json:"title"`
	Status         int    `gorm:"column:status;type:tinyint(1);default:0;NOT NULL" json:"status"`
	Remark         string `gorm:"column:remark;type:varchar(255);NOT NULL" json:"remark"`
	Sort           uint   `gorm:"column:sort;type:smallint(5) unsigned;default:0;NOT NULL" json:"sort"`
	Pid            uint   `gorm:"column:pid;type:smallint(5) unsigned;default:0;NOT NULL" json:"pid"`
	Level          uint   `gorm:"column:level;type:tinyint(1) unsigned;default:0;NOT NULL" json:"level"`
	Type           int    `gorm:"column:type;type:tinyint(1);default:0;NOT NULL" json:"type"`
	GroupId        string `gorm:"column:group_id;type:varchar(15);NOT NULL" json:"group_id"`
	Icon           string `gorm:"column:icon;type:varchar(1000)" json:"icon"` // 图标
	LeftMenuAction uint   `gorm:"column:left_menu_action;type:tinyint(1);default:0;NOT NULL" json:"left_menu_action"`
}

func (m *Node) TableName() string {
	return "node"
}

func (m *Node) Zree() {

}

func (m *Node) AccessZree(pid, role_id int) []map[string]interface{} {
	return m.acceesszree(pid, role_id)
}

func (m *Node) acceesszree(pid, role_id int) []map[string]interface{} {
	nodeM := []Node{}
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

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

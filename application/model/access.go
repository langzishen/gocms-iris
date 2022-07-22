package model

// 权限表----------RBAC
type Access struct {
	RoleId uint   `gorm:"column:role_id;type:int(10) unsigned;NOT NULL" json:"role_id"`
	NodeId uint   `gorm:"column:node_id;type:int(10) unsigned;NOT NULL" json:"node_id"`
	Level  int    `gorm:"column:level;type:tinyint(1);NOT NULL" json:"level"`
	Pid    int    `gorm:"column:pid;type:int(10);NOT NULL" json:"pid"`
	Module string `gorm:"column:module;type:varchar(50);NOT NULL" json:"module"`
}

func (m *Access) TableName() string {
	return "access"
}

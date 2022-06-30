package model

// 数据授权表-------RBAC增强版
type DataAccess struct {
	RoleId     uint   `gorm:"column:role_id;type:int(10) unsigned;NOT NULL" json:"role_id"`                        // 用户组ID
	NodeId     uint   `gorm:"column:node_id;type:int(10) unsigned;NOT NULL" json:"node_id"`                        // 节点ID
	Tid        uint   `gorm:"column:tid;type:int(10) unsigned;NOT NULL" json:"tid"`                                // 分类ID
	Model      string `gorm:"column:model;type:varchar(100);NOT NULL" json:"model"`                                // 模型名称
	CreateTime uint   `gorm:"column:create_time;type:int(10) unsigned;NOT NULL;autoCreateTime" json:"create_time"` // 添加时间
}

func (m *DataAccess) TableName() string {
	return "data_access"
}

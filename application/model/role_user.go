package model

//用户组与用户关系表 ---------RBAC
type RoleUser struct {
	RoleId uint   `gorm:"column:role_id;type:int(10) unsigned;default:0;NOT NULL" json:"role_id"`
	UserId string `gorm:"column:user_id;type:char(32);default:0;NOT NULL" json:"user_id"`
}

func (m *RoleUser) TableName() string {
	return "role_user"
}

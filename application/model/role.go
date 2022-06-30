package model

//用户组表----------RBAC
type Role struct {
	Id         uint   `gorm:"column:id;type:int(10) unsigned;primary_key;AUTO_INCREMENT" json:"id"`
	Name       string `gorm:"column:name;type:varchar(30);NOT NULL" json:"name"`
	Pid        int    `gorm:"column:pid;type:smallint(5);default:0;NOT NULL" json:"pid"`
	Status     int    `gorm:"column:status;type:tinyint(1);default:0;NOT NULL" json:"status"`
	Remark     string `gorm:"column:remark;type:varchar(255);NOT NULL" json:"remark"`
	Ename      string `gorm:"column:ename;type:varchar(5);NOT NULL" json:"ename"`
	CreateTime uint   `gorm:"column:create_time;type:int(10) unsigned;NOT NULL;autoCreateTime" json:"create_time"`
	UpdateTime uint   `gorm:"column:update_time;type:int(10) unsigned;NOT NULL;autoUpdateTime" json:"update_time"`
}

func (m *Role) TableName() string {
	return "role"
}

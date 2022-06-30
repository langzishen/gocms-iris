package model

//用户表---后台管理员----RBAC
type User struct {
	Id         uint   `gorm:"column:id;type:int(10) unsigned;primary_key" json:"id"`
	ObjectId   string `gorm:"column:object_id;type:varchar(32);NOT NULL" json:"object_id"` // 用户唯一标识
	Account    string `gorm:"column:account;type:varchar(40);NOT NULL" json:"account"`     // 用户名
	Nickname   string `gorm:"column:nickname;type:varchar(200);NOT NULL" json:"nickname"`
	Logo       string `gorm:"column:logo;type:varchar(1000);NOT NULL" json:"logo"`
	Age        uint   `gorm:"column:age;type:tinyint(3);NOT NULL" json:"age"`
	Sex        uint   `gorm:"column:sex;type:tinyint(1);NOT NULL" json:"sex"`
	Password   string `gorm:"column:password;type:char(32);NOT NULL" json:"password"`
	Phone      string `gorm:"column:phone;type:varchar(13)" json:"phone"`
	Email      string `gorm:"column:email;type:varchar(150)" json:"email"`
	CreateTime uint   `gorm:"column:create_time;type:int(10) unsigned;NOT NULL;autoCreateTime" json:"create_time"`
	UpdateTime uint   `gorm:"column:update_time;type:int(10) unsigned;default:0;NOT NULL;autoUpdateTime" json:"update_time"`
	Status     int    `gorm:"column:status;type:tinyint(1);default:1;NOT NULL" json:"status"`            // 状态：默认1，1为正常0为禁用可恢复，-1为删除
	TypeId     uint   `gorm:"column:type_id;type:tinyint(2) unsigned;default:0;NOT NULL" json:"type_id"` // 9为后台超级管理员
}

func (m *User) TableName() string {
	return "user"
}

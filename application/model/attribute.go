package model

//属性表
type Attribute struct {
	Id         int    `gorm:"column:id;type:int(10);primary_key;AUTO_INCREMENT" json:"id"`     // id
	ModelEname string `gorm:"column:model_ename;type:varchar(80);NOT NULL" json:"model_ename"` // 对应模型名称
	Attrname   string `gorm:"column:attrname;type:varchar(50);NOT NULL" json:"attrname"`       // 属性名称
	Attrvalue  string `gorm:"column:attrvalue;type:varchar(50);NOT NULL" json:"attrvalue"`     // 属性值
	Sort       int    `gorm:"column:sort;type:mediumint(9);default:999;NOT NULL" json:"sort"`  // 排序值
	Status     int    `gorm:"column:status;type:tinyint(1);default:1;NOT NULL" json:"status"`  // 状态 1:启用,0:禁用
}

func (m *Attribute) TableName() string {
	return "attribute"
}

package model

type Photo struct {
	Id         uint   `gorm:"column:id;type:int(11) unsigned;primary_key;AUTO_INCREMENT;comment:逐渐自增" json:"id"`
	Tid        string `gorm:"column:tid;type:varchar(255);comment:分类ID;NOT NULL" json:"tid"`
	Title      string `gorm:"column:title;type:varchar(255);comment:标题;NOT NULL" json:"title"`
	Img        string `gorm:"column:img;type:varchar(2000);comment:图片地址;NOT NULL" json:"img"`
	CreateTime uint   `gorm:"column:create_time;type:int(11) unsigned;comment:创建时间;NOT NULL;autoCreateTime" json:"create_time"`
	UpdateTime uint   `gorm:"column:update_time;type:int(11) unsigned;comment:更新时间;NOT NULL;autoUpdateTime" json:"update_time"`
	CreaterId  uint   `gorm:"column:creater_id;type:int(11) unsigned;comment:创建者;NOT NULL" json:"creater_id"`
	Status     int    `gorm:"column:status;type:tinyint(1);default:1;comment:状态;NOT NULL" json:"status"`
}

func (m *Photo) TableName() string {
	return "photo"
}

package model

// 模型表
type Model struct {
	Id         int    `gorm:"column:id;type:int(10);primary_key;AUTO_INCREMENT" json:"id"`          // 模型表
	Ename      string `gorm:"column:ename;type:varchar(20);NOT NULL" json:"ename"`                  // 模块名称
	Cname      string `gorm:"column:cname;type:varchar(50);NOT NULL" json:"cname"`                  // 显示名称
	Notes      string `gorm:"column:notes;type:text" json:"notes"`                                  // 应用描述
	Menugroup  int    `gorm:"column:menugroup;type:tinyint(2);default:2;NOT NULL" json:"menugroup"` // 属于大菜单
	Sort       int    `gorm:"column:sort;type:tinyint(3);default:1;NOT NULL" json:"sort"`           // 排序
	Author     string `gorm:"column:author;type:varchar(30);NOT NULL" json:"author"`                // 开发作者
	Version    string `gorm:"column:version;type:varchar(15);NOT NULL" json:"version"`              // 版本
	CreateTime int    `gorm:"column:create_time;type:int(10);default:0;NOT NULL;autoCreateTime" json:"create_time"`
	UpdateTime int    `gorm:"column:update_time;type:int(10);default:0;NOT NULL;autoUpdateTime" json:"update_time"`
	Status     int    `gorm:"column:status;type:tinyint(1);default:1;NOT NULL" json:"status"` // 状态
}

func (m *Model) TableName() string {
	return "model"
}

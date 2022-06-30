package model

//文章表
type Article struct {
	Id          int    `gorm:"column:id;type:int(10);primary_key;AUTO_INCREMENT" json:"id"`                           // id
	Tid         string `gorm:"column:tid;type:varchar(50);NOT NULL" json:"tid"`                                       // 所属分类
	Title       string `gorm:"column:title;type:varchar(100);NOT NULL" json:"title"`                                  // 文章标题
	Keywords    string `gorm:"column:keywords;type:varchar(100)" json:"keywords"`                                     // 关键字
	Description string `gorm:"column:description;type:varchar(200)" json:"description"`                               // 描述
	Img         string `gorm:"column:img;type:varchar(100)" json:"img"`                                               // 预览图片
	Content     string `gorm:"column:content;type:longtext;NOT NULL" json:"content"`                                  // 内容
	CreateTime  int    `gorm:"column:create_time;type:int(10);NOT NULL;autoCreateTime" json:"create_time"`            // 录入时间
	UpdateTime  int    `gorm:"column:update_time;type:int(10);NOT NULL;autoUpdateTime" json:"update_time"`            // 修改时间
	CreaterId   int    `gorm:"column:creater_id;type:int(10);NOT NULL" json:"creater_id"`                             // 录入人
	Sort        int    `gorm:"column:sort;type:smallint(5);default:0;NOT NULL" json:"sort"`                           // 排序值
	Apv         int    `gorm:"column:apv;type:smallint(5);default:0;NOT NULL" json:"apv"`                             // 浏览量
	Rewrite     string `gorm:"column:rewrite;type:varchar(50)" json:"rewrite"`                                        // URL重写值
	Status      int    `gorm:"column:status;type:tinyint(1);default:1;NOT NULL" json:"status"`                        // 状态 1:启用,0:禁用
	Template    string `gorm:"column:template;type:varchar(50)" json:"template"`                                      // 使用模板
	Attrtj      string `gorm:"column:attrtj;type:varchar(30)" json:"attrtj"`                                          // 推荐属性
	Outurl      string `gorm:"column:outurl;type:varchar(150)" json:"outurl"`                                         // 外部网址
	IsLoginShow uint   `gorm:"column:is_login_show;type:tinyint(1) unsigned;default:0;NOT NULL" json:"is_login_show"` // 是否需要登录才显示
}

func (m *Article) TableName() string {
	return "article"
}

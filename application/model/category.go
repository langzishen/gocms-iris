package model

//分类表
type Category struct {
	Classid          int    `gorm:"column:classid;type:int(10);primary_key;AUTO_INCREMENT" json:"classid"`        // 栏目id
	Classpid         int    `gorm:"column:classpid;type:int(10);default:0;NOT NULL" json:"classpid"`              // 栏目父id
	Classpids        string `gorm:"column:classpids;type:varchar(100)" json:"classpids"`                          // 栏目父ids
	Classchild       int    `gorm:"column:classchild;type:tinyint(1);default:0;NOT NULL" json:"classchild"`       // 是否有下级
	Classchildids    string `gorm:"column:classchildids;type:varchar(2000)" json:"classchildids"`                 // 栏目下级ids
	Classarrchildids string `gorm:"column:classarrchildids;type:mediumtext" json:"classarrchildids"`              // 栏目下级对象
	Classtitle       string `gorm:"column:classtitle;type:varchar(100);NOT NULL" json:"classtitle"`               // 栏目标题
	Classsort        int    `gorm:"column:classsort;type:smallint(5);default:0;NOT NULL" json:"classsort"`        // 排序
	Classstatus      int    `gorm:"column:classstatus;type:tinyint(1);default:1;NOT NULL" json:"classstatus"`     // 状态
	Classkeywords    string `gorm:"column:classkeywords;type:varchar(150)" json:"classkeywords"`                  // 关键字
	Classdescription string `gorm:"column:classdescription;type:varchar(200)" json:"classdescription"`            // 描述
	Classmodule      string `gorm:"column:classmodule;type:varchar(30)" json:"classmodule"`                       // 所属模型
	Classrewrite     string `gorm:"column:classrewrite;type:varchar(100)" json:"classrewrite"`                    // URL重写值
	Classtemplate    string `gorm:"column:classtemplate;type:varchar(50)" json:"classtemplate"`                   // 栏目模版
	Newstemplate     string `gorm:"column:newstemplate;type:varchar(50)" json:"newstemplate"`                     // 文章模版
	Classimg         string `gorm:"column:classimg;type:varchar(100)" json:"classimg"`                            // 栏目预览图
	Classapv         int    `gorm:"column:classapv;type:int(10);default:0;NOT NULL" json:"classapv"`              // 栏目浏览量
	Classdomain      string `gorm:"column:classdomain;type:varchar(100)" json:"classdomain"`                      // 栏目二级域名
	Classouturl      string `gorm:"column:classouturl;type:varchar(150)" json:"classouturl"`                      // 栏目外部网址
	Classmenushow    int    `gorm:"column:classmenushow;type:tinyint(1);default:1;NOT NULL" json:"classmenushow"` // 前台菜单中显示
}

func (m *Category) TableName() string {
	return "category"
}

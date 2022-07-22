package model

// 登陆记录表
type LogLogin struct {
	Id            uint   `gorm:"column:id;type:int(10) unsigned;primary_key;AUTO_INCREMENT" json:"id"`         // 主键自增
	UserId        uint   `gorm:"column:user_id;type:int(10) unsigned;default:0;NOT NULL" json:"user_id"`       // 用户ID
	LastLoginTime uint   `gorm:"column:last_login_time;type:int(10) unsigned;NOT NULL" json:"last_login_time"` // 最后的登陆时间
	LoginApp      string `gorm:"column:login_app;type:varchar(80);NOT NULL" json:"login_app"`                  // 登录的应用平台
	From          uint   `gorm:"column:from;type:tinyint(3) unsigned;default:1;NOT NULL" json:"from"`          // 来源：1为web2为android3为ios
	LastLoginIp   string `gorm:"column:last_login_ip;type:varchar(16);NOT NULL" json:"last_login_ip"`          // 最后登录的Ip
	DeviceToken   string `gorm:"column:device_token;type:varchar(200);NOT NULL" json:"device_token"`           // 登录的设备唯一标识
	Count         uint   `gorm:"column:count;type:int(10) unsigned;NOT NULL" json:"count"`                     // 登录的次数
}

func (m *LogLogin) TableName() string {
	return "log_login"
}

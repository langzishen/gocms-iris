package service

import (
	"gocms/application/model"
	"time"
)

type LogLogin struct {
	Base
}

//返回最后一次登录的信息
func (m *LogLogin) AddLogLogin(appName string, user_id uint) {
	var logLoginM model.LogLogin
	res := InitDB().Select("count").Where("user_id=? and login_app=? and `from`=1", user_id, appName).Last(&logLoginM)
	var count int
	if res.Error != nil {
		count = 1
	} else {
		count += int(logLoginM.Count)
	}

	data := model.LogLogin{UserId: user_id, LastLoginTime: uint(time.Now().Unix()), LoginApp: appName, From: 1, LastLoginIp: "0.0.0.0", DeviceToken: "", Count: uint(count)}
	InitDB().Create(&data)
}

//返回最后一次登录的信息
func (m *LogLogin) LastLogLogin(appName string, user_id uint) (int, string, int64) {
	logLoginM := model.LogLogin{}
	InitDB().Where("user_id=? and login_app=? and `from`=1", user_id, appName).Last(&logLoginM)
	var count int64
	InitDB().Model(logLoginM).Where("user_id=? and login_app=? and `from`=1", user_id, appName).Count(&count)
	return int(logLoginM.LastLoginTime), logLoginM.LastLoginIp, count
}

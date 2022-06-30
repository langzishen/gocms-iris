package model

import (
	"fmt"
	"gocms/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"strconv"
)

var DbConf = config.InitConfig().Db

//dsn := "用户名:密码@tcp(127.0.0.1:端口号)/数据库名?charset=utf8mb4&parseTime=True&loc=Local"
var Dsn = DbConf.User + ":" + DbConf.Password + "@tcp(" + DbConf.Host + ":" + strconv.Itoa(DbConf.Port) + ")/" + DbConf.Database + "?charset=utf8mb4&parseTime=True"

func InitDB() *gorm.DB {
	Db, Err := gorm.Open(mysql.Open(Dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "",
			SingularTable: true,
		},
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if Err != nil {
		//panic(Err.Error())
		fmt.Println(Err.Error())
	}
	return Db
}

package service

import (
	"encoding/json"
	"fmt"
	"gocms/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"strconv"
	"strings"
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

type Base struct {
	ModelName string
}

/**
获取列表
	list:模型切片 , 返回的数组
	field:字段切片
	order:排序：默认空字符串表示默认排序
	page:分页中的当前页,如果为-1,不进行分页处理,limit也不起作用,如果为0,不分页但limit起作用返回所有数据的limit条数,如果大于0,正常分页功能
	limit:返回的条数
*/
func (s *Base) List(list interface{}, field []string, where map[string]map[string]interface{}, order string, page int, limit int) ([]map[string]interface{}, int64) {
	tx := InitDB()
	//m := modelM.([]interface{})
	//tx = tx.Model(modelM)
	//tx = tx.Debug()
	if page > -1 { //分页处理
		if len(field) > 0 { //非空[]string{}，如果为空默认查询所有字段，即 '*'
			tx.Statement.Selects = field
			field_str := ""
			for _, Statement_v := range tx.Statement.Selects {
				field_str += "`" + Statement_v + "`" + ","
			}
			field_str = strings.Trim(field_str, ",")
			tx.Statement.Selects = []string{"SQL_CALC_FOUND_ROWS " + field_str}
		} else {
			tx.Statement.Selects = []string{"SQL_CALC_FOUND_ROWS *"}
		}
	} else { //不分页分页处理，常规查询 不加SQL_CALC_FOUND_ROWS
		if len(field) > 0 {
			tx.Select(field)
		}
	}

	if len(where) > 0 {
		for where_key, where_val := range where {
			if where_key == "string" { //不是map方式，直接是字符串，例: `filed="val AND id =1"`
				for _, where_val_val := range where_val {
					tx = tx.Where(where_val_val)
				}
			} else {
				switch strings.ToLower(where_key) {
				case "gt":
					for where_val_key, where_val_val := range where_val {
						tx = tx.Clauses(clause.Gt{Column: where_val_key, Value: where_val_val})
					}
				case "lt":
					for where_val_key, where_val_val := range where_val {
						tx = tx.Clauses(clause.Lt{Column: where_val_key, Value: where_val_val})
					}
				case "gte":
					for where_val_key, where_val_val := range where_val {
						tx = tx.Clauses(clause.Gte{Column: where_val_key, Value: where_val_val})
					}
				case "lte":
					for where_val_key, where_val_val := range where_val {
						tx = tx.Clauses(clause.Lte{Column: where_val_key, Value: where_val_val})
					}
				case "like":
					for where_val_key, where_val_val := range where_val {
						tx = tx.Clauses(clause.Like{Column: where_val_key, Value: where_val_val})
					}
				default: //默认：`=`相当于直接where(map[string]interface{})
					for where_val_key, where_val_val := range where_val {
						tx = tx.Clauses(clause.Eq{Column: where_val_key, Value: where_val_val})
					}
				}
			}
		}
	}
	if order != "" {
		tx = tx.Order(order)
	}
	if page > -1 { //page=-1时不进行分页处理,也不处理limit,返回全部数据
		if page > 0 {
			offset := (page - 1) * limit
			if offset != 0 {
				tx = tx.Offset(offset)
			}
			if limit != 0 {
				tx = tx.Limit(limit)
			}
		} else { //如果page为0时不分页,只返回limit条数数据
			if limit != 0 {
				tx = tx.Limit(limit)
			}
		}
	}
	res := tx.Find(list)

	var res_list []map[string]interface{}
	list_byte, _ := json.Marshal(list)
	json.Unmarshal(list_byte, &res_list) //返回json后的切片数组

	if page > -1 {
		type TotalCount struct {
			Count int64 `gorm:"column:total_count" json:"total_count"`
		}
		var count2 TotalCount
		//model.InitDB().Debug().Raw("SELECT FOUND_ROWS() AS `total_count`").Scan(&count2)
		InitDB().Raw("SELECT FOUND_ROWS() AS `total_count`").Scan(&count2)
		return res_list, count2.Count
	} else {
		var count int64
		count = res.RowsAffected
		return res_list, count
	}
}

# gocms

gocms是go语言实现的一套cms内容管理系统,服务端使用go语言,主体框架使用[iris](https://github.com/kataras/iris)框架,采用mvc架构,管理后台前端使用[topjui](https://www.topjui.com/),数据库使用mysql,ORM使用[goorm](https://gorm.io),内置RBAC、UPLOAD、OSS等扩展。

#### 环境搭建：

	1.go version:1.16+
	2.首次运行需修改项目根目录下config/config.json配置文件。项目运行起来后可以在后台的系统管理中修改部分配置。
	3.导入数据库文件:db/gocms.sql到mysql中,数据库名称需与config.json保持一致
	4.启动方式：go run main.go
	5.前台访问地址：域名/index/index/index
	6.后台访问地址：域名/boss/index/index,(默认)，也可以在main.go中修改`app.Party("/boss", controller.AdminAuth)`,中的`“/boss”`为其他路径。
	7.后台的超级管理员账号/密码：admin/amin
	8.`/static/uploads`、`sessions.db`需要有写入权限。



#### 主要包依赖：

* github.com/kataras/iris/v12
* github.com/golang/freetype
* gorm.io/driver/mysql
* github.com/go-sql-driver/mysql
* gorm.io/gorm
* github.com/aliyun/aliyun-oss-go-sdk